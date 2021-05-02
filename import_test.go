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
			c := qt.New(t)

			//  load std libs
			loadStd()

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
			nonTestImportPaths: []string{"github.com/hashicorp/nomad/drivers/shared/executor", "github.com/pkg/errors", "golang.org/x/sys/windows", "rsc.io/quote", "testdata/mod1/api", "testdata/mod1/version", "github.com/LK4D4/joincontext"},
			root:               "testdata/mod1/",
			wantTestModules:    []string{"github.com/frankban/quicktest", "github.com/shirou/gopsutil", "rsc.io/quote"},
			wantNonTestModules: []string{"github.com/hashicorp/nomad", "github.com/pkg/errors", "golang.org/x/sys", "rsc.io/quote", "testdata/mod1", "github.com/LK4D4/joincontext"},
		},
		{
			name:               "root with NO ending slash",
			testImportPaths:    []string{"github.com/frankban/quicktest", "github.com/shirou/gopsutil/mem", "rsc.io/quote"},
			nonTestImportPaths: []string{"github.com/hashicorp/nomad/drivers/shared/executor", "github.com/pkg/errors", "golang.org/x/sys/windows", "rsc.io/quote", "testdata/mod1/api", "testdata/mod1/version", "github.com/LK4D4/joincontext"},
			root:               "testdata/mod1",
			wantTestModules:    []string{"github.com/frankban/quicktest", "github.com/shirou/gopsutil", "rsc.io/quote"},
			wantNonTestModules: []string{"github.com/hashicorp/nomad", "github.com/pkg/errors", "golang.org/x/sys", "rsc.io/quote", "testdata/mod1", "github.com/LK4D4/joincontext"},
		},

		{
			name:               "with duplicates",
			testImportPaths:    []string{"rsc.io/quote", "github.com/frankban/quicktest", "rsc.io/quote", "github.com/shirou/gopsutil/mem", "rsc.io/quote"},
			nonTestImportPaths: []string{"github.com/hashicorp/nomad/drivers/shared/executor", "github.com/pkg/errors", "golang.org/x/sys/windows", "rsc.io/quote", "testdata/mod1/api", "testdata/mod1/version", "github.com/LK4D4/joincontext"},
			root:               "testdata/mod1",
			wantTestModules:    []string{"rsc.io/quote", "github.com/frankban/quicktest", "github.com/shirou/gopsutil"},
			wantNonTestModules: []string{"github.com/hashicorp/nomad", "github.com/pkg/errors", "golang.org/x/sys", "rsc.io/quote", "testdata/mod1", "github.com/LK4D4/joincontext"},
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

func Test_getTestModules(t *testing.T) {
	tests := []struct {
		name string
		root string
		want []string
	}{
		// TODO: undo the comment

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
			want: []string{"github.com/frankban/quicktest", "github.com/shirou/gopsutil"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := qt.New(t)

			//  load std libs
			loadStd()

			got, err := getTestModules(tt.root)
			c.Assert(err, qt.IsNil)
			c.Assert(got, qt.DeepEquals, tt.want)
		})
	}
}

// {
// 	name:       "main module that has nested module inside",
// 	root:       "testdata/mod2",
// 	importPath: "github.com/hashicorp/vault/api",
// 	want:       "github.com/hashicorp/vault/api",
// },
// {
// 	name:       "module that is nested module inside another main",
// 	root:       "testdata/mod2/nestedModule1",
// 	importPath: "crawshaw.io/sqlite",
// 	want:       "crawshaw.io/sqlite",
// },
