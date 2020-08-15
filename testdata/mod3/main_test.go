package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
)

func TestMod3(t *testing.T) {
	x := "hello"
	y := "h" + "e" + "llo"
	if !cmp.Equal(x, y) {
		t.Errorf("mod3: x is not equal to y")
	}
}

func TestSomething(t *testing.T) {
	assert.Equal(t, 123, 123, "mod3: they should be equal")
}
