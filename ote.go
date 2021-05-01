// // ote: updates a packages go.mod file with a comment next to all dependencies that are test dependencies; identifying them as such.
// //
// // It maybe useful in places where it is important to audit all dependencies that are going to run in production.
// //
// // Install: go get github.com/komuw/ote
// // Usage: ote .
// //
// package main

// import (
// 	"flag"
// 	"fmt"
// 	"io"
// 	"log"
// 	"os"
// 	"path/filepath"
// 	"strings"

// 	"golang.org/x/mod/modfile"
// 	"golang.org/x/tools/go/packages"
// )

// // This lists are got from https://github.com/golang/go/blob/master/src/go/build/syslist.go
// // They should be synced periodically
// const goosList = `aix
// 				  android
// 				  darwin
// 				  dragonfly
// 				  freebsd
// 				  hurd
// 				  illumos
// 				  js
// 				  linux
// 				  nacl
// 				  netbsd
// 				  openbsd
// 				  plan9
// 				  solaris
// 				  windows
// 				  zos `
// const goarchList = `386
// 					amd64
// 					amd64p32
// 					arm
// 					armbe
// 					arm64
// 					arm64be
// 					ppc64
// 					ppc64le
// 					mips
// 					mipsle
// 					mips64
// 					mips64le
// 					mips64p32
// 					mips64p32le
// 					ppc
// 					riscv
// 					riscv64
// 					s390
// 					s390x
// 					sparc
// 					sparc64
// 					wasm `
// const cGo = "cgo"

// func getModFile(gomodFile string) (*modfile.File, error) {
// 	modContents, err := os.ReadFile(filepath.Clean(gomodFile))
// 	if err != nil {
// 		return nil, err
// 	}
// 	f, err := modfile.Parse(gomodFile, modContents, nil)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return f, nil
// }

// func getPackage(pattern string, gomodFile string, mainModule bool) (*packages.Package, error) {
// 	patterns := []string{fmt.Sprintf("pattern=%s", pattern)}
// 	buildFlags := (strings.Join(strings.Split(goosList, " "), ",") +
// 		strings.Join(strings.Split(goarchList, " "), ",") +
// 		cGo)

// 	pkgNeeds := packages.NeedModule
// 	if mainModule {
// 		pkgNeeds = packages.NeedImports | packages.NeedModule
// 	}
// 	cfg := &packages.Config{
// 		Mode:       pkgNeeds,
// 		Tests:      false,
// 		BuildFlags: []string{fmt.Sprintf("-tags=%s", buildFlags)},
// 		Dir:        filepath.Dir(gomodFile),
// 	}
// 	pkgs, err := packages.Load(cfg, patterns...)
// 	if err != nil {
// 		return nil, err
// 	}
// 	pkg := pkgs[0]

// 	return pkg, nil
// }

// // getModules finds all the modules that have been used/imported by a module
// func getModules(pattern string, gomodFile string) ([]string, error) {
// 	modulePaths := []string{}

// 	mainPkg, err := getPackage(pattern, gomodFile, true)
// 	if err != nil {
// 		return modulePaths, err
// 	}
// 	impPaths := make([]string, 0, len(mainPkg.Imports))
// 	for impPath := range mainPkg.Imports {
// 		impPaths = append(impPaths, impPath)
// 	}

// 	for _, v := range impPaths {
// 		pkg, err := getPackage(v, gomodFile, false)
// 		if err != nil {
// 			// maybe we should continue, instead of returning?
// 			return modulePaths, err
// 		}
// 		if pkg.Module == nil {
// 			// something like `fmt` has a `.Module` that is nil
// 			continue
// 		}

// 		if pkg.Module.Path != "" {
// 			modulePaths = append(modulePaths, pkg.Module.Path)
// 		}
// 	}

// 	return modulePaths, nil
// }

// // getDeps finds all the dependencies of a given module
// func getDeps(gomodFile string) ([]modfile.Require, error) {
// 	requires := []modfile.Require{}

// 	modContents, err := os.ReadFile(filepath.Clean(gomodFile))
// 	if err != nil {
// 		return requires, err
// 	}

// 	f, err := modfile.Parse(gomodFile, modContents, nil)
// 	if err != nil {
// 		return requires, err
// 	}

// 	for _, v := range f.Require {
// 		requires = append(requires, *v)
// 	}

// 	return requires, nil
// }

// func getTestDeps(impPaths []string, allDeps []modfile.Require) []modfile.Require {
// 	// find which deps exist in `allDeps` but not in `impPaths`; this are the test dependencies

// 	testRequires := []modfile.Require{}

// 	// contains tells whether a contains x.
// 	contains := func(a []string, x string) bool {
// 		for _, n := range a {
// 			if x == n {
// 				return true
// 			}
// 		}
// 		return false
// 	}

// 	for _, r := range allDeps {
// 		if !contains(impPaths, r.Mod.Path) {
// 			testRequires = append(testRequires, r)
// 		}
// 	}

// 	return testRequires
// }

// func updateMod(testRequires []modfile.Require, f *modfile.File) error {
// 	notIndirect := []modfile.Require{}
// 	for _, v := range testRequires {
// 		// we do not want to add a `//test` comment to any requires that already have `//indirect` comment
// 		if !v.Indirect {
// 			notIndirect = append(notIndirect, v)
// 		}
// 	}

// 	if len(testRequires) < 0 {
// 		// if there are no test dependencies, we need to go through all the deps and
// 		// remove any test comments that there may be there.
// 		for _, fr := range f.Require {
// 			line := fr.Syntax
// 			setTest(line, false)
// 		}

