package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestMod3(t *testing.T) {
	x := "hello"
	y := "h" + "e" + "llo"
	if !cmp.Equal(x, y) {
		t.Errorf("mod3: x is not equal to y")
	}

}
