package mt

import (
	"reflect"
	"strings"

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

func IntSliceToArray[T integers](ss []T) djson.Array {
	var ret djson.Array
	for _, s := range ss {
		ret = append(ret, s)
	}

	return ret
}

func HasNoneEmpty(opt *djson.JSON, keys ...string) bool {
	if opt == nil {
		return false
	}
	for _, key := range keys {
		if opt.String(key) == "" {
			return false
		}
	}
	return true
}

func GetStringIfNoneEmpty(opt *djson.JSON, keys ...string) string {
	if len(keys) == 0 {
		return ""
	}

	defaultRet := ""
	if len(keys) > 2 {
		defaultRet = keys[1]
	}

	if opt == nil || !HasNoneEmpty(opt, keys[0]) {
		return defaultRet
	}

	return opt.String(keys[0])
}

func GetIntIfNoneEmpty(opt *djson.JSON, key string, defInt ...int64) int64 {

	var defaultRet int64
	if len(defInt) > 1 {
		defaultRet = defInt[0]
	}

	if opt == nil || !HasNoneEmpty(opt, key) {
		return defaultRet
	}

	return opt.Int(key)
}

func UpdateValuesIfExist(dst *djson.JSON, src *djson.JSON, keys ...string) {
	if dst == nil || src == nil {
		return
	}

	for _, key := range keys {
		if key != "" && src.HasKey(key) {
			dst.Put(key, src.Interface(key))
		}
	}
}

func UpdateNullIfNotEmpty(dst interface{}, src *djson.JSON, key string, allowEmpty ...bool) bool {
	if dst == nil || src == nil || key == "" || !src.HasKey(key) {
		return false
	}

	if len(allowEmpty) > 0 && allowEmpty[0] && !HasNoneEmpty(src, key) {
		return false
	}

	switch obj := dst.(type) {
	case *null.String:
		obj.SetValid(src.String(key))
	case *null.Bool:
		obj.SetValid(src.Bool(key))
	case *null.Int:
		obj.SetValid(int(src.Int(key)))
	case *null.Int8:
		obj.SetValid(int8(src.Int(key)))
	case *null.Int16:
		obj.SetValid(int16(src.Int(key)))
	case *null.Int32:
		obj.SetValid(int32(src.Int(key)))
	case *null.Int64:
		obj.SetValid(int64(src.Int(key)))
	case *null.Uint:
		obj.SetValid(uint(src.Int(key)))
	case *null.Uint8:
		obj.SetValid(uint8(src.Int(key)))
	case *null.Uint16:
		obj.SetValid(uint16(src.Int(key)))
	case *null.Uint32:
		obj.SetValid(uint32(src.Int(key)))
	case *null.Uint64:
		obj.SetValid(uint64(src.Int(key)))
	case *null.Float32:
		obj.SetValid(float32(src.Float(key)))
	case *null.Float64:
		obj.SetValid(src.Float(key))
	default:
		return false
	}

	return true
}

const (
	ALLOW_BOOL   int8 = 0x01
	ALLOW_STRING int8 = 0x02
	ALLOW_INT    int8 = 0x04
	ALLOW_FLOAT  int8 = 0x08
	ALLOW_OBJECT int8 = 0x10
	ALLOW_ARRAY  int8 = 0x20
	ALLOW_ALL    int8 = 0x3F // 0x01 | 0x02 | 0x04 | 0x08 | 0x10 | 0x20
)

func CheckMask[T integers](s T, v T) bool {
	return s&v != 0
}

func UpdateJsonIfValidNull(dst *djson.JSON, key string, src interface{}, allow ...int8) bool {
	if src == nil || dst == nil || key == "" {
		return false
	}

	allowFlag := ALLOW_ALL
	if len(allow) > 0 {
		allowFlag = allow[0]
	}

	switch t := src.(type) {
	case null.String:
		if t.Valid {
			if t.String == "" {
				if CheckMask(allowFlag, ALLOW_STRING) {
					dst.Put(key, "")
				}
			} else { // can be json string
				ss := djson.New().Parse(t.String)
				if ss.IsArray() {
					if CheckMask(allowFlag, ALLOW_ARRAY) {
						dst.Put(key, ss)
					}
				} else if ss.IsObject() {
					if CheckMask(allowFlag, ALLOW_OBJECT) {
						dst.Put(key, ss)
					}
				} else {
					if CheckMask(allowFlag, ALLOW_STRING) {
						dst.Put(key, t.String)
					}
				}
			}
		}
	case null.Int:
		if t.Valid && CheckMask(allowFlag, ALLOW_INT) {
			dst.Put(key, t.Int)
		}
	case null.Int8:
		if t.Valid && CheckMask(allowFlag, ALLOW_INT) {
			dst.Put(key, t.Int8)
		}
	case null.Int16:
		if t.Valid && CheckMask(allowFlag, ALLOW_INT) {
			dst.Put(key, t.Int16)
		}
	case null.Int32:
		if t.Valid && CheckMask(allowFlag, ALLOW_INT) {
			dst.Put(key, t.Int32)
		}
	case null.Int64:
		if t.Valid && CheckMask(allowFlag, ALLOW_INT) {
			dst.Put(key, t.Int64)
		}
	case null.Uint:
		if t.Valid && CheckMask(allowFlag, ALLOW_INT) {
			dst.Put(key, t.Uint)
		}
	case null.Uint8:
		if t.Valid && CheckMask(allowFlag, ALLOW_INT) {
			dst.Put(key, t.Uint8)
		}
	case null.Uint16:
		if t.Valid && CheckMask(allowFlag, ALLOW_INT) {
			dst.Put(key, t.Uint16)
		}
	case null.Uint32:
		if t.Valid && CheckMask(allowFlag, ALLOW_INT) {
			dst.Put(key, t.Uint32)
		}
	case null.Uint64:
		if t.Valid && CheckMask(allowFlag, ALLOW_INT) {
			dst.Put(key, t.Uint64)
		}
	case null.Bool:
		if t.Valid && CheckMask(allowFlag, ALLOW_BOOL) {
			dst.Put(key, t.Bool)
		}
	case null.Float32:
		if t.Valid && CheckMask(allowFlag, ALLOW_FLOAT) {
			dst.Put(key, t.Float32)
		}
	case null.Float64:
		if t.Valid && CheckMask(allowFlag, ALLOW_FLOAT) {
			dst.Put(key, t.Float64)
		}
	default:
		return false
	}

	return true
}

func JsonStringToObject(ss string) *djson.JSON {
	dd := djson.New().Parse(ss)
	if dd.IsObject() {
		return dd
	}
	return djson.NewObject()
}

func JsonStringToArray(ss string) *djson.JSON {
	dd := djson.New().Parse(ss)
	if dd.IsArray() {
		return dd
	}
	return djson.NewArray()
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

func JsonStringToInt64Slice(ss string) []int64 {
	dd := djson.New().Parse(ss)
	if dd.IsArray() {
		return djson.JsonToInt64Slice(dd)
	}
	return []int64{}
}

func JsonStringToInt8Slice(ss string) []int8 {
	res := make([]int8, 0)
	dd := djson.New().Parse(ss)
	if dd.IsArray() {
		dd.Seek()
		for dd.Next() {
			res = append(res, int8(dd.Scan().Int()))
		}
		return res
	}
	return res
}

func JsonStringToInterfaceSlice(ss string) []interface{} {
	res := make([]interface{}, 0)
	dd := djson.New().Parse(ss)
	if dd.IsArray() {
		arr := djson.JsonToInt64Slice(dd)
		for _, v := range arr {
			res = append(res, v)
		}
	}

	return res
}

func JsonStringToInt32Slice(ss string) []int32 {
	dd := djson.New().Parse(ss)
	if dd.IsArray() {
		res := []int32{}
		dd.Seek()
		for dd.Next() {
			res = append(res, int32(dd.Scan().Int()))
		}
		return res
	}
	return []int32{}
}

func JsonToUint64Slice(js *djson.JSON, key ...string) []uint64 {
	if js == nil || !js.IsArray() || js.Size() == 0 {
		return []uint64{}
	}

	ss := make([]uint64, 0)
	js.Seek()
	if len(key) > 0 {
		for js.Next() {
			ec := js.Scan()
			ss = append(ss, uint64(ec.Int(key[0])))
		}
	} else {
		for js.Next() {
			ec := js.Scan()
			ss = append(ss, uint64(ec.Int()))
		}
	}
	return ss
}

func RemoveDuplicatedTag(sliceList *djson.JSON) *djson.JSON {
	allKeys := make(map[string]bool)
	list := djson.NewArray()
	sliceList.Seek()
	for sliceList.Next() {
		val := sliceList.Scan().String()
		val = strings.TrimLeft(val, "#")

		if val == "" {
			continue
		}

		if _, value := allKeys[val]; !value {
			allKeys[val] = true
			list.PutArray(val)
		}
	}

	return list
}

func AddCPUTagForNotEnrolledUser(id string) string {
	tag := "CPU_"
	return tag + id
}

func MustGetStringOfSlice(x string) string {
	d := djson.New().Parse(x)
	if !d.IsArray() {
		return "[]"
	}

	return x
}

func IsEmptyJsonString(x string) bool {
	xx := strings.ReplaceAll(x, " ", "")
	return IsIn(xx, "", "[]", "{}")
}

func JsonPtrStr(d *djson.JSON, key ...string) *string {
	if d != nil {
		if len(key) == 0 {
			return PtrStr(d.ToString())
		} else {
			if !d.HasKey(key[0]) {
				return nil
			}

			return PtrStr(d.String(key[0]))
		}
	}

	return nil
}

func JsonPtrBool(d *djson.JSON, key string) *bool {
	if d == nil || !d.HasKey(key) {
		return nil
	}
	return PtrBool(d.Bool(key))
}

func JsonPtrFloat64(d *djson.JSON, key string) *float64 {
	if d == nil || !d.HasKey(key) {
		return nil
	}
	return PtrFloat64(d.Float(key))
}

func JsonPtrFloat32(d *djson.JSON, key string) *float32 {
	if d == nil || !d.HasKey(key) {
		return nil
	}
	return PtrFloat32(d.Float(key))
}

func JsonPtrInt(d *djson.JSON, key string) *int {
	if d == nil || !d.HasKey(key) {
		return nil
	}
	return PtrInt(d.Int(key))
}

func JsonPtrInt64(d *djson.JSON, key string) *int64 {
	if d == nil || !d.HasKey(key) {
		return nil
	}
	return PtrInt64(d.Int(key))
}

func JsonPtrInt8(d *djson.JSON, key string) *int8 {
	if d == nil || !d.HasKey(key) {
		return nil
	}
	return PtrInt8(d.Int(key))
}

func JsonPtrArray(d *djson.JSON, key string) *djson.JSON {
	if d == nil || !d.HasKey(key) {
		return nil
	}
	return djson.MustGetArray(d, key)
}

func JsonPtrObject(d *djson.JSON, key string) *djson.JSON {
	if d == nil || !d.HasKey(key) {
		return nil
	}
	return djson.MustGetObject(d, key)
}
