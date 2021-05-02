package main

import (
	"bytes"
	"testing"

	qt "github.com/frankban/quicktest"
	"golang.org/x/mod/modfile"
)

func Test_getModFile(t *testing.T) {
	tests := []struct {
		name      string
		gomodFile string
		want      string
	}{

		{
			name:      "mod1",
			gomodFile: "testdata/mod1/go.mod",
			want:      "testdata/mod1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := qt.New(t)

			got, err := getModFile(tt.gomodFile)
			c.Assert(err, qt.IsNil)
			c.Assert(got.Module.Mod.Path, qt.DeepEquals, tt.want)
		})
	}
}

func Test_updateMod(t *testing.T) {
	f, _ := getModFile("testdata/mod1/go.mod")
	t.Cleanup(func() {
		f.Cleanup()
	})

	tests := []struct {
		name            string
		trueTestModules []string
		f               *modfile.File
	}{
		{

			name:            "mod1",
			trueTestModules: []string{"github.com/frankban/quicktest", "github.com/shirou/gopsutil"},
			f:               f,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := qt.New(t)

			err := updateMod(tt.trueTestModules, tt.f)
			c.Assert(err, qt.IsNil)
		})
	}
}

func Test_writeMod(t *testing.T) {
	f, _ := getModFile("testdata/mod1/go.mod")
	t.Cleanup(func() {
		f.Cleanup()
	})

	tests := []struct {
		name            string
		trueTestModules []string
		f               *modfile.File
		gomodFile       string
		readonly        bool
		want            []string
	}{
		{

			name:            "mod1",
			trueTestModules: []string{"github.com/frankban/quicktest", "github.com/shirou/gopsutil"},
			f:               f,
			gomodFile:       "testdata/mod1/go.mod",
			readonly:        true,
			want: []string{
				"module testdata/mod1",
				"github.com/frankban/quicktest v1.12.1 // test",
				"github.com/shirou/gopsutil v2.20.9+incompatible // test",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := qt.New(t)

			err := updateMod(tt.trueTestModules, tt.f)
			c.Assert(err, qt.IsNil)

			w := &bytes.Buffer{}
			errW := writeMod(tt.f, tt.gomodFile, w, tt.readonly)
			c.Assert(errW, qt.IsNil)

			for _, v := range tt.want {
				c.Assert(w.String(), qt.Contains, v)
			}
		})
	}
}
