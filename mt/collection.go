package mt

import (
	"reflect"

	"github.com/volatiletech/null/v8"
)

func Index[T comparable](vs []T, t T) int {
	for i, v := range vs {
		if v == t {
			return i
		}
	}
	return -1
}

func Include[T comparable](vs []T, t T) bool {
	return Index(vs, t) >= 0
}

func Any[T comparable](vs []T, f func(T) bool) bool {
	for _, v := range vs {
		if f(v) {
			return true
		}
	}
	return false
}

func All[T comparable](vs []T, f func(T) bool) bool {
	for _, v := range vs {
		if !f(v) {
			return false
		}
	}
	return true
}

func Filter[T comparable](vs []T, f func(T) bool) []T {
	vsf := make([]T, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

func Map[T comparable](vs []T, f func(T) T) []T {
	vsm := make([]T, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}

func IsIn[T comparable](v T, opts ...T) bool {
	for _, e := range opts {
		if e == v {
			return true
		}
	}

	return false
}

func IsInEnum(x any, c any) bool {

	xx := reflect.ValueOf(x)
	cc := reflect.ValueOf(c)

	nullInt, nullIntOk := NullToInt64(c)
	ccc, nullStrOk := c.(null.String)
	nullStr := ccc.String

	if xx.Kind() == reflect.Ptr {
		xx = xx.Elem()
	}

	if xx.Kind() != reflect.Struct {
		return false
	}

	for i := 0; i < xx.NumField(); i++ {
		xxx := xx.Field(i)

		if xxx.Kind() == reflect.Struct {
			continue
		}

		if xxx.Kind() == reflect.String {
			if nullStrOk {
				if xxx.String() == nullStr {
					return true
				}
			} else {
				if xxx.String() == cc.String() {
					return true
				}
			}
		} else {
			if xxx.CanInt() {
				if nullIntOk {
					if xxx.Int() == nullInt {
						return true
					}
				} else {
					if cc.CanInt() && xxx.Int() == cc.Int() {
						return true
					}
				}
			}
		}
	}

	return false
}

func IsInInt64(v int64, opts ...int64) bool {
	for _, e := range opts {
		if e == v {
			return true
		}
	}

	return false
}

func IsInInt(v int, opts ...int) bool {
	for _, e := range opts {
		if e == v {
			return true
		}
	}

	return false
}

func IsInStr(v string, opts ...string) bool {
	for _, e := range opts {
		if e == v {
			return true
		}
	}

	return false
}
