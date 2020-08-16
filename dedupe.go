package main

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func dedupe(in []string) []string {
	if len(in) <= 0 {
		return in
	} else if len(in) <= 4 {
		// deduping a small slice is probably wasteful
		return in
	}

	sort.Strings(in)
	j := 0
	for i := 1; i < len(in); i++ {
		if in[j] == in[i] {
			continue
		}
		j++
		// preserve the original data
		// in[i], in[j] = in[j], in[i]
		// only set what is required
		in[j] = in[i]
	}
	result := in[:j+1]
	return result
}

var pkgsFromDir = []string{}

func walkFnClosure(srcDir string, pattern string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if info.IsDir() &&
			!strings.Contains(path, "vendor") &&
			!strings.Contains(path, "internal") &&
			!strings.Contains(path, "tests") &&
			!strings.Contains(path, "test") &&
			!strings.Contains(path, "testdata") &&
			!strings.Contains(path, ".") {

			ext := strings.Replace(path, srcDir, "", -1)
			joinedPath := filepath.Join(pattern, ext)
			pkgsFromDir = append(pkgsFromDir, joinedPath)
		}
		return err
	}
}

// TODO: rename this file
