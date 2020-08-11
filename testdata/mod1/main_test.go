package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
)

func TestSomething(t *testing.T) {

	assert.Equal(t, 123, 123, "they should be equal")

	if !cmp.Equal(2, 2) {
		t.Errorf("mod1 test failure")
	}

}
