package main

import (
	"path/filepath"
	"sort"
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

	mP := "github.com/komuw/ote"
	if !cmp.Equal(f.Module.Mod.Path, mP) {
		t.Errorf("\ngot \n\t%#+v \nwanted \n\t%#+v", f.Module.Mod.Path, mP)
	}

}

func TestGetPackage(t *testing.T) {
	mP := "github.com/komuw/ote"
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

	expectedModules := []string{
		"golang.org/x/mod",
		"github.com/Shopify/sarama",
		"github.com/nats-io/nats.go",
		"golang.org/x/tools",
	}

	mP := "github.com/komuw/ote"
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
	thisMod := "github.com/komuw/ote"
	modulePaths, err := getModules(thisMod, gomodFile)
	if err != nil {
		t.Fatal(err)
	}

	allDeps, err := getDeps(thisMod)
	if err != nil {
		t.Fatal(err)
	}

	_ = getTestDeps(modulePaths, allDeps)
}

func TestRun(t *testing.T) {
	fp := "testdata/mod1"
	readonly := true
	run(fp, readonly)

}
