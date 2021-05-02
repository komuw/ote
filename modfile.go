package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"golang.org/x/mod/modfile"
)

func getModFile(gomodFile string) (*modfile.File, error) {
	modContents, err := os.ReadFile(filepath.Clean(gomodFile))
	if err != nil {
		return nil, err
	}
	f, err := modfile.Parse(gomodFile, modContents, nil)
	if err != nil {
		return nil, err
	}

	return f, nil
}

// updateMod updates the in-memory modfile
func updateMod(trueTestModules []string, f *modfile.File) error {
	if len(trueTestModules) < 0 {
		// if there are no test dependencies, we need to go through all the deps and
		// remove any test comments that there may be there.
		for _, fr := range f.Require {
			line := fr.Syntax
			setTest(line, false)
		}
		return nil
	}

	for _, ni := range trueTestModules {
		for _, fr := range f.Require {
			if ni == fr.Mod.Path {
				// add test comment
				line := fr.Syntax
				setTest(line, true)
			}
		}
	}

	// contains tells whether a contains x.
	contains := func(a []string, x modfile.Require) bool {
		for _, n := range a {
			if x.Mod.Path == n {
				return true
			}
		}
		return false
	}

	for _, fr := range f.Require {
		line := fr.Syntax
		if isTest(line) {
			if !contains(trueTestModules, *fr) {
				// Remove test comment for any module that may be used in both test files and non-test files.
				// If a module has a test comment but is not in testRequires, it should be removed.
				setTest(line, false)
			}
		}
	}

	return nil
}

// writeMod updates the on-disk modfile
func writeMod(f *modfile.File, gomodFile string, w io.Writer, readonly bool) error {
	f.SortBlocks()
	f.Cleanup()
	b, err := f.Format()
	if err != nil {
		return err
	}

	i, err := os.Stat(gomodFile)
	if err != nil {
		return err
	}

	if readonly {
		fmt.Fprintln(w, string(b))
	} else {
		fi, err := os.OpenFile(gomodFile, os.O_RDWR, i.Mode())
		if err != nil {
			return err
		}
		defer fi.Close()

		_, err = fi.Write(b)
		fmt.Fprintln(w, "successfully updated go.mod file.")
		return err
	}

	return nil
}