// 		return nil
// 	}

// 	for _, ni := range notIndirect {
// 		for _, fr := range f.Require {
// 			if ni.Mod == fr.Mod {
// 				// add test comment
// 				line := fr.Syntax
// 				setTest(line, true)
// 			}
// 		}
// 	}

// 	// contains tells whether a contains x.
// 	contains := func(a []modfile.Require, x modfile.Require) bool {
// 		for _, n := range a {
// 			if x.Mod == n.Mod {
// 				return true
// 			}
// 		}
// 		return false
// 	}

// 	for _, fr := range f.Require {
// 		line := fr.Syntax
// 		if isTest(line) {
// 			if !contains(testRequires, *fr) {
// 				// Remove test comment for any module that may be used in both test files and non-test files.
// 				// If a module has a test comment but is not in testRequires, it should be removed.
// 				setTest(line, false)
// 			}
// 		}
// 	}

// 	return nil
// }

// func writeMod(f *modfile.File, gomodFile string, w io.Writer, readonly bool) error {
// 	f.SortBlocks()
// 	f.Cleanup()
// 	b, err := f.Format()
// 	if err != nil {
// 		return err
// 	}

// 	i, err := os.Stat(gomodFile)
// 	if err != nil {
// 		return err
// 	}

// 	if readonly {
// 		fmt.Fprintln(w, string(b))
// 	} else {
// 		fi, err := os.OpenFile(gomodFile, os.O_RDWR, i.Mode())
// 		if err != nil {
// 			return err
// 		}
// 		_, err = fi.Write(b)
// 		fmt.Fprintln(w, "successfully updated go.mod file.")
// 		return err
// 	}

// 	return nil
// }

// func main() {
// 	f, r := cli()

// 	err := run(f, os.Stdout, r)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }

// func cli() (string, bool) {
// 	var f string
// 	var r bool

// 	flag.Usage = func() {
// 		fmt.Fprintf(flag.CommandLine.Output(),
// 			`ote updates a packages go.mod file with a comment next to all dependencies that are test dependencies; identifying them as such.

// Usage:

// -f string
// 	path to directory containing the go.mod file. By default, it uses the current directory. (default ".")
// -r
//         (readonly) write to stdout instead of updating go.mod file.

// examples:
// 	ote .
// 		update go.mod in the current directory
// 	ote -f /tmp/myPkg
// 		update go.mod in the /tmp/myPkg directory

// 	ote -r
// 		(readonly) write to stdout instead of updating go.mod file.
// 	ote -f /tmp/myPkg -r
// 	        (readonly) write to stdout instead of updating go.mod file in the /tmp/myPkg directory.

// 	`)
// 	}
// 	flag.StringVar(
// 		&f,
// 		"f",
// 		".",
// 		"path to directory containing the go.mod file. By default, it uses the current directory.")
// 	flag.BoolVar(
// 		&r,
// 		"r",
// 		false,
// 		"(readonly) display how the updated go.mod file would look like, without actually updating the file.")
// 	flag.Parse()

// 	return f, r
// }

// func run(fp string, w io.Writer, readonly bool) error {
// 	gomodFile := filepath.Join(fp, "go.mod")

// 	f, err := getModFile(gomodFile)
// 	if err != nil {
// 		return err
// 	}
// 	thisMod := f.Module.Mod.Path

// 	modulePaths, err := getModules(thisMod, gomodFile)
// 	if err != nil {
// 		return err
// 	}

// 	allDeps, err := getDeps(gomodFile)
// 	if err != nil {
// 		return err
// 	}

// 	testRequires := getTestDeps(modulePaths, allDeps)
// 	err = updateMod(testRequires, f)
// 	if err != nil {
// 		return err
// 	}

// 	err = writeMod(f, gomodFile, w, readonly)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

// TODO: better errors

func main() {
	f, r := cli()

	err := run(f, os.Stdout, r)
	if err != nil {
		log.Fatal(err)
	}
}

func run(fp string, w io.Writer, readonly bool) error {
	err := loadStd()
	if err != nil {
		log.Fatal("unable to load Stdlibs: ", err)
	}

	gomodFile := filepath.Join(fp, "go.mod")
	f, err := getModFile(gomodFile)
	if err != nil {
		return err
	}

	trueTestModules, err := getTestModules(fp)
	if err != nil {
		return err
	}
	err = updateMod(trueTestModules, f)
	if err != nil {
		return err
	}

	err = writeMod(f, gomodFile, w, readonly)
	if err != nil {
		return err
	}

	return nil
}

// TODO: remove this:
// Usage:
//   go run . -f testdata/mod2/ -r

func cli() (string, bool) {
	var f string
	var r bool

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(),
			`ote updates a packages go.mod file with a comment next to all dependencies that are test dependencies; identifying them as such.
Usage:
-f string
	path to directory containing the go.mod file. By default, it uses the current directory. (default ".")
-r	
        (readonly) write to stdout instead of updating go.mod file.
examples:
	ote .
		update go.mod in the current directory
	ote -f /tmp/myPkg
		update go.mod in the /tmp/myPkg directory
	ote -r
		(readonly) write to stdout instead of updating go.mod file.
	ote -f /tmp/myPkg -r
	        (readonly) write to stdout instead of updating go.mod file in the /tmp/myPkg directory.
	`)
	}
	flag.StringVar(
		&f,
		"f",
		".",
		"path to directory containing the go.mod file. By default, it uses the current directory.")
	flag.BoolVar(
		&r,
		"r",
		false,
		"(readonly) display how the updated go.mod file would look like, without actually updating the file.")
	flag.Parse()

	return f, r
}
