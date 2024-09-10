package mt

import (
	"reflect"

	"github.com/lokks307/djson/v2"
	"github.com/volatiletech/null/v8"
)

func HasEmptyValue(vv ...any) bool {
	if len(vv) == 0 {
		return false
	}

	for _, v := range vv {
		if v == nil {
			return true
		}

		switch t := v.(type) {
		case int, int8, int16, int32, int64:
			el := reflect.ValueOf(v)
			if el.Int() <= 0 {
				return true
			}
		case uint, uint8, uint16, uint32, uint64:
			el := reflect.ValueOf(v)
			if el.Uint() == 0 {
				return true
			}
		case *int, *int8, *int16, *int32, *int64:
			ptr := reflect.ValueOf(v)
			if ptr.IsNil() {
				return true
			}
			if ptr.Elem().Int() <= 0 {
				return true
			}
		case *uint, *uint8, *uint16, *uint32, *uint64:
			ptr := reflect.ValueOf(v)
			if ptr.IsNil() {
				return true
			}
			if ptr.Elem().Uint() == 0 {
				return true
			}
		case null.Int, null.Int8, null.Int16, null.Int32, null.Int64, null.Uint, null.Uint8, null.Uint16, null.Uint32, null.Uint64:
			vx, ok := NullToInt64(v)
			if !ok || vx <= 0 {
				return true
			}
		case string:
			if t == "" {
				return true
			}
		case *string:
			if t == nil || *t == "" {
				return true
			}
		case *djson.JSON:
			if t == nil {
				return true
			}
		default:
			rv := reflect.ValueOf(v)
			rk := rv.Kind()
			switch rk {
			case reflect.Chan, reflect.Func, reflect.Map, reflect.Pointer, reflect.UnsafePointer, reflect.Interface, reflect.Slice:
				return rv.IsNil()
			}

		}
	}

	return false
}
