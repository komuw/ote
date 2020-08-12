package version

import (
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestFoo(t *testing.T) {
	t.Run("numbers", func(t *testing.T) {
		c := qt.New(t)
		numbers := []int{42, 47}
		c.Assert(numbers, qt.DeepEquals, []int{42, 47})
	})

}
