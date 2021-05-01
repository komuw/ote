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

func walkDirFn(path string, d fs.DirEntry, err error) error {
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
	if filepath.Ext(fName) != ".go" {
		return nil
	}

	impPaths, errF := fetchImports(path)
	if errF != nil {
		return errF
	}
	if strings.Contains(fName, "_test.go") {
		testImportPaths = append(testImportPaths, impPaths...)
	} else {
		nonTestImportPaths = append(nonTestImportPaths, impPaths...)
	}

	return nil
}

var (
	//TODO: turn into
	// type importPaths string
	// testImportPaths    = []importPaths{}

	testImportPaths    = []string{}
	nonTestImportPaths = []string{}
)

func fetchImports(file string) ([]string, error) {
	fset := token.NewFileSet()
	// file := "/Users/komuw/mystuff/ote/testdata/mod2/main.go"
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

var stdLibPkgs = make(map[string]struct{})

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

//
// Usage:
//     fetchModule("/Users/komuw/mystuff/ote/testdata/mod2/", "github.com/hashicorp/nomad/drivers/shared/executor")
func fetchModule(root, importPath string) (string, error) {
	cfg := &packages.Config{
		Mode:  packages.NeedModule,
		Tests: false,
		// BuildFlags: []string{fmt.Sprintf("-tags=%s", buildFlags)},
		Dir: filepath.Dir(root),
	}
	pkgs, err := packages.Load(
		cfg,
		fmt.Sprintf("pattern=%s", importPath),
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
	// testModules := []string{}
	// nonTestModules := []string{}

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
