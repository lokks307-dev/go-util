package mt

import (
	"fmt"
	"testing"
)

func TestSumInt(t *testing.T) {

	fmt.Println(SumInt64(int64(1), float64(2.0), float32(4.5), float32(3.5)))
}
