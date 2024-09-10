package mt

import (
	"reflect"
	"strings"

	"github.com/lokks307/djson/v2"
	"github.com/volatiletech/null/v8"
)

func AddIfNotEmptyArray(dst *djson.JSON, key string, v any) {
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

func AddIfNotEmpty(dst *djson.JSON, key string, v any) {
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

func UpdateNullIfNotEmpty(dst any, src *djson.JSON, key string, allowEmpty ...bool) bool {
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

func UpdateJsonIfValidNull(dst *djson.JSON, key string, src any, allow ...int8) bool {
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
	res := make([]string, 0)
	dd := djson.New().Parse(ss)
	if dd.IsArray() {
		dd.Seek()
		for dd.Next() {
			res = append(res, dd.Scan().String())
		}
		return res
	}

	return res
}

func JsonStringToInterfaceSlice(ss string) []any {
	res := make([]any, 0)
	dd := djson.New().Parse(ss)
	if dd.IsArray() {
		dd.Seek()
		for dd.Next() {
			res = append(res, dd.Scan().Interface())
		}
		return res
	}

	return res
}

func JsonStringToIntegerSlice[T integers](ss string) []T {
	res := make([]T, 0)
	dd := djson.New().Parse(ss)
	if dd.IsArray() {
		dd.Seek()
		for dd.Next() {
			res = append(res, T(dd.Scan().Int()))
		}
		return res
	}
	return res
}

func JsonStringToIntSlice(ss string) []int {
	return JsonStringToIntegerSlice[int](ss)
}

func JsonStringToInt8Slice(ss string) []int8 {
	return JsonStringToIntegerSlice[int8](ss)
}

func JsonStringToInt16Slice(ss string) []int16 {
	return JsonStringToIntegerSlice[int16](ss)
}

func JsonStringToInt32Slice(ss string) []int32 {
	return JsonStringToIntegerSlice[int32](ss)
}

func JsonStringToInt64Slice(ss string) []int64 {
	return JsonStringToIntegerSlice[int64](ss)
}

func JsonStringToUintSlice(ss string) []uint {
	return JsonStringToIntegerSlice[uint](ss)
}

func JsonStringToUint8Slice(ss string) []uint8 {
	return JsonStringToIntegerSlice[uint8](ss)
}

func JsonStringToUint16Slice(ss string) []uint16 {
	return JsonStringToIntegerSlice[uint16](ss)
}

func JsonStringToUint32Slice(ss string) []uint32 {
	return JsonStringToIntegerSlice[uint32](ss)
}

func JsonStringToUint64Slice(ss string) []uint64 {
	return JsonStringToIntegerSlice[uint64](ss)
}

func JsonToStringSlice(js *djson.JSON, key ...string) []string {
	return djson.JsonToStringSlice(js, key...)
}

func JsonToIntegerSlice[T integers](js *djson.JSON, key ...string) []T {
	if js == nil || !js.IsArray() || js.Size() == 0 {
		return []T{}
	}

	ss := make([]T, 0)
	js.Seek()
	if len(key) > 0 {
		for js.Next() {
			ec := js.Scan()
			ss = append(ss, T(ec.Int(key[0])))
		}
	} else {
		for js.Next() {
			ec := js.Scan()
			ss = append(ss, T(ec.Int()))
		}
	}
	return ss
}

func JsonToIntSlice(js *djson.JSON, key ...string) []int {
	return JsonToIntegerSlice[int](js, key...)
}

func JsonToInt8Slice(js *djson.JSON, key ...string) []int8 {
	return JsonToIntegerSlice[int8](js, key...)
}

func JsonToInt16Slice(js *djson.JSON, key ...string) []int16 {
	return JsonToIntegerSlice[int16](js, key...)
}

func JsonToInt326Slice(js *djson.JSON, key ...string) []int32 {
	return JsonToIntegerSlice[int32](js, key...)
}

func JsonToInt64Slice(js *djson.JSON, key ...string) []int64 {
	return JsonToIntegerSlice[int64](js, key...)
}

func JsonToUintSlice(js *djson.JSON, key ...string) []uint {
	return JsonToIntegerSlice[uint](js, key...)
}

func JsonToUint8Slice(js *djson.JSON, key ...string) []uint8 {
	return JsonToIntegerSlice[uint8](js, key...)
}

func JsonToUint16Slice(js *djson.JSON, key ...string) []uint16 {
	return JsonToIntegerSlice[uint16](js, key...)
}

func JsonToUint32Slice(js *djson.JSON, key ...string) []uint32 {
	return JsonToIntegerSlice[uint32](js, key...)
}

func JsonToUint64Slice(js *djson.JSON, key ...string) []uint64 {
	return JsonToIntegerSlice[uint64](js, key...)
}

func JsonElementToIntegerSlice[T integers](js *djson.JSON, el string) []T {
	if js == nil {
		return []T{}
	}

	if el == "" {
		return JsonToIntegerSlice[T](js)
	}

	if !js.IsArray(el) {
		return []T{}
	}

	dd, _ := js.Array(el)
	return JsonToIntegerSlice[T](dd)
}

func JsonElementToStringSlice(js *djson.JSON, el string) []string {
	if js == nil {
		return []string{}
	}

	if el == "" {
		return JsonToStringSlice(js)
	}

	if !js.IsArray(el) {
		return []string{}
	}

	dd, _ := js.Array(el)
	return JsonToStringSlice(dd)
}

func JsonElementToIntSlice(js *djson.JSON, el string) []int {
	return JsonElementToIntegerSlice[int](js, el)
}
func JsonElementToInt8Slice(js *djson.JSON, el string) []int8 {
	return JsonElementToIntegerSlice[int8](js, el)
}
func JsonElementToInt16Slice(js *djson.JSON, el string) []int16 {
	return JsonElementToIntegerSlice[int16](js, el)
}
func JsonElementToInt32Slice(js *djson.JSON, el string) []int32 {
	return JsonElementToIntegerSlice[int32](js, el)
}
func JsonElementToInt64Slice(js *djson.JSON, el string) []int64 {
	return JsonElementToIntegerSlice[int64](js, el)
}
func JsonElementToUntSlice(js *djson.JSON, el string) []uint {
	return JsonElementToIntegerSlice[uint](js, el)
}
func JsonElementToUint8Slice(js *djson.JSON, el string) []uint8 {
	return JsonElementToIntegerSlice[uint8](js, el)
}
func JsonElementToUint16Slice(js *djson.JSON, el string) []uint16 {
	return JsonElementToIntegerSlice[uint16](js, el)
}
func JsonElementToUint32Slice(js *djson.JSON, el string) []uint32 {
	return JsonElementToIntegerSlice[uint32](js, el)
}
func JsonElementToUint64Slice(js *djson.JSON, el string) []uint64 {
	return JsonElementToIntegerSlice[uint64](js, el)
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

func JsonPtrFloat32(d *djson.JSON, key string) *float32 {
	return JsonToPtrFloat32(d, key)
}

func JsonToPtrFloat32(d *djson.JSON, key string, def ...float32) *float32 {
	defx := ToFloat64Slice(def)
	return PtrFloat32(JsonToPtrFloat64(d, key, defx...))
}

func JsonPtrFloat64(d *djson.JSON, key string) *float64 {
	return JsonToPtrFloat64(d, key)
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

func JsonPtrInt(d *djson.JSON, key string) *int {
	return JsonToPtrInt(d, key)
}

func JsonToPtrInt(d *djson.JSON, key string, def ...int) *int {
	defx := ToInt64Slice(def)
	return PtrInt(JsonToPtrInt64(d, key, defx...))
}

func JsonPtrInt8(d *djson.JSON, key string) *int8 {
	return JsonToPtrInt8(d, key)
}

func JsonToPtrInt8(d *djson.JSON, key string, def ...int8) *int8 {
	defx := ToInt64Slice(def)
	return PtrInt8(JsonToPtrInt64(d, key, defx...))
}

func JsonPtrInt16(d *djson.JSON, key string) *int16 {
	return JsonToPtrInt16(d, key)
}

func JsonToPtrInt16(d *djson.JSON, key string, def ...int16) *int16 {
	defx := ToInt64Slice(def)
	return PtrInt16(JsonToPtrInt64(d, key, defx...))
}

func JsonPtrInt32(d *djson.JSON, key string) *int32 {
	return JsonToPtrInt32(d, key)
}

func JsonToPtrInt32(d *djson.JSON, key string, def ...int32) *int32 {
	defx := ToInt64Slice(def)
	return PtrInt32(JsonToPtrInt64(d, key, defx...))
}

func JsonPtrInt64(d *djson.JSON, key string) *int64 {
	return JsonToPtrInt64(d, key)
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

func JsonPtrBool(d *djson.JSON, key string) *bool {
	return JsonToPtrBool(d, key)
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

func JsonPtrArray(d *djson.JSON, key string) *djson.JSON {
	if d == nil || !d.HasKey(key) {
		return nil
	}

	arr, ok := d.Array(key)
	if !ok || !arr.IsArray() {
		return nil
	}

	return arr
}

func JsonPtrObject(d *djson.JSON, key string) *djson.JSON {
	if d == nil || !d.HasKey(key) {
		return nil
	}

	obj, ok := d.Object(key)
	if !ok || !obj.IsObject() {
		return nil
	}

	return obj
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

func JsonToNullInt(d *djson.JSON, key string) null.Int {
	return null.IntFromPtr(JsonToPtrInt(d, key))
}

func JsonToNullInt8(d *djson.JSON, key string) null.Int8 {
	return null.Int8FromPtr(JsonToPtrInt8(d, key))
}

func JsonToNullInt16(d *djson.JSON, key string) null.Int16 {
	return null.Int16FromPtr(JsonToPtrInt16(d, key))
}
func JsonToNullInt32(d *djson.JSON, key string) null.Int32 {
	return null.Int32FromPtr(JsonToPtrInt32(d, key))
}

func JsonToNullInt64(d *djson.JSON, key string) null.Int64 {
	return null.Int64FromPtr(JsonToPtrInt64(d, key))
}

func JsonToNullString(d *djson.JSON, key string) null.String {
	return null.StringFromPtr(JsonToPtrStr(d, key))
}

func JsonToNullBool(d *djson.JSON, key string) null.Bool {
	return null.BoolFromPtr(JsonToPtrBool(d, key))
}
