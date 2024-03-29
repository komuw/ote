package main

import (
	"errors"
	"fmt"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"slices"
	"strings"
	"sync"

	"golang.org/x/tools/go/packages"
)

const stdlib = "std"

// once is used to ensure that the stdLibPkgs map is populated only once
var (
	once       = &sync.Once{}         //nolint:gochecknoglobals
	stdLibPkgs = map[string]struct{}{ //nolint:gochecknoglobals
		"C": {}, // cGo. see: https://blog.golang.org/cgo
	}
)

func isStdLibPkg(pkg, std string) (bool, error) {
	var err error
	once.Do(func() {
		pkgs, errL := packages.Load(nil, std)
		if errL != nil {
			// set stdLibPkgs to empty when error occurs.
			stdLibPkgs = map[string]struct{}{}
			err = errL
		}
		if len(pkgs) < 10 {
			// it means an error occured since
			// we will always have more than 10 pkgs in the Go standard lib
			for _, p := range pkgs {
				if len(p.Errors) > 0 {
					err = errors.New(p.Errors[0].Msg)
				}
			}
		}
		for _, p := range pkgs {
			stdLibPkgs[p.PkgPath] = struct{}{}
		}
	})

	_, ok := stdLibPkgs[pkg]
	return ok, err
}

// fetchImports returns all the imports found in one .go file
func fetchImports(file string) ([]string, error) {
	fset := token.NewFileSet()
	var src interface{} = nil
	mode := parser.ImportsOnly
	f, err := parser.ParseFile(fset, file, src, mode)
	if err != nil {
		return nil, err
	}

	impPaths := make([]string, 0)
	for _, impPath := range f.Imports {
		if impPath != nil {
			if impPath.Path != nil {
				p := impPath.Path.Value
				p = strings.Trim(p, "\"")
				isStdlib, _ := isStdLibPkg(p, stdlib) // ignore error, since when error happens; the other value will be false
				if !isStdlib {
					impPaths = append(impPaths, p)
				}
			}
		}
	}

	return impPaths, nil
}

// Usage:
//
//	fetchModule("testdata/modfiles/mod1/", "github.com/hashicorp/nomad/drivers/shared/executor")
func fetchModule(root, importPath string) (string, error) {
	cfg := &packages.Config{
		Mode:  packages.NeedModule,
		Tests: false,
		Dir:   root,
	}
	pkgs, err := packages.Load(
		cfg,
		fmt.Sprintf("pattern=%s", importPath),
		// to show the Go commands that this method call will invoke;
		// run ote while the env var `export GOPACKAGESDEBUG=true` is set on the commandline
	)
	if err != nil {
		return "", err
	}

	if len(pkgs) > 1 {
		return "", fmt.Errorf("import %s produced greater than 1 packages", importPath)
	}
	if len(pkgs) <= 0 {
		return "", fmt.Errorf("import %s does not belong to any package", importPath)
	}

	pkg := pkgs[0]
	if pkg.Module == nil {
		if len(pkg.Errors) > 0 {
			if pkg.Errors[0].Kind == packages.ListError && strings.Contains(pkg.Errors[0].Msg, "build constraints exclude all Go files") {
				// If a package depends on an import that requires an explicit build tag,
				// then, `packages.Load()` is going to fail. The only way to fix it would be to pass in
				// the explicit build tag required by that dependency to `packages.Config.BuildFlags`
				// However, we do not know in advance what the required build tags are. That's why the call
				// to `packages.Load()/packages.Config` do not have any build tag specified.
				//
				// As an example, `/testdata/modfiles/mod1/` depends on `golang.org/x/sys/windows`. That dependency
				// has the build tag; `+build windows`: https://github.com/golang/sys/blob/0f9fa26af87c481a6877a4ca1330699ba9a30673/windows/aliases.go#L5-L8
				// and thus `packages.Load()` fails with error;
				// `build constraints exclude all Go files in /pkg/mod/golang.org/x/sys@v0.0.1/windows`
				//
				// We can always add more errors here as/when we discover them
				return "", nil
			} else {
				return "", fmt.Errorf("unable to find module for import %s : %s", importPath, pkg.Errors[0].Msg)
			}
		} else {
			// this can be raised if an import path is inside a file that has some build tag
			// that ote didn't take into account.
			// see: https://github.com/komuw/ote/issues/3
			return "", fmt.Errorf("import %s does not belong to any module", importPath)
		}
	}

	return pkg.Module.Path, nil
}

