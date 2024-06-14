package mt

func FindIndex[T comparable](slice []T, item T) int {
	for i := range slice {
		if slice[i] == item {
			return i
		}
	}
	return -1
}

func FindIndexStringSlice(slice []string, item string) int {
	return FindIndex(slice, item)
}

func FindIndexInt64Slice(slice []int64, item int64) int {
	return FindIndex(slice, item)
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
	if len(a) != len(b) {
		return false
	}

	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func IsSameSlice[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}

	aa := make(map[T]bool)

	for i := range a {
		aa[a[i]] = false
	}

	for j := range b {
		if _, ok := aa[b[j]]; ok {
			aa[b[j]] = true
		} else {
			return false
		}
	}

	for _, k := range aa {
		if !k {
			return false
		}
	}

	return true
}
