package mt

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
