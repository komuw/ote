package main

import (
	"testing"

	qt "github.com/frankban/quicktest"
)

func Test_loadStd(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "run", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := qt.New(t)
			err := loadStd()
			c.Assert(err, qt.IsNil)

			got := isStdLibPkg("archive/tar")
			c.Assert(got, qt.IsTrue)

			got2 := isStdLibPkg("nonExistent")
			c.Assert(got2, qt.IsFalse)
		})
	}
}

func Test_fetchImports(t *testing.T) {
	tests := []struct {
		name string
		file string
		want []string
	}{
		{
			name: "nonTestFile",
			file: "testdata/mod2/main.go",
			want: []string{"fmt", "testdata/mod2/api", "testdata/mod2/version"},
		},
		{
			name: "testFile",
			file: "testdata/mod2/version/ver_test.go",
			want: []string{"fmt", "testing", "github.com/frankban/quicktest", "github.com/shirou/gopsutil/mem", "rsc.io/quote"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := qt.New(t)

			got, err := fetchImports(tt.file)
			t.Log(got)
			c.Assert(err, qt.IsNil)
			c.Assert(got, qt.DeepEquals, tt.want)
		})
	}
}
