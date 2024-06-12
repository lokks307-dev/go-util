package mt

import (
	"reflect"

	"github.com/lokks307/djson/v2"
	"github.com/volatiletech/null/v8"
)

func Ptr[T any](v T) *T {
	return &v
}

func NullToInt64(v interface{}) (int64, bool) {
	var vv int64
	switch x := v.(type) {
	case null.Int:
		if !x.Valid {
			return 0, false
		}
		vv = int64(x.Int)
	case null.Int8:
		if !x.Valid {
			return 0, false
		}
		vv = int64(x.Int8)
	case null.Int16:
		if !x.Valid {
			return 0, false
		}
		vv = int64(x.Int16)
	case null.Int32:
		if !x.Valid {
			return 0, false
		}
		vv = int64(x.Int32)
	case null.Int64:
		if !x.Valid {
			return 0, false
		}
		vv = int64(x.Int64)
	case null.Uint:
		if !x.Valid {
			return 0, false
		}
		vv = int64(x.Uint)
	case null.Uint8:
		if !x.Valid {
			return 0, false
		}
		vv = int64(x.Uint8)
	case null.Uint16:
		if !x.Valid {
			return 0, false
		}
		vv = int64(x.Uint16)
	case null.Uint32:
		if !x.Valid {
			return 0, false
		}
		vv = int64(x.Uint32)
	case null.Uint64:
		if !x.Valid {
			return 0, false
		}
		vv = int64(x.Uint64)
	default:
		return 0, false
	}

	return vv, true
}

func AnyToInt64(v interface{}) (int64, bool) {
	var vv int64

	switch v.(type) {
	case int, int8, int16, int32, int64:
		el := reflect.ValueOf(v)
		vv = el.Int()
	case uint, uint8, uint16, uint32, uint64:
		el := reflect.ValueOf(v)
		vv = int64(el.Uint())
	case *int, *int8, *int16, *int32, *int64:
		ptr := reflect.ValueOf(v)
		if ptr.IsNil() {
			return 0, false
		}
		vv = ptr.Elem().Int()
	case *uint, *uint8, *uint16, *uint32, *uint64:
		ptr := reflect.ValueOf(v)
		if ptr.IsNil() {
			return 0, false
		}
		vv = int64(ptr.Elem().Uint())
	case null.Int, null.Int8, null.Int16, null.Int32, null.Int64, null.Uint, null.Uint8, null.Uint16, null.Uint32, null.Uint64:
		vx, ok := NullToInt64(v)
		if !ok {
			return 0, false
		}
		vv = vx
	default:
		return 0, false
	}

	return vv, true
}

func AnyToFloat64(v interface{}) (float64, bool) {
	var vv float64

	switch t := v.(type) {
	case float32, float64:
		el := reflect.ValueOf(v)
		vv = el.Float()
	case *float32, *float64:
		ptr := reflect.ValueOf(v)
		if ptr.IsNil() {
			return 0, false
		}
		vv = ptr.Elem().Float()
	case null.Float32:
		if !t.Valid {
			return 0, false
		}
		vv = float64(t.Float32)
	case null.Float64:
		if !t.Valid {
			return 0, false
		}
		vv = t.Float64
	}

	return vv, true
}

func PtrInt64(v interface{}) *int64 {
	vv, ok := AnyToInt64(v)
	if !ok {
		return nil
	}

	return &vv
}

func PtrInt32(v interface{}) *int32 {
	vv, ok := AnyToInt64(v)
	if !ok {
		return nil
	}

	vx := int32(vv)
	return &vx
}

func PtrInt16(v interface{}) *int16 {
	vv, ok := AnyToInt64(v)
	if !ok {
		return nil
	}

	vx := int16(vv)
	return &vx
}

func PtrInt8(v interface{}) *int8 {
	vv, ok := AnyToInt64(v)
	if !ok {
		return nil
	}

	vx := int8(vv)
	return &vx
}

func PtrInt(v interface{}) *int {
	vv, ok := AnyToInt64(v)
	if !ok {
		return nil
	}

	vx := int(vv)
	return &vx
}

func PtrFloat32(v interface{}) *float32 {
	vv, ok := AnyToFloat64(v)
	if !ok {
		return nil
	}
	vx := float32(vv)
	return &vx
}

func PtrFloat64(v interface{}) *float64 {
	vv, ok := AnyToFloat64(v)
	if !ok {
		return nil
	}
	return &vv
}

func PtrStr(v interface{}) *string {
	switch t := v.(type) {
	case *string:
		return t
	case string:
		return &t
	case null.String:
		if !t.Valid {
			return nil
		}

		return &t.String
	}
	return nil
}

func PtrBool(v interface{}) *bool {
	switch t := v.(type) {
	case *bool:
		return t
	case bool:
		return &t
	case null.Bool:
		if !t.Valid {
			return nil
		}

		return &t.Bool
	}
	return nil
}

func PtrJsonToFloat32(d *djson.JSON, key string) *float32 {
	return PtrFloat32(PtrJsonToFloat64(d, key))
}

func PtrJsonToFloat64(d *djson.JSON, key string) *float64 {
	if d == nil || key == "" {
		return nil
	}

	if d.HasKey(key) {
		return PtrFloat64(d.Float(key))
	}

	return nil
}

func PtrJsonToInt(d *djson.JSON, key string) *int {
	return PtrInt(PtrJsonToInt64(d, key))
}

