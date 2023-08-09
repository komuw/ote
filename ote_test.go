package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"
	"go.akshayshah.org/attest"
)

func Test_run(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		fp       string
		readonly bool
		wantErr  string
	}{
		{
			name:     "testdata/modfiles/mod1",
			fp:       "testdata/modfiles/mod1",
			readonly: true,
			wantErr:  "",
		},

		{
			name:     "testdata/modfiles/mod2",
			fp:       "testdata/modfiles/mod2",
			readonly: true,
			wantErr:  "",
		},

		{
			name:     "testdata/modfiles/mod3",
			fp:       "testdata/modfiles/mod3",
			readonly: true,
			wantErr:  "",
		},

		{
			name:     "testdata/modfiles/mod4",
			fp:       "testdata/modfiles/mod4",
			readonly: true,
			wantErr:  "",
		},

		{
			name:     "testdata/modfiles/mod5",
			fp:       "testdata/modfiles/mod5",
			readonly: true,
			wantErr:  "",
		},

		{
			name:     "testdata/modfiles/mod6",
			fp:       "testdata/modfiles/mod6",
			readonly: true,
			wantErr:  "",
		},

		{
			name:     "testdata/modfiles/mod7",
			fp:       "testdata/modfiles/mod7",
			readonly: true,
			wantErr:  "",
		},

		{
			name:     "testdata/modfiles/nonExistentPackage",
			fp:       "testdata/modfiles/nonExistentPackage",
			readonly: true,
			wantErr:  "no such file or directory",
		},
	}
	for _, tt := range tests {
		tt := tt // capture range variable

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			w := &bytes.Buffer{}
			err := run(tt.fp, w, tt.readonly)
			if len(tt.wantErr) > 0 {
				attest.Subsequence(t, err.Error(), tt.wantErr)
				return
			}

			attest.Ok(t, err)

			got := w.String()
			path := getDataPath(t, "ote_test.go", tt.name)
			dealWithTestData(t, path, got)
		})
	}
}

func Test_cli(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		file     string
		readonly bool
		version  bool
	}{
		{
			name:     "current directory",
			file:     ".",
			readonly: false,
			version:  false,
		},
	}
	for _, tt := range tests {
		tt := tt // capture range variable

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := qt.New(t)

			f, r, v := cli()
			c.Assert(f, qt.Equals, tt.file)
			c.Assert(r, qt.Equals, tt.readonly)
			c.Assert(v, qt.Equals, tt.version)
		})
	}
}

const oteWriteDataForTests = "OTE_WRITE_DATA_FOR_TESTS"

// dealWithTestData asserts that gotContent is equal to data found at path.
//
// If the environment variable [oteWriteDataForTests] is set, this func
// will write gotContent to path instead.
func dealWithTestData(t *testing.T, path, gotContent string) {
	t.Helper()

	path = strings.ReplaceAll(path, ".go", "")

	p, e := filepath.Abs(path)
	attest.Ok(t, e)

	writeData := os.Getenv(oteWriteDataForTests) != ""
	if writeData {
		attest.Ok(t,
			os.WriteFile(path, []byte(gotContent), 0o644),
		)
		t.Logf("\n\t written testdata to %s\n", path)
		return
	}

	b, e := os.ReadFile(p)
	attest.Ok(t, e)

	expectedContent := string(b)
	attest.Equal(t, gotContent, expectedContent, attest.Sprintf("path: %s", path))
}

func getDataPath(t *testing.T, testPath, testName string) string { //nolint:unparam
	t.Helper()

	s := strings.ReplaceAll(testName, " ", "_")
	tName := strings.ReplaceAll(s, "/", "_")

	path := filepath.Join("testdata", "text_files", testPath, tName) + ".txt"

	return path
}
