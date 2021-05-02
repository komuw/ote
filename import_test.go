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
			c.Assert(err, qt.IsNil)
			c.Assert(got, qt.DeepEquals, tt.want)
		})
	}
}

func Test_fetchModule(t *testing.T) {
	tests := []struct {
		name       string
		root       string
		importPath string
		want       string
	}{
		{
			name:       "root with ending slash",
			root:       "testdata/mod2/",
			importPath: "github.com/hashicorp/nomad/drivers/shared/executor",
			want:       "github.com/hashicorp/nomad",
		},
		{
			name:       "root with NO ending slash",
			root:       "testdata/mod2",
			importPath: "github.com/hashicorp/nomad/drivers/shared/executor",
			want:       "github.com/hashicorp/nomad",
		},
		{
			name:       "short import path",
			root:       "testdata/mod2/",
			importPath: "rsc.io/quote",
			want:       "rsc.io/quote",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := qt.New(t)

			got, err := fetchModule(tt.root, tt.importPath)
			c.Assert(err, qt.IsNil)
			c.Assert(got, qt.Equals, tt.want)
		})
	}
}

func Test_getAllmodules(t *testing.T) {
	tests := []struct {
		name               string
		testImportPaths    []string
		nonTestImportPaths []string
		root               string
		wantTestModules    []string
		wantNonTestModules []string
	}{
		{
			name:               "root with ending slash",
			testImportPaths:    []string{"github.com/frankban/quicktest", "github.com/shirou/gopsutil/mem", "rsc.io/quote"},
			nonTestImportPaths: []string{"github.com/hashicorp/nomad/drivers/shared/executor", "github.com/pkg/errors", "golang.org/x/sys/windows", "rsc.io/quote", "testdata/mod2/api", "testdata/mod2/version", "github.com/LK4D4/joincontext"},
			root:               "testdata/mod2/",
			wantTestModules:    []string{"github.com/frankban/quicktest", "github.com/shirou/gopsutil", "rsc.io/quote"},
			wantNonTestModules: []string{"github.com/hashicorp/nomad", "github.com/pkg/errors", "golang.org/x/sys", "rsc.io/quote", "testdata/mod2", "github.com/LK4D4/joincontext"},
		},
		{
			name:               "root with NO ending slash",
			testImportPaths:    []string{"github.com/frankban/quicktest", "github.com/shirou/gopsutil/mem", "rsc.io/quote"},
			nonTestImportPaths: []string{"github.com/hashicorp/nomad/drivers/shared/executor", "github.com/pkg/errors", "golang.org/x/sys/windows", "rsc.io/quote", "testdata/mod2/api", "testdata/mod2/version", "github.com/LK4D4/joincontext"},
			root:               "testdata/mod2",
			wantTestModules:    []string{"github.com/frankban/quicktest", "github.com/shirou/gopsutil", "rsc.io/quote"},
			wantNonTestModules: []string{"github.com/hashicorp/nomad", "github.com/pkg/errors", "golang.org/x/sys", "rsc.io/quote", "testdata/mod2", "github.com/LK4D4/joincontext"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := qt.New(t)

			gotTestModules, gotNonTestModules, err := getAllmodules(tt.testImportPaths, tt.nonTestImportPaths, tt.root)
			c.Assert(err, qt.IsNil)
			c.Assert(gotTestModules, qt.DeepEquals, tt.wantTestModules)
			c.Assert(gotNonTestModules, qt.DeepEquals, tt.wantNonTestModules)
		})
	}
}
