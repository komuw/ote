package main

import (
	"bytes"
	"testing"

	qt "github.com/frankban/quicktest"
)

func Test_run(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		fp       string
		readonly bool
		wantErr  string
		want     string
	}{

		{
			name:     "testdata/mod1",
			fp:       "testdata/mod1",
			readonly: true,
			wantErr:  "",
			want: `module testdata/mod1

go 1.16

replace github.com/shirou/gopsutil => github.com/hashicorp/gopsutil v2.18.13-0.20200531184148-5aca383d4f9d+incompatible

require (
	github.com/LK4D4/joincontext v0.0.0-20171026170139-1724345da6d5
	github.com/StackExchange/wmi v0.0.0-20210224194228-fe8f1750fd46 // indirect
	github.com/go-ole/go-ole v1.2.5 // indirect
	github.com/hashicorp/nomad v1.0.4 // The test comment should be removed by ote
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.7.0 // indirect
	golang.org/x/net v0.0.0-20210501142056-aec3718b3fa0 // indirect
	golang.org/x/sys v0.0.0-20210426230700-d19ff857e887
	rsc.io/quote v1.5.2
)

exclude golang.org/x/net v1.2.3

retract (
	v1.0.1 // Contains retractions only.
	v1.0.0 // Published accidentally.
)

require (
	github.com/frankban/quicktest v1.12.1 // test
	github.com/shirou/gopsutil v2.20.9+incompatible // test; PriorComment
)

`,
		},

		{
			name:     "testdata/mod2",
			fp:       "testdata/mod2",
			readonly: true,
			wantErr:  "",
			want: `module testdata/mod2

go 1.16

require (
	github.com/StackExchange/wmi v0.0.0-20210224194228-fe8f1750fd46 // indirect
	github.com/chromedp/chromedp v0.7.1
	github.com/fatih/color v1.10.0 // indirect
	github.com/frankban/quicktest v1.10.2 // indirect
	github.com/go-ole/go-ole v1.2.5 // indirect
	github.com/golang/snappy v0.0.3 // indirect
	github.com/google/go-cmp v0.5.5 // indirect
	github.com/hashicorp/go-hclog v0.14.1 // indirect
	github.com/hashicorp/vault/api v1.1.0
	github.com/kr/text v0.2.0 // indirect
	github.com/pierrec/lz4 v2.5.2+incompatible // indirect
	github.com/stretchr/testify v1.6.1 // indirect
	golang.org/x/crypto v0.0.0-20210415154028-4f45737414dc // indirect
	golang.org/x/net v0.0.0-20210316092652-d523dce5a7f4 // indirect
	golang.org/x/text v0.3.5 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776 // indirect
)

require (
	github.com/shirou/gopsutil v3.21.4+incompatible // test
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // test
)

`,
		},

		{
			name:     "testdata/mod3",
			fp:       "testdata/mod3",
			readonly: true,
			wantErr:  "",
			want: `module testdata/mod3

go 1.16

require (
	github.com/ethereum/go-ethereum v1.10.2
	github.com/kingzbauer/africastalking-go v0.0.2-alpha.1
)

require go.uber.org/goleak v1.1.10 // test

`,
		},

		{
			name:     "testdata/mod4",
			fp:       "testdata/mod4",
			readonly: true,
			wantErr:  "",
			want: `module testdata/mod4

go 1.16

require (
	github.com/alexedwards/scs/v2 v2.4.0
	github.com/aws/aws-sdk-go v1.38.31
	github.com/go-kit/kit v0.10.0
	github.com/ishidawataru/sctp v0.0.0-20210226210310-f2269e66cdee
	github.com/ory/herodot v0.9.5
	github.com/rs/zerolog v1.21.0
	github.com/sirupsen/logrus v1.8.1
	github.com/zeebo/errs/v2 v2.0.3
	go.uber.org/zap v1.13.0
)

require github.com/benweissmann/memongo v0.1.1 // test

`,
		},

		{
			// TODO: This module should actually not raise an error.
			// It should succeed.
			// This should be fixed in:
			// https://github.com/komuw/ote/issues/3
			// https://github.com/komuw/ote/issues/27
			// We should update this test when those issues are solved
			name:     "testdata/mod5",
			fp:       "testdata/mod5",
			readonly: true,
			wantErr:  "build constraints exclude all Go files in",
			want:     ``,
		},

		{
			name:     "testdata/nonExistentPackage",
			fp:       "testdata/nonExistentPackage",
			readonly: true,
			wantErr:  "no such file or directory",
			want:     ``,
		},
	}
	for _, tt := range tests {
		tt := tt // capture range variable

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := qt.New(t)

			w := &bytes.Buffer{}
			err := run(tt.fp, w, tt.readonly)
			if len(tt.wantErr) > 0 {
				c.Assert(err.Error(), qt.Contains, tt.wantErr)
			} else {
				c.Assert(err, qt.IsNil)
			}

			got := w.String()
			c.Assert(got, qt.Equals, tt.want)
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
