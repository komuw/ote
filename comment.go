// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"strings"

	"golang.org/x/mod/modfile"
)

// NB: most of the code here is taken from https://github.com/golang/mod/blob/v0.3.0/modfile/rule.go

// TODO: add license. see: https://github.com/rogpeppe/go-internal/pull/37/files

var (
	slashSlash  = []byte("//")
	testComment = "// test"
)

// isTest reports whether line has a "// test" comment
func isTest(line *modfile.Line) bool {
	if len(line.Suffix) == 0 {
		return false
	}
	f := strings.Fields(strings.TrimPrefix(line.Suffix[0].Token, string(slashSlash)))
	return (len(f) == 1 && f[0] == "test" || len(f) > 1 && f[0] == "test;")
}

// setTest sets line to have a "// test" comment.
func setTest(line *modfile.Line) {
	if isTest(line) {
		return
	}

	// Adding comment.
	if len(line.Suffix) == 0 {
		// New comment.
		line.Suffix = []modfile.Comment{{Token: testComment, Suffix: true}}
		return
	}

	com := &line.Suffix[0]
	text := strings.TrimSpace(strings.TrimPrefix(com.Token, string(slashSlash)))
	if text == "" {
		// Empty comment.
		com.Token = testComment
		return
	}

	// Insert at beginning of existing comment.
	com.Token = "// test; " + text
	return

}
