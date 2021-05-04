// ote: updates a package's go.mod file with a comment next to all dependencies that are test dependencies; identifying them as such.
//
// It maybe useful in places where it is important to audit all dependencies that are going to run in production.
//
// Install:
//    go install github.com/komuw/ote@latest
//
// Usage:
//    ote .
//
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

// Usage:
//   1.
//    go run . -f testdata/mod1/ -r
//
//   2.
//    export export GOPACKAGESDEBUG=true && \
//    go run . -f testdata/mod1/ -r
func main() {
	f, r := cli()

	err := run(f, os.Stdout, r)
	if err != nil {
		log.Fatal(err)
	}
}

func run(fp string, w io.Writer, readonly bool) error {
	e := loadStd()
	if e != nil {
		return e
	}

	gomodFile := filepath.Join(fp, "go.mod")
	f, errM := getModFile(gomodFile)
	if errM != nil {
		return errM
	}
	// `f.Cleanup` will be called inside `writeMod`

	trueTestModules, errT := getTestModules(fp)
	if errT != nil {
		return errT
	}
	errU := updateMod(trueTestModules, f)
	if errU != nil {
		return errU
	}

	errW := writeMod(f, gomodFile, w, readonly)
	if errW != nil {
		return errW
	}

	return nil
}

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
