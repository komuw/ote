package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/kylelemons/godebug/diff"
)

var (
	fp        = "testdata/mod1"
	gomodFile = filepath.Join(fp, "go.mod")
)

func TestGetModFile(t *testing.T) {
	f, err := getModFile(gomodFile)
	if err != nil {
		t.Fatal(err)
	}

	mP := "testdata/mod1"
	if !cmp.Equal(f.Module.Mod.Path, mP) {
		t.Errorf("\ngot \n\t%#+v \nwanted \n\t%#+v", f.Module.Mod.Path, mP)
	}

}

func TestGetPackage(t *testing.T) {
	mP := "testdata/mod1"
	pkg, err := getPackage(mP, gomodFile, true)

	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(pkg.Module.Path, mP) {
		t.Errorf("\ngot \n\t%#+v \nwanted \n\t%#+v", pkg.Module.Path, mP)
	}

}

func TestGetModules(t *testing.T) {
	trans := cmp.Transformer("Sort", func(in []string) []string {
		out := append([]string(nil), in...)
		sort.Strings(out)
		return out
	})

	expectedModules := []string{"github.com/Shopify/sarama", "github.com/nats-io/nats.go"}

	mP := "testdata/mod1"
	m, err := getModules(mP, gomodFile)
	if err != nil {
		t.Fatal(err)
	}

	if !cmp.Equal(m, expectedModules, trans) {
		t.Errorf("\ngot \n\t%#+v \nwanted \n\t%#+v", m, expectedModules)
	}

}

