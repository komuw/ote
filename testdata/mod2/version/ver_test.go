package version

import (
	"fmt"
	"testing"

	qt "github.com/frankban/quicktest"

	"github.com/shirou/gopsutil/mem"
)

func TestFoo(t *testing.T) {
	t.Run("numbers", func(t *testing.T) {
		c := qt.New(t)
		numbers := []int{42, 47}
		c.Assert(numbers, qt.DeepEquals, []int{42, 47})
	})

}

func TestBaa(t *testing.T) {

	v, err := mem.VirtualMemory()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(v)

}
