package main

import (
	"bytes"
	"testing"

	qt "github.com/frankban/quicktest"
)

func Test_run(t *testing.T) {
	tests := []struct {
		name     string
		fp       string
		readonly bool
		want     string
	}{
		{
			name:     "ote's own modfile",
			fp:       ".",
			readonly: true,
			want: `module github.com/komuw/ote

go 1.16

require (
	github.com/frankban/quicktest v1.12.1 // test
	golang.org/x/mod v0.4.2
	golang.org/x/sys v0.0.0-20210426230700-d19ff857e887 // indirect
	golang.org/x/tools v0.1.0
)

`,
		},
		{
			name:     "testdata/mod1",
			fp:       "testdata/mod1",
			readonly: true,
			want: `module testdata/mod1

go 1.16

replace github.com/shirou/gopsutil => github.com/hashicorp/gopsutil v2.18.13-0.20200531184148-5aca383d4f9d+incompatible

require (
	github.com/LK4D4/joincontext v0.0.0-20171026170139-1724345da6d5
	github.com/StackExchange/wmi v0.0.0-20210224194228-fe8f1750fd46 // indirect
	github.com/frankban/quicktest v1.12.1 // test
	github.com/go-ole/go-ole v1.2.5 // indirect
	github.com/hashicorp/nomad v1.0.4 // The test comment should be removed by ote
	github.com/pkg/errors v0.9.1
	github.com/shirou/gopsutil v2.20.9+incompatible // test; PriorComment
	github.com/stretchr/testify v1.7.0 // indirect
	golang.org/x/net v0.0.0-20210501142056-aec3718b3fa0 // indirect
	golang.org/x/sys v0.0.0-20210426230700-d19ff857e887
	rsc.io/quote v1.5.2
)

`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := qt.New(t)

			w := &bytes.Buffer{}
			err := run(tt.fp, w, tt.readonly)
			c.Assert(err, qt.IsNil)

			got := w.String()
			c.Assert(got, qt.Equals, tt.want)
		})
	}
}
