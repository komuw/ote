package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"golang.org/x/mod/modfile"
)

type lineMod struct {
	name string
	ver  string
	coms modfile.Comments
}

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
	if len(trueTestModules) <= 0 {
		// if there are no test dependencies, we need to go through all the deps and
		// remove any test comments that there may be there.
		for _, fr := range f.Require {
			line := fr.Syntax
			setTest(line, false)
		}
		return nil
	}

	lineMods := []lineMod{}
	for _, ni := range trueTestModules {
		for _, fr := range f.Require {
			if ni == fr.Mod.Path {
				// add test comment
				line := fr.Syntax
				setTest(line, true)
				lineMods = append(lineMods, lineMod{name: line.Token[0], ver: line.Token[1], coms: line.Comments})
			}
		}
	}

	for _, fr := range f.Require {
		line := fr.Syntax
		if isTest(line) {
			if !contains(trueTestModules, fr.Mod.Path) {
				// Remove test comment for any module that may be used in both test files and non-test files.
				// If a module has a test comment but is not in testRequires, it should be removed.
				setTest(line, false)
			}
		}
	}

	addTestRequireBlock(f, lineMods)

	return nil
}

func addTestRequireBlock(f *modfile.File, lineMods []lineMod) {
	// Add a new require block after the last "require".
	// This new block will house test-only requirements
	// eg.
	/*
		require (
				github.com/fatih/color v1.12.0 // test
				github.com/frankban/quicktest v1.12.1 // test
			)
	*/

	if len(lineMods) <= 0 {
		return
	}
	for _, y := range lineMods {
		// since test-only deps are in their own require blocks,
		// drop them from the main one.
		_ = f.DropRequire(y.name)
	}

	testLines := []*modfile.Line{}
	for _, y := range lineMods {
		testLines = append(testLines, &modfile.Line{
			Token:    []string{y.name, y.ver},
			Comments: y.coms,
		})
	}
	newTestBlock := &modfile.LineBlock{
		Token: []string{"require"},
		Line:  testLines,
	}
	f.Syntax.Stmt = append(f.Syntax.Stmt, newTestBlock)
	f.Syntax.Cleanup()
}

// writeMod updates the on-disk modfile
func writeMod(f *modfile.File, gomodFile string, w io.Writer, readonly bool) error {
	f.SortBlocks()
	f.Cleanup()
	b, errF := f.Format()
	if errF != nil {
		return errF
	}

	if readonly {
		fmt.Fprintln(w, string(b))
	} else {
		i, errS := os.Stat(gomodFile)
		if errS != nil {
			return errS
		}

		fi, errO := os.OpenFile(gomodFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, i.Mode())
		if errO != nil {
			return errO
		}

		_, errW := fi.Write(b)
		if errW != nil {
			return errW
		}

		errC := fi.Close()
		if errC != nil {
			return errC
		}

		fmt.Fprintln(w, "successfully updated go.mod file.")
		return errC
	}

	return nil
}