func PtrJsonToInt8(d *djson.JSON, key string) *int8 {
	return PtrInt8(PtrJsonToInt64(d, key))
}

func PtrJsonToInt16(d *djson.JSON, key string) *int16 {
	return PtrInt16(PtrJsonToInt64(d, key))
}

func PtrJsonToInt32(d *djson.JSON, key string) *int32 {
	return PtrInt32(PtrJsonToInt64(d, key))
}

func PtrJsonToInt64(d *djson.JSON, key string) *int64 {
	if d == nil || key == "" {
		return nil
	}

	if d.HasKey(key) {
		return PtrInt64(d.Int(key))
	}

	return nil
}

func PtrJsonToStr(d *djson.JSON, key string) *string {
	if d == nil || key == "" {
		return nil
	}

	if d.HasKey(key) {
		return PtrStr(d.String(key, ""))
	}

	return nil
}

func PtrJsonToBool(d *djson.JSON, key string) *bool {
	if d == nil || key == "" {
		return nil
	}

	if d.HasKey(key) {
		return PtrBool(d.Bool(key))
	}

	return nil
}

func PtrIntPositiveOrNil[T intergers](a T) *T {
	if a <= 0 {
		return nil
	}
	return &a
}

func PtrIntNonZeroOrNil[T intergers](a T) *T {
	if a == 0 {
		return nil
	}
	return &a
}

func PtrToInt64(s interface{}) int64 {
	v, _ := AnyToInt64(s)
	return v
}

func PtrToInt32(s interface{}) int32 {
	v, _ := AnyToInt64(s)
	return int32(v)
}

func PtrToInt(s interface{}) int {
	v, _ := AnyToInt64(s)
	return int(v)
}

func PtrToFloat64(s interface{}) float64 {
	v, _ := AnyToFloat64(s)
	return v
}

func PtrToFloat32(s interface{}) float32 {
	v, _ := AnyToFloat64(s)
	return float32(v)
}

func JsonToPtrFloat32(d *djson.JSON, key string, def ...float32) *float32 {
	defx := ToFloat64Slice(def)
	return PtrFloat32(JsonToPtrFloat64(d, key, defx...))
}

func JsonToPtrFloat64(d *djson.JSON, key string, def ...float64) *float64 {
	if d == nil || key == "" {
		return nil
	}

	if d.HasKey(key) {
		return PtrFloat64(d.Float(key))
	}

	if len(def) > 0 {
		return &def[0]
	} else {
		return nil
	}
}

func JsonToPtrInt(d *djson.JSON, key string, def ...int) *int {
	defx := ToInt64Slice(def)
	return PtrInt(JsonToPtrInt64(d, key, defx...))
}

func JsonToPtrInt8(d *djson.JSON, key string, def ...int8) *int8 {
	defx := ToInt64Slice(def)
	return PtrInt8(JsonToPtrInt64(d, key, defx...))
}

func JsonToPtrInt16(d *djson.JSON, key string, def ...int16) *int16 {
	defx := ToInt64Slice(def)
	return PtrInt16(JsonToPtrInt64(d, key, defx...))
}

func JsonToPtrInt32(d *djson.JSON, key string, def ...int32) *int32 {
	defx := ToInt64Slice(def)
	return PtrInt32(JsonToPtrInt64(d, key, defx...))
}

func JsonToPtrInt64(d *djson.JSON, key string, def ...int64) *int64 {
	if d == nil || key == "" {
		return nil
	}

	if d.HasKey(key) {
		return PtrInt64(d.Int(key))
	}

	if len(def) > 0 {
		return &def[0]
	} else {
		return nil
	}
}

func JsonToPtrStr(d *djson.JSON, key string, def ...string) *string {
	if d == nil || key == "" {
		return nil
	}

	if d.HasKey(key) {
		return PtrStr(d.String(key, ""))
	}

	if len(def) > 0 {
		return &def[0]
	} else {
		return nil
	}
}

func JsonToPtrBool(d *djson.JSON, key string, def ...bool) *bool {
	if d == nil || key == "" {
		return nil
	}

	if d.HasKey(key) {
		return PtrBool(d.Bool(key))
	}

	if len(def) > 0 {
		return &def[0]
	} else {
		return nil
	}
}

func IsZeroInt64(d *int64) bool {
	return d == nil || *d == 0
}

func JsonToPtrStrIfNotNull(d *djson.JSON, key string, def ...string) *string {
	if d == nil || key == "" {
		return nil
	}

	if len(def) > 0 {
		if d.String(key, "") == "null" {
			return &def[0]
		}
	}
	if d.HasKey(key) {
		return PtrStr(d.String(key))
	}

	return nil
}

func JsonToPtrInt32IfNotNull(d *djson.JSON, key string, def ...int32) *int32 {
	if d == nil || key == "" {
		return nil
	}

	if len(def) > 0 {
		if d.String(key, "") == "null" {
			return &def[0]
		}
	}
	if d.HasKey(key) {
		return PtrInt32(d.Int(key))
	}

	return nil
}

func JsonToPtrInt64IfNotNull(d *djson.JSON, key string, def ...int64) *int64 {
	if d == nil || key == "" {
		return nil
	}

	if len(def) > 0 {
		if d.String(key, "") == "null" {
			return &def[0]
		}
	}

	if d.HasKey(key) {
		return PtrInt64(d.Int(key))
	}

	return nil
}
