// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"strings"

	"golang.org/x/mod/modfile"
)

// NB: most of the code in this file is taken from https://github.com/golang/mod/blob/v0.3.0/modfile/rule.go
// The license of that repo is found at: https://github.com/golang/mod/blob/v0.3.0/LICENSE
// and is included below:

//
//
//
// Copyright (c) 2009 The Go Authors. All rights reserved.

// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:

//    * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//    * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//    * Neither the name of Google Inc. nor the names of its
// contributors may be used to endorse or promote products derived from
// this software without specific prior written permission.

// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
//
//

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

// setTest sets line to have(or not have) a "// test" comment.
func setTest(line *modfile.Line, add bool) {
	if isTest(line) == add {
		return
	}

	if add {
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
	} else {
		// Removing comment.
		f := strings.Fields(line.Suffix[0].Token)
		if len(f) == 2 {
			// Remove whole comment.
			line.Suffix = nil
			return
		}

		// Remove comment prefix.
		com := &line.Suffix[0]
		i := strings.Index(com.Token, "test;")
		com.Token = "//" + com.Token[i+len("test;"):]
	}
}
