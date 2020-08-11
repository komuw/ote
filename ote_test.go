package main

import (
	"sort"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestGetModFile(t *testing.T) {
	f, err := getModFile()
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
	pkg, err := getPackage(mP, true)

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
	m, err := getModules(mP)
	if err != nil {
		t.Fatal(err)
	}

	if !cmp.Equal(m, expectedModules, trans) {
		t.Errorf("\ngot \n\t%#+v \nwanted \n\t%#+v", m, expectedModules)
	}

}

func TestGetDeps(t *testing.T) {
	thisMod := "github.com/komuw/ote"
	_, err := getDeps(thisMod)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetTestDeps(t *testing.T) {
	thisMod := "github.com/komuw/ote"

	modulePaths, err := getModules(thisMod)
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

	p := "testdata/mod1"
	_ = p

}
