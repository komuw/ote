package main

import (
	"fmt"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"strings"

	"golang.org/x/mod/modfile"
	"golang.org/x/mod/module"
	"golang.org/x/tools/go/packages"
)

// This lists are got from https://github.com/golang/go/blob/master/src/go/build/syslist.go
// They should be synced periodically
const goosList = "aix android darwin dragonfly freebsd hurd illumos ios js linux nacl netbsd openbsd plan9 solaris windows zos "
const goarchList = "386 amd64 amd64p32 arm armbe arm64 arm64be ppc64 ppc64le mips mipsle mips64 mips64le mips64p32 mips64p32le ppc riscv riscv64 s390 s390x sparc sparc64 wasm "
const cGo = "cgo"

var stdLibPkgs = map[string]struct{}{
	"C": {}, // cGo. see: https://blog.golang.org/cgo
}

func loadStd() error {
	pkgs, err := packages.Load(nil, "std")
	if err != nil {
		return err
	}
	for _, p := range pkgs {
		stdLibPkgs[p.PkgPath] = struct{}{}
	}
	return nil
}

func isStdLibPkg(pkg string) bool {
	_, ok := stdLibPkgs[pkg]
	return ok
}

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
				if !isStdLibPkg(p) {
					impPaths = append(impPaths, p)
				}
			}
		}
	}

	return impPaths, nil
}

//
// Usage:
//     fetchModule("testdata/mod1/", "github.com/hashicorp/nomad/drivers/shared/executor")
func fetchModule(root, importPath string) (string, error) {
	buildFlags := (strings.Join(strings.Split(goosList, " "), ",") +
		strings.Join(strings.Split(goarchList, " "), ",") +
		cGo)
	cfg := &packages.Config{
		Mode:       packages.NeedModule,
		Tests:      false,
		BuildFlags: []string{fmt.Sprintf("-tags=%s", buildFlags)},
		Dir:        root,
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
	if len(pkgs) < 0 {
		return "", fmt.Errorf("import %s does not belong to any package", importPath)
	}

	pkg := pkgs[0]
	if pkg.Module == nil {
		return "", fmt.Errorf("import %s does not belong to any module", importPath)
	}

	mRequire := modfile.Require{
		Mod:      module.Version{Path: pkg.Module.Path, Version: pkg.Module.Version},
		Indirect: pkg.Module.Indirect,
	}
	// TODO: remove this
	_ = mRequire

	return pkg.Module.Path, nil
}

func getAllmodules(testImportPaths []string, nonTestImportPaths []string, root string) (testModules []string, nonTestModules []string, err error) {
	// todo: these two for loops can be made concurrent.

	for _, v := range testImportPaths {
		m, err := fetchModule(root, v)
		if err != nil {
			return testModules, nonTestModules, err
		}
		testModules = append(testModules, m)
	}

	for _, v := range nonTestImportPaths {
		m, err := fetchModule(root, v)
		if err != nil {
			return testModules, nonTestModules, err
		}
		nonTestModules = append(nonTestModules, m)
	}

	return dedupe(testModules), dedupe(nonTestModules), nil
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

			allGoFiles = append(allGoFiles, path)
			return nil
		},
	)
	if err != nil {
		return []string{}, err
	}

	fetchToAnalyze := func() []string {
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
		return tobeAnalyzed
	}

	tobeAnalyzed := fetchToAnalyze()
	fetchPaths := func() ([]string, []string, error) {
		// TODO: turn into
		// type importPaths string
		// testImportPaths = []importPaths{}

		testImportPaths := []string{}
		nonTestImportPaths := []string{}

		for _, path := range tobeAnalyzed {
			impPaths, errF := fetchImports(path)
			if errF != nil {
				return []string{}, []string{}, errF
			}

			if strings.Contains(path, "_test.go") {
				// this takes care of both;
				// (i) test files
				// (ii) example files(https://blog.golang.org/examples)
				testImportPaths = append(testImportPaths, impPaths...)
			} else {
				nonTestImportPaths = append(nonTestImportPaths, impPaths...)
			}
		}
		return testImportPaths, nonTestImportPaths, nil
	}

	testImportPaths, nonTestImportPaths, err := fetchPaths()
	if err != nil {
		return []string{}, err
	}

	testModules, nonTestModules, err := getAllmodules(testImportPaths, nonTestImportPaths, root)
	if err != nil {
		return []string{}, err
	}
	trueTestModules := difference(testModules, nonTestModules)

	return trueTestModules, nil
}
