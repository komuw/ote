package main

import (
	"testing"

	qt "github.com/frankban/quicktest"
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
		{
			name:      "mod2",
			gomodFile: "testdata/mod2/go.mod",
			want:      "testdata/mod2",
		},
		{
			name:      "mod3",
			gomodFile: "testdata/mod3/go.mod",
			want:      "testdata/mod3",
		},
		{
			name:      "mod5",
			gomodFile: "testdata/mod5/go.mod",
			want:      "testdata/mod5",
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
