package mt

import (
	"reflect"

	"github.com/lokks307/djson/v2"
	"github.com/volatiletech/null/v8"
)

func AddIfNotEmptyArray(dst *djson.JSON, key string, v interface{}) {
	if dst == nil || key == "" {
		return
	}

	switch t := v.(type) {
	case string:
		tt := djson.New().Parse(t)
		if tt.IsArray() {
			dst.Put(key, tt)
		}
	case *string:
		if t != nil && *t != "" {
			tt := djson.New().Parse(*t)
			if tt.IsArray() {
				dst.Put(key, tt)
			}
		}
	case null.String:
		if t.Valid && t.String != "" {
			tt := djson.New().Parse(t.String)
			if tt.IsArray() {
				dst.Put(key, tt)
			}
		}
	}
}

func AddIfNotEmpty(dst *djson.JSON, key string, v interface{}) {
	if dst == nil || key == "" {
		return
	}

	switch t := v.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		el := reflect.ValueOf(v)
		vv := el.Int()
		if vv != 0 {
			dst.Put(key, vv)
		}
	case *int, *int8, *int16, *int32, *int64, *uint, *uint8, *uint16, *uint32, *uint64:
		ptr := reflect.ValueOf(v)
		if !ptr.IsNil() {
			vv := ptr.Elem().Int()
			if vv != 0 {
				dst.Put(key, vv)
			}
		}
	case float32, float64:
		el := reflect.ValueOf(v)
		vv := el.Float()
		if !IsAlmostEqual(vv, 0) {
			dst.Put(key, vv)
		}
	case *float32, *float64:
		ptr := reflect.ValueOf(v)
		if !ptr.IsNil() {
			vv := ptr.Elem().Float()
			if !IsAlmostEqual(vv, 0) {
				dst.Put(key, vv)
			}
		}
	case null.Int, null.Int8, null.Int16, null.Int32, null.Int64, null.Uint, null.Uint8, null.Uint16, null.Uint32, null.Uint64:
		vv, ok := NullToInt64(v)
		if ok && vv != 0 {
			dst.Put(key, vv)
		}
	case bool:
		dst.Put(key, t)
	case *bool:
		if t != nil {
			dst.Put(key, *t)
		}
	case *string:
		if t != nil && *t != "" {
			dst.Put(key, *t)
		}
	case string:
		if t != "" {
			dst.Put(key, t)
		}
	case null.String:
		if t.Valid && t.String != "" {
			dst.Put(key, t.String)
		}
	}
}

func PutIfHasKey(dst *djson.JSON, src *djson.JSON, keys ...string) {
	if len(keys) < 1 || dst == nil || src == nil {
		return
	}

	srcKey := keys[0]
	dstKey := srcKey

	if len(keys) > 1 {
		dstKey = keys[1]
	}

	if src.HasKey(srcKey) {
		dst.Put(dstKey, src.Interface(srcKey))
	}
}

func IntSliceToArray[T intergers](ss []T) djson.Array {
	var ret djson.Array
	for _, s := range ss {
		ret = append(ret, s)
	}

	return ret
}
