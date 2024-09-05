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

func TestTemplateSet(t *testing.T) {

	a := NewSet[string]()
	a.Add("a", "b")
	fmt.Println(a.IsIn("a"))
	fmt.Println(a.ToSlice())
}
