package main

import (
	"fmt"
	"testing"

	"github.com/shirou/gopsutil/mem"
)

func TestBaa(t *testing.T) {
	v, err := mem.VirtualMemory()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(v)
}
