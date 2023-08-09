//go:build !darwin && !linux && !windows
// +build !darwin,!linux,!windows

package main

import (
	// this import is shared in test files and non-test files
	// thus it should not be rendered by `ote` with a `// test` comment.
	// this position will only be held if `ote` is able to ignore the build tags in this file.
	// the tags make it so that this file can't be built in macOS,linux or windows(the places where we run `ote` tests/ci)
	"github.com/kr/pretty"
)

func MyPrint(x interface{}) {
	pretty.Println(x)
}
