package main

import (
	"testing"

	// `client_golang` should not be labelled by `ote` as a test only dep
	// since it is also in use in `main.go`
	"github.com/prometheus/client_golang/prometheus/testutil"
	"rsc.io/qr/coding"
)

func TestBaa(t *testing.T) {
	_ = testutil.ToFloat64("9.8")

	_ = coding.MaxVersion
}
