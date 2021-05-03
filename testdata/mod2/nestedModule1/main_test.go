package main

import (
	"testing"

	// this is only in the nested module not in the main module.
	"crawshaw.io/sqlite"
)

func TestBaa(t *testing.T) {
	_ = sqlite.SQLITE_DENY
}
