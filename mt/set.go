package mt

func IsOptionSet(opt []bool) bool {
	return len(opt) > 0 && opt[0]
}

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

func (m *IntSet) IsIn(v int64) bool {
	if m.intMap == nil {
		return false
	}

	_, ok := m.intMap[v]
	return ok
}

func (m *IntSet) Size() int {
	return len(m.intMap)
}

func (m *IntSet) ToSlice() []int64 {
	ret := make([]int64, 0)
	for v := range m.intMap {
		ret = append(ret, v)
	}
	return ret
}

type Set[T comparable] struct {
	inMap map[T]bool
}

func NewSet[T comparable]() *Set[T] {
	return &Set[T]{
		inMap: make(map[T]bool),
	}
}

func (m *Set[T]) Add(vv ...T) {
	for _, v := range vv {
		m.inMap[v] = true
	}
}

func (m *Set[T]) Remove(v T) {
	_, ok := m.inMap[v]
	if ok {
		m.inMap[v] = false
	}
}

func (m *Set[T]) Clear() {
	m.inMap = make(map[T]bool)
}

func (m *Set[T]) IsIn(v T) bool {
	if m.inMap == nil {
		return false
	}

	sv, ok := m.inMap[v]
	return ok && sv
}

func (m *Set[T]) Size() int {
	return len(m.inMap)
}

func (m *Set[T]) ToSlice() []T {
	ret := make([]T, 0)
	for v, sv := range m.inMap {
		if sv {
			ret = append(ret, v)
		}
	}
	return ret
}
