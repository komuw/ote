package main

import (
	"testing"

	qt "github.com/frankban/quicktest"
)

func Test_fetchImports(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		file string
		want []string
	}{
		{
			name: "nonTestFile",
			file: "testdata/mod1/main.go",
			want: []string{"testdata/mod1/api", "testdata/mod1/version"},
		},
		{
			name: "testFile",
			file: "testdata/mod1/version/ver_test.go",
			want: []string{"github.com/frankban/quicktest", "github.com/shirou/gopsutil/mem", "rsc.io/quote"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := qt.New(t)

			got, err := fetchImports(tt.file)
			c.Assert(err, qt.IsNil)
			c.Assert(got, qt.DeepEquals, tt.want)
		})
	}
}

func Test_fetchModule(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		root       string
		importPath string
		want       string
	}{
		{
			name:       "root with ending slash",
			root:       "testdata/mod1/",
			importPath: "github.com/hashicorp/nomad/drivers/shared/executor",
			want:       "github.com/hashicorp/nomad",
		},
		{
			name:       "root with NO ending slash",
			root:       "testdata/mod1",
			importPath: "github.com/hashicorp/nomad/drivers/shared/executor",
			want:       "github.com/hashicorp/nomad",
		},
		{
			name:       "short import path",
			root:       "testdata/mod1/",
			importPath: "rsc.io/quote",
			want:       "rsc.io/quote",
		},
		{
			name:       "main module that has nested module inside",
			root:       "testdata/mod2",
			importPath: "github.com/hashicorp/vault/api",
			want:       "github.com/hashicorp/vault/api",
		},
		{
			name:       "module that is nested module inside another main",
			root:       "testdata/mod2/nestedModule1",
			importPath: "crawshaw.io/sqlite",
			want:       "crawshaw.io/sqlite",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := qt.New(t)

			got, err := fetchModule(tt.root, tt.importPath)
			c.Assert(err, qt.IsNil)
			c.Assert(got, qt.Equals, tt.want)
		})
	}
}

func Test_getAllmodules(t *testing.T) {
	t.Parallel()

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
			nonTestImportPaths: []string{"github.com/hashicorp/nomad/drivers/shared/executor", "github.com/pkg/errors", "golang.org/x/sys/windows", "rsc.io/quote", "testdata/mod1/api", "testdata/mod1/version", "github.com/LK4D4/joincontext"},
			root:               "testdata/mod1/",
			wantTestModules:    []string{"github.com/frankban/quicktest", "github.com/shirou/gopsutil"},
			wantNonTestModules: []string{"github.com/hashicorp/nomad", "github.com/pkg/errors", "golang.org/x/sys", "testdata/mod1", "github.com/LK4D4/joincontext"},
		},
		{
			name:               "root with NO ending slash",
			testImportPaths:    []string{"github.com/frankban/quicktest", "github.com/shirou/gopsutil/mem", "rsc.io/quote"},
			nonTestImportPaths: []string{"github.com/hashicorp/nomad/drivers/shared/executor", "github.com/pkg/errors", "golang.org/x/sys/windows", "rsc.io/quote", "testdata/mod1/api", "testdata/mod1/version", "github.com/LK4D4/joincontext"},
			root:               "testdata/mod1",
			wantTestModules:    []string{"github.com/frankban/quicktest", "github.com/shirou/gopsutil"},
			wantNonTestModules: []string{"github.com/hashicorp/nomad", "github.com/pkg/errors", "golang.org/x/sys", "testdata/mod1", "github.com/LK4D4/joincontext"},
		},

		{
			name:               "with duplicates",
			testImportPaths:    []string{"github.com/frankban/quicktest", "github.com/frankban/quicktest", "github.com/shirou/gopsutil/mem", "github.com/frankban/quicktest"},
			nonTestImportPaths: []string{"github.com/hashicorp/nomad/drivers/shared/executor", "github.com/pkg/errors", "golang.org/x/sys/windows", "rsc.io/quote", "testdata/mod1/api", "testdata/mod1/version", "github.com/LK4D4/joincontext"},
			root:               "testdata/mod1",
			wantTestModules:    []string{"github.com/frankban/quicktest", "github.com/shirou/gopsutil"},
			wantNonTestModules: []string{"github.com/hashicorp/nomad", "github.com/pkg/errors", "golang.org/x/sys", "rsc.io/quote", "testdata/mod1", "github.com/LK4D4/joincontext"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := qt.New(t)

			gotTestModules, gotNonTestModules, err := getAllmodules(tt.testImportPaths, tt.nonTestImportPaths, tt.root)
			c.Assert(err, qt.IsNil)
			c.Assert(gotTestModules, qt.DeepEquals, tt.wantTestModules)
			c.Assert(gotNonTestModules, qt.DeepEquals, tt.wantNonTestModules)
		})
	}
}

func Test_getTestModules(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		root string
		want []string
	}{
		{
			name: "root with ending slash",
			root: "testdata/mod1/",
			want: []string{"github.com/frankban/quicktest", "github.com/shirou/gopsutil"},
		},
		{
			name: "root with NO ending slash",
			root: "testdata/mod1",
			want: []string{"github.com/frankban/quicktest", "github.com/shirou/gopsutil"},
		},
		{
			name: "main module that has nested module inside",
			root: "testdata/mod2",
			want: []string{"github.com/shirou/gopsutil", "gopkg.in/natefinch/lumberjack.v2"},
		},
		{
			name: "module that is nested module inside another main",
			root: "testdata/mod2/nestedModule1",
			want: []string{"crawshaw.io/sqlite"},
		},
		{
			name: "module that is nested module inside another main take II",
			root: "testdata/mod2/nestedModule2",
			want: []string{"github.com/sirupsen/logrus"},
		},
		{
			name: "module with vendor directory",
			root: "testdata/mod3",
			want: []string{"go.uber.org/goleak"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := qt.New(t)

			got, err := getTestModules(tt.root)
			c.Assert(err, qt.IsNil)
			c.Assert(got, qt.DeepEquals, tt.want)
		})
	}
}

func Test_isStdLibPkg(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		pkg  string
		std  string
		want bool
	}{
		{
			name: "fmt ok",
			pkg:  "fmt",
			std:  stdlib,
			want: true,
		},
		{
			name: "bogus package",
			pkg:  "bogusPkg",
			std:  stdlib,
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := qt.New(t)

			got, err := isStdLibPkg(tt.pkg, tt.std)
			c.Assert(err, qt.IsNil)
			c.Assert(got, qt.Equals, tt.want)
		})
	}
}
