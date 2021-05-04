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
	fmod1, _ := getModFile("testdata/mod1/go.mod")
	t.Cleanup(func() {
		fmod1.Cleanup()
	})

	tests := []struct {
		name            string
		trueTestModules []string
		f               *modfile.File
	}{
		{

			name:            "mod1",
			trueTestModules: []string{"github.com/frankban/quicktest", "github.com/shirou/gopsutil"},
			f:               fmod1,
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
	fmod1, _ := getModFile("testdata/mod1/go.mod")
	t.Cleanup(func() {
		fmod1.Cleanup()
	})
	fmod4, _ := getModFile("testdata/mod4/go.mod")
	t.Cleanup(func() {
		fmod4.Cleanup()
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
			f:               fmod1,
			gomodFile:       "testdata/mod1/go.mod",
			readonly:        true,
			want: []string{
				"module testdata/mod1",
				"github.com/frankban/quicktest v1.12.1 // test",
				"github.com/shirou/gopsutil v2.20.9+incompatible // test",
			},
		},
		{

			name:            "mod4",
			trueTestModules: []string{"github.com/benweissmann/memongo"},
			f:               fmod4,
			gomodFile:       "testdata/mod4/go.mod",
			readonly:        true,
			want: []string{
				"module testdata/mod4",
				"github.com/alexedwards/scs/v2 v2.4.0",
				"github.com/aws/aws-sdk-go v1.38.31",
				"github.com/benweissmann/memongo v0.1.1 // test",
				"github.com/go-kit/kit v0.10.0",
				"github.com/ishidawataru/sctp v0.0.0-20210226210310-f2269e66cdee",
				"github.com/rs/zerolog v1.21.0",
				"github.com/sirupsen/logrus v1.8.1",
				"github.com/zeebo/errs/v2 v2.0.3",
				"go.uber.org/zap v1.13.0",
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
