package main

import (
	"testing"

	"go.uber.org/goleak"
)

func TestBaa(t *testing.T) {
	_ = goleak.Find()
}
