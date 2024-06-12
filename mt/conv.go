package mt

import (
	"strconv"
	"strings"

	"github.com/lokks307/djson/v2"
)

func StringSliceToAnySlice(ss []string) []interface{} {
	is := make([]interface{}, len(ss))
	for i, v := range ss {
		is[i] = MysqlRealEscapeString(v)
	}

	return is
}

func IntSliceToAnySlice[T intergers](ss []T) []interface{} {
	is := make([]interface{}, len(ss))
	for i, v := range ss {
		is[i] = v
	}

	return is
}

func FormatStringSlice(n int) string {
	return strings.TrimSuffix(strings.Repeat(`'%s',`, n), ",")
}

func MapKeyToInt64Slice(m map[int64]bool) []int64 {
	l := make([]int64, 0)

	for k := range m {
		l = append(l, k)
	}

	return l
}

func JsonStringToStringSlice(ss string) []string {
	dd := djson.New().Parse(ss)
	if dd.IsArray() {
		return djson.ArrayJsonToStringSlice(dd)
	}
	return []string{}
}

func JsonStringToIntSlice(ss string) []int {
	dd := djson.New().Parse(ss)
	if dd.IsArray() {
		return djson.JsonToIntSlice(dd)
	}
	return []int{}
}

func ToIntSlice[T intergers](ss []T) []int {
	r := make([]int, 0)
	for _, s := range ss {
		r = append(r, int(s))
	}

	return r
}

func ToInt8Slice[T intergers](ss []T) []int8 {
	r := make([]int8, 0)
	for _, s := range ss {
		r = append(r, int8(s))
	}

	return r
}

func ToInt32Slice[T intergers](ss []T) []int32 {
	r := make([]int32, 0)
	for _, s := range ss {
		r = append(r, int32(s))
	}

	return r
}

func ToInt64Slice[T intergers](ss []T) []int64 {
	r := make([]int64, 0)
	for _, s := range ss {
		r = append(r, int64(s))
	}

	return r
}

func ToFloat32Slice[T numbers](ss []T) []float32 {
	r := make([]float32, len(ss))
	for i, s := range ss {
		r[i] = float32(s)
	}

	return r
}

func ToFloat64Slice[T numbers](ss []T) []float64 {
	r := make([]float64, len(ss))
	for i, s := range ss {
		r[i] = float64(s)
	}

	return r
}

func ToString[T intergers](v T) string {
	return strconv.FormatInt(int64(v), 10)
}

func ToStringSlice[T intergers](ss []T) []string {
	r := make([]string, len(ss))
	for i, s := range ss {
		r[i] = ToString(s)
	}

	return r
}
