package main

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
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

	mP := "my/mod1"
	if !cmp.Equal(f.Module.Mod.Path, mP) {
		t.Errorf("\ngot \n\t%#+v \nwanted \n\t%#+v", f.Module.Mod.Path, mP)
	}

}

func TestGetPackage(t *testing.T) {
	mP := "my/mod1"
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

	mP := "my/mod1"
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
	thisMod := "my/mod1"
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
	modContents, err := ioutil.ReadFile(gomodFile)
	if err != nil {
		t.Fatal(err)
	}

	if strings.Contains(string(modContents), testComment) {
		t.Errorf("%#+v contained %#+v which is unexpected", gomodFile, testComment)
	}

	fp := "testdata/mod1"
	readonly := true
	buf := new(bytes.Buffer)
	run(fp, buf, readonly)

	if !strings.Contains(buf.String(), testComment) {
		t.Errorf("re-rendered %#+v did NOT contain %#+v which is unexpected", gomodFile, testComment)
	}

}
