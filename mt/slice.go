package mt

import "slices"

func FindIndex[T comparable](slice []T, item T) int {
	return slices.Index(slice, item)
}

func FindIndexStringSlice(slice []string, item string) int {
	return slices.Index(slice, item)
}

func FindIndexInt64Slice(slice []int64, item int64) int {
	return slices.Index(slice, item)
}

func FiltOutInt64Slice(slice []int64, item int64) []int64 {
	nslice := make([]int64, 0)
	for i := range slice {
		if slice[i] != item {
			nslice = append(nslice, slice[i])
		}
	}

	return nslice
}

func IsSameOrderSlice[T comparable](a, b []T) bool {
	return slices.Equal(a, b)
}

func IsSameSlice[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}

	setA := NewSet[T]()
	for i := range a {
		setA.Add(a[i])
	}

	for j := range b {
		if setA.IsIn(b[j]) {
			setA.Remove(b[j])
		} else {
			return false
		}
	}

	return setA.Size() == 0
}
