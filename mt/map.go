package mt

func MapToSlice[T comparable](ss map[T]bool) []T {
	v := make([]T, 0)
	for s := range ss {
		v = append(v, s)
	}

	return v
}

func GetUniqueSlice[T comparable](ss []T) []T {
	set := NewSet[T]()
	for _, s := range ss {
		set.Add(s)
	}
	return set.ToSlice()
}
