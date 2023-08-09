// ote: updates a package's go.mod file with a comment next to all dependencies that are test dependencies; identifying them as such.
//
// It maybe useful in places where it is important to audit all dependencies that are going to run in production.
//
// Install:
//
//	go install github.com/komuw/ote@latest
//
// Usage:
//
//	ote .
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
//
//	1.
//	 go run . -f testdata/modfiles/mod1/ -r
//
//	2.
//	 export export GOPACKAGESDEBUG=true && \
//	 go run . -f testdata/modfiles/mod1/ -r
func main() {
	f, r, v := cli()
	if v {
		// show version
		fmt.Println(version())
		return
	}

	err := run(f, os.Stdout, r)
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}

func run(fp string, w io.Writer, readonly bool) error {
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

	if err := writeMod(f, gomodFile, w, readonly); err != nil {
		return err
	}

	return nil
}

func cli() (string, bool, bool) {
	var v bool
	var f string
	var r bool

	flag.Usage = func() {
		_, _ = fmt.Fprintf(flag.CommandLine.Output(),
			`ote updates a packages go.mod file with a comment next to all dependencies that are test dependencies; identifying them as such.

Usage:
-f string  # path to directory containing the go.mod file. By default, it uses the current directory. (default ".")
-r bool    # (readonly) write to stdout instead of updating go.mod file.
-v bool    # display version of ote in use.

examples:
	ote .                    # update go.mod in the current directory
	ote -f /tmp/someDir      # update go.mod in the /tmp/someDir directory
	ote -r                   # (readonly) write to stdout instead of updating go.mod file.
	ote -f /tmp/someDir -r   # (readonly) write to stdout instead of updating go.mod file in the /tmp/someDir directory.
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
	flag.BoolVar(
		&v,
		"v",
		false,
		"display version of ote in use.")
	flag.Parse()

	return f, r, v
}