func TestGetDeps(t *testing.T) {
	_, err := getDeps(gomodFile)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetTestDeps(t *testing.T) {
	thisMod := "testdata/mod1"
	modulePaths, err := getModules(thisMod, gomodFile)
	if err != nil {
		t.Fatal(err)
	}

	allDeps, err := getDeps(gomodFile)
	if err != nil {
		t.Fatal(err)
	}

	_ = getTestDeps(modulePaths, allDeps)
}

func TestCli(t *testing.T) {
	f, r := cli()

	if !cmp.Equal(f, ".") {
		t.Errorf("\ngot \n\t%#+v \nwanted \n\t%#+v", f, ".")
	}

	if !cmp.Equal(r, false) {
		t.Errorf("\ngot \n\t%#+v \nwanted \n\t%#+v", r, false)
	}
}

func TestRun(t *testing.T) {

	tt := []struct {
		fp                      string
		modFilePath             string
		expectedModfile         []byte
		expectedModifiedModfile []byte
	}{
		{
			fp:          ".",
			modFilePath: "go.mod",
			expectedModfile: []byte(`module github.com/komuw/ote

go 1.14

require (
	github.com/google/go-cmp v0.5.2 // test
	github.com/kylelemons/godebug v1.1.0 // test
	golang.org/x/mod v0.3.0
	golang.org/x/tools v0.0.0-20201012192620-5bd05386311b
)
`),
			expectedModifiedModfile: []byte(`module github.com/komuw/ote

go 1.14

require (
	github.com/google/go-cmp v0.5.2 // test
	github.com/kylelemons/godebug v1.1.0 // test
	golang.org/x/mod v0.3.0
	golang.org/x/tools v0.0.0-20201012192620-5bd05386311b
)

`),
		},

		{
			fp:          "testdata/mod1",
			modFilePath: "testdata/mod1/go.mod",
			expectedModfile: []byte(`module testdata/mod1

go 1.14

require (
	github.com/Shopify/sarama v1.27.0
	github.com/golang/protobuf v1.4.2 // indirect
	github.com/google/go-cmp v0.5.1
	github.com/nats-io/nats-server/v2 v2.1.7 // indirect
	github.com/nats-io/nats.go v1.10.0
	github.com/stretchr/testify v1.6.1 // SomePriorComment
	golang.org/x/net v0.0.0-20200625001655-4c5254603344 // indirect
	google.golang.org/protobuf v1.25.0 // indirect
)
`),
			expectedModifiedModfile: []byte(`module testdata/mod1

go 1.14

require (
	github.com/Shopify/sarama v1.27.0
	github.com/golang/protobuf v1.4.2 // indirect
	github.com/google/go-cmp v0.5.1 // test
	github.com/nats-io/nats-server/v2 v2.1.7 // indirect
	github.com/nats-io/nats.go v1.10.0
	github.com/stretchr/testify v1.6.1 // test; SomePriorComment
	golang.org/x/net v0.0.0-20200625001655-4c5254603344 // indirect
	google.golang.org/protobuf v1.25.0 // indirect
)

`),
		},

		{
			fp:          "testdata/mod3",
			modFilePath: "testdata/mod3/go.mod",
			expectedModfile: []byte(`module testdata/mod3

go 1.14

// ote should remove the //test comment from go-cmp since it is also used in main.go
// it should also add a //test comment to testify
require (
	github.com/google/go-cmp v0.5.1 // test
	github.com/pkg/json v0.0.0-20200630040052-6ff993914616
	github.com/stretchr/testify v1.3.0
)
`),
			expectedModifiedModfile: []byte(`module testdata/mod3

go 1.14

// ote should remove the //test comment from go-cmp since it is also used in main.go
// it should also add a //test comment to testify
require (
	github.com/google/go-cmp v0.5.1
	github.com/pkg/json v0.0.0-20200630040052-6ff993914616
	github.com/stretchr/testify v1.3.0 // test
)

`),
		},

		{
			fp:          "testdata/mod5",
			modFilePath: "testdata/mod5/go.mod",
			expectedModfile: []byte(`module testdata/mod5

go 1.15

require (
	github.com/golang/protobuf v1.4.2 // indirect
	github.com/google/go-cmp v0.5.2
	github.com/komuw/kama v0.0.0-20201012123531-9c57efc1ae36
	github.com/nats-io/nats-server/v2 v2.1.8 // indirect
	github.com/nats-io/nats.go v1.10.0
	golang.org/x/crypto v0.0.0-20200820211705-5c72a883971a // indirect
	google.golang.org/protobuf v1.25.0 // indirect
)
`),
			expectedModifiedModfile: []byte(`module testdata/mod5

go 1.15

require (
	github.com/golang/protobuf v1.4.2 // indirect
	github.com/google/go-cmp v0.5.2 // test
	github.com/komuw/kama v0.0.0-20201012123531-9c57efc1ae36 // test
	github.com/nats-io/nats-server/v2 v2.1.8 // indirect
	github.com/nats-io/nats.go v1.10.0
	golang.org/x/crypto v0.0.0-20200820211705-5c72a883971a // indirect
	google.golang.org/protobuf v1.25.0 // indirect
)

`),
		},
	}

	for _, v := range tt {
		v := v // capture range variable
		t.Run(fmt.Sprintf("runing test for modFilePath: %s", v.modFilePath), func(t *testing.T) {
			originalMod, err := ioutil.ReadFile(v.modFilePath)
			if err != nil {
				t.Fatal(err)
			}

			// NB: you can use `github.com/kylelemons/godebug/diff` to find out how to make two strings equal
			// diff := diff.Diff(string(v.expectedModfile), string(originalMod))
			if !cmp.Equal(originalMod, v.expectedModfile) {
				diff := diff.Diff(string(v.expectedModfile), string(originalMod))
				t.Errorf("\n original modfile mismatch, diff: \n======================\n%s\n======================\n", diff)
			}

			readonly := true
			oteMod := new(bytes.Buffer)
			run(v.fp, oteMod, readonly)

			if !cmp.Equal(oteMod.Bytes(), v.expectedModifiedModfile) {
				diff := diff.Diff(string(v.expectedModifiedModfile), string(oteMod.Bytes()))
				t.Errorf("\n modified modfile mismatch, diff: \n======================\n%s\n======================\n", diff)
			}
		})
	}
}
