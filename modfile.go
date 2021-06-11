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

	testLines := []*modfile.Line{}
	for _, ni := range trueTestModules {
		for _, fr := range f.Require {
			if ni == fr.Mod.Path {
				// add test comment
				line := fr.Syntax
				setTest(line, true)
				testLines = append(testLines, line)
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

	addTestRequireBlock(f, testLines)

	return nil
}

func addTestRequireBlock(f *modfile.File, testLines []*modfile.Line) {
	// Add a new require block after the last "require".
	// This new block will house test-only requirements
	// eg.
	/*
		require (
				github.com/fatih/color v1.12.0 // test
				github.com/frankban/quicktest v1.12.1 // test
				github.com/go-xorm/builder v0.3.4 // test
			)
	*/

	if len(testLines) <= 0 {
		return
	}

	type aha struct {
		name string
		ver  string
		com  string
	}
	yes := []aha{}
	for _, t := range testLines {
		yes = append(yes, aha{name: t.Token[0], ver: t.Token[1], com: testComment})
	}

	for _, y := range yes {
		// since test-only deps are in their own require blocks,
		// drop them from the main one.
		f.DropRequire(y.name)
	}

	xx := []*modfile.Line{}
	for _, y := range yes {
		xx = append(xx,
			&modfile.Line{Token: []string{y.name, y.ver, y.com}},
		)
	}

	newTestBlock := &modfile.LineBlock{
		Token: []string{"require"},
		Line:  xx,
		// Line: []*modfile.Line{
		// 	&modfile.Line{Token: []string{"github.com/fatih/color", "v1.12.0", "// test"}},
		// 	&modfile.Line{Token: []string{"github.com/go-xorm/builder", "v0.3.4", "// test"}},
		// 	&modfile.Line{Token: []string{"github.com/frankban/quicktest", "v1.12.1", "// test"}},
		// },
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
