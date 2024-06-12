package mt

import "reflect"

func Max[T numbers](a, b T) T {
	if a > b {
		return a
	}

	return b
}

func Min[T numbers](a, b T) T {
	if a > b {
		return b
	}

	return a
}

func Pow[T intergers](base, exp T) T {
	var result T = 1
	for {
		if exp&1 == 1 {
			result *= base
		}
		exp >>= 1
		if exp == 0 {
			break
		}
		base *= base
	}

	return result
}

func Divide[T intergers](base, divider T) (T, T) {
	return base / divider, base % divider
}

func SumInt64(tt ...interface{}) int64 {
	var totalInt int64
	var totalFloat float64

	for _, t := range tt {
		v := reflect.ValueOf(t)
		switch t.(type) {
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
			totalInt += v.Int()
		case float32, float64:
			totalFloat += v.Float()
		default:
		}
	}

	return totalInt + int64(totalFloat)
}

func SumInt(tt ...interface{}) int {
	return int(SumInt64(tt...))
}

func IsInRange[T numbers](a, b, c T) bool {
	var start, end T
	if b > c {
		start = c
		end = b
	} else {
		start = b
		end = c
	}

	return start <= a && a < end
}

func AbsInt(v int) int {
	if v < 0 {
		return v * -1
	}

	return v
}

func AbsInt64(v int64) int64 {
	if v < 0 {
		return v * -1
	}

	return v
}
