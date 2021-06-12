package main

import (
	"errors"
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
				if len(line.Token) == 2 {
					// eg; `github.com/shirou/gopsutil v1.21.5`
					lineMods = append(lineMods, lineMod{name: line.Token[0], ver: line.Token[1], coms: line.Comments})
				} else if len(line.Token) == 3 {
					// eg; `require rsc.io/quote v1.5.2`
					lineMods = append(lineMods, lineMod{name: line.Token[1], ver: line.Token[2], coms: line.Comments})
				} else {
					return errors.New("unrecognised line.Token format")
				}
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

	findLastRequire := func(f *modfile.File) int {
		/*
			The code of this function is mostly file is taken from https://github.com/golang/mod/
			The license of that repo is found at: https://github.com/golang/mod/blob/v0.3.0/LICENSE
			and is included below:

				Copyright (c) 2009 The Go Authors. All rights reserved.

				Redistribution and use in source and binary forms, with or without
				modification, are permitted provided that the following conditions are
				met:

				   * Redistributions of source code must retain the above copyright
				notice, this list of conditions and the following disclaimer.
				   * Redistributions in binary form must reproduce the above
				copyright notice, this list of conditions and the following disclaimer
				in the documentation and/or other materials provided with the
				distribution.
				   * Neither the name of Google Inc. nor the names of its
				contributors may be used to endorse or promote products derived from
				this software without specific prior written permission.

				THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
				"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
				LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
				A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
				OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
				SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
				LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
				DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
				THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
				(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
				OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
		*/

		var pos int

		for i, stmt := range f.Syntax.Stmt {
			switch stmt := stmt.(type) {
			case *modfile.Line:
				if len(stmt.Token) == 0 || stmt.Token[0] != "require" {
					continue
				}
				pos = i
			case *modfile.LineBlock:
				if len(stmt.Token) == 0 || stmt.Token[0] != "require" {
					continue
				}
				pos = i
			}
		}

		return pos
	}

	insertAt := func(slice []modfile.Expr, index int, value *modfile.LineBlock) []modfile.Expr {
		/*
			Insert inserts the value into the slice at the specified index,
			which must be in range.
			from: https://github.com/golang/go/wiki/SliceTricks#insert
		*/
		slice = append(slice, nil /* use the zero value of the element type */)
		copy(slice[index+1:], slice[index:])
		slice[index] = value
		return slice
	}

	if len(lineMods) <= 0 {
		return
	}
	for _, y := range lineMods {
		// since test-only deps are in their own require blocks,
		// drop them from the main one.
		_ = f.DropRequire(y.name)
	}

	testLines := []*modfile.Line{
		{Token: []string{"// created using: https://github.com/komuw/ote \n"}},
	}
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
	index := findLastRequire(f) + 1
	f.Syntax.Stmt = insertAt(f.Syntax.Stmt, index, newTestBlock)
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
