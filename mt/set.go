package mt

type IntSet struct {
	intMap map[int64]bool
}

func NewIntSet() *IntSet {
	return &IntSet{
		intMap: make(map[int64]bool),
	}
}

func (m *IntSet) Add(vv ...interface{}) {
	for _, v := range vv {
		vit, ok := AnyToInt64(v)
		if ok {
			m.intMap[vit] = true
		}
	}
}

func (m *IntSet) ToSlice() []int64 {
	ret := make([]int64, 0)
	for v := range m.intMap {
		ret = append(ret, v)
	}
	return ret
}
