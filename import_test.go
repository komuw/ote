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
