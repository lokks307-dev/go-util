package mt

import (
	"fmt"
	"testing"
)

func TestSet(t *testing.T) {
	var b *int
	c := 1
	b = &c
	a := NewIntSet()
	a.Add(b, 2)

	fmt.Println(a.ToSlice())
}