// getAllTestModules finds all the Go modules used only in test files.
func getAllTestModules(testImportPaths, nonTestImportPaths []string, root string) (testModules []string, err error) {
	// There could be some import paths that exist in both test files & non-test files.
	// In hashicorp/nomad we found that to be about 50% of imports.
	// In juju/juju it is about 80%
	// see: https://github.com/komuw/ote/issues/22
	//
	// Given that, it then only makes sense to filter out this import paths that are common
	// before calling fetchModule(which is one of the most expensive calls in ote)

	testOnlyMods := []string{}
	for _, v := range testImportPaths {
		m, errF := fetchModule(root, v)
		if errF != nil {
			return testModules, errF
		}
		testOnlyMods = append(testOnlyMods, m)
	}

	nonTestMods := []string{}
	for _, v := range nonTestImportPaths {
		m, errF := fetchModule(root, v)
		if errF != nil {
			return testModules, errF
		}
		nonTestMods = append(nonTestMods, m)
	}

	existsInBoth := []string{}
	for _, a := range nonTestMods {
		if slices.Contains(testOnlyMods, a) {
			existsInBoth = append(existsInBoth, a)
		}
	}

	testModules = difference(testOnlyMods, existsInBoth)

	return dedupe(testModules), nil
}

func getTestModules(root string) ([]string, error) {
	allGoFiles := []string{}
	nonMainModFileDirs := []string{}
	err := filepath.WalkDir(
		// note: WalkDir reads an entire directory into memory before proceeding to walk that directory.
		// see documentation of filepath.WalkDir
		root,
		func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if d.IsDir() {
				// return on directories since we don't want to parse them
				return nil
			}
			if !d.Type().IsRegular() {
				// non regular files. nothing to parse
				return nil
			}

			fName := d.Name()
			if filepath.Ext(fName) == ".mod" {
				if path != filepath.Join(root, "go.mod") {
					nonMainModFileDirs = append(nonMainModFileDirs, filepath.Dir(path))
				}
			}
			if filepath.Ext(fName) != ".go" {
				return nil
			}
			if strings.Contains(path, "vendor/") {
				// ignore files inside vendor/ directory
				return nil
			}

			allGoFiles = append(allGoFiles, path)
			return nil
		},
	)
	if err != nil {
		return []string{}, err
	}

	filesTobeAnalyzed := fetchToAnalyze(allGoFiles, nonMainModFileDirs)
	testImportPaths, nonTestImportPaths, err := getAllImports(filesTobeAnalyzed)
	if err != nil {
		return []string{}, err
	}

	trueTestModules, err := getAllTestModules(testImportPaths, nonTestImportPaths, root)
	if err != nil {
		return []string{}, err
	}

	return trueTestModules, nil
}

// fetchToAnalyze returns only the list of Go files that need to be analyzed.
// it excludes files that are in directories which are nested go modules.
func fetchToAnalyze(allGoFiles, nonMainModFileDirs []string) []string {
	notToBetAnalzyed := []string{}
	for _, goFile := range allGoFiles {
		for _, mod := range nonMainModFileDirs {
			if strings.Contains(goFile, mod) {
				// this file should not be analyzed because it belongs
				// to a nested module
				notToBetAnalzyed = append(notToBetAnalzyed, goFile)
			}
		}
	}

	tobeAnalyzed := difference(allGoFiles, notToBetAnalzyed)
	return tobeAnalyzed // no need to dedupe. files are unlikely to be duplicates.
}

// getAllImports aggregates all imports from a list of .go files
func getAllImports(files []string) ([]string, []string, error) {
	// TODO: turn into
	// type importPaths string
	// testImportPaths = []importPaths{}

	testImportPaths := []string{}
	nonTestImportPaths := []string{}

	for _, filePath := range files {
		impPaths, errF := fetchImports(filePath)
		if errF != nil {
			return []string{}, []string{}, errF
		}

		if strings.Contains(filePath, "_test.go") {
			// this takes care of both;
			// (i) test files
			// (ii) example files(https://blog.golang.org/examples)
			testImportPaths = append(testImportPaths, impPaths...)
		} else {
			nonTestImportPaths = append(nonTestImportPaths, impPaths...)
		}
	}

	// dedupe, since one importPath is likely to have been used in multiple Go files.
	return dedupe(testImportPaths), dedupe(nonTestImportPaths), nil
}
