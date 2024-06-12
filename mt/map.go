package mt

func MapToSlice[T comparable](ss map[T]bool) []T {
	v := make([]T, 0)
	for s := range ss {
		v = append(v, s)
	}

	return v
}

func GetUniqueSlice[T comparable](ss []T) []T {
	v := make(map[T]bool)
	for _, s := range ss {
		v[s] = true
	}

	return MapToSlice(v)
}
