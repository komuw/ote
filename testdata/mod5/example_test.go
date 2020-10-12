package main_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/komuw/kama"
)

func TestSomething(t *testing.T) {
	res := kama.Dir(v.variable)

	if !cmp.Equal(2, 2) {
		t.Errorf("\ngot \n\t%#+v \nwanted \n\t%#+v", res, v.expected)
	}

}
