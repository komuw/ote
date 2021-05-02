package main

import (
	"testing"

	qt "github.com/frankban/quicktest"
)

func Test_difference(t *testing.T) {

	tests := []struct {
		name            string
		testModules     []string
		nonTestModules  []string
		trueTestModules []string
	}{
		{
			name:            "empty slices",
			testModules:     []string{},
			nonTestModules:  []string{},
			trueTestModules: []string{},
		},
		{
			name:            "some slices",
			testModules:     []string{"github.com/frankban/quicktest", "github.com/shirou/gopsutil", "rsc.io/quote"},
			nonTestModules:  []string{"github.com/hashicorp/nomad", "github.com/pkg/errors", "rsc.io/quote", "testdata/mod1", "github.com/LK4D4/joincontext"},
			trueTestModules: []string{"github.com/frankban/quicktest", "github.com/shirou/gopsutil"},
		},
		{
			name:            "duplicate items I",
			testModules:     []string{"github.com/frankban/quicktest", "rsc.io/quote", "github.com/shirou/gopsutil", "rsc.io/quote"},
			nonTestModules:  []string{"github.com/hashicorp/nomad", "github.com/pkg/errors", "rsc.io/quote", "testdata/mod1", "github.com/LK4D4/joincontext"},
			trueTestModules: []string{"github.com/frankban/quicktest", "github.com/shirou/gopsutil"},
		},
		{
			name:            "duplicate items II",
			testModules:     []string{"github.com/frankban/quicktest", "rsc.io/quote", "github.com/shirou/gopsutil", "rsc.io/quote"},
			nonTestModules:  []string{"rsc.io/quote", "github.com/hashicorp/nomad", "github.com/pkg/errors", "rsc.io/quote", "testdata/mod1", "github.com/LK4D4/joincontext"},
			trueTestModules: []string{"github.com/frankban/quicktest", "github.com/shirou/gopsutil"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := qt.New(t)
			got := difference(tt.testModules, tt.nonTestModules)
			c.Assert(got, qt.DeepEquals, tt.trueTestModules)
		})
	}
}

func Test_dedupe(t *testing.T) {
	tests := []struct {
		name string
		in   []string
		want []string
	}{
		{
			name: "empty slice",
			in:   []string{},
			want: []string{},
		},
		{
			name: "nil slice",
			in:   nil,
			want: nil,
		},
		{
			name: "small slice",
			in:   []string{"a", "a"},
			want: []string{"a"},
		},
		{
			name: "large slice",
			in:   []string{"a", "a", "b", "a"},
			want: []string{"a", "b"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := qt.New(t)
			got := dedupe(tt.in)
			c.Assert(got, qt.DeepEquals, tt.want)
		})
	}
}
