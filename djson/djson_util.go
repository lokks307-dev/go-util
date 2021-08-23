package djson

import (
	"reflect"
	"strings"

	"github.com/volatiletech/null/v8"
)

func (m *DJSON) Size() int {
	return m.Length()
}

func (m *DJSON) Length() int {
	if m.JsonType == JSON_NULL {
		return 0
	}

	if m.JsonType == JSON_ARRAY {
		return m.Array.Length()
	}

	if m.JsonType == JSON_OBJECT {
		return m.Object.Length()
	}

	return 1
}

func (m *DJSON) HasKey(key interface{}) bool {
	switch tkey := key.(type) {
	case string:
		if m.JsonType == JSON_OBJECT {
			return m.Object.HasKey(tkey)
		}
	case int:
		if m.JsonType == JSON_ARRAY {
			return tkey >= 0 && m.Array.Size() > tkey
		}
	}

	return false
}

func (m *DJSON) toFieldsValue(val reflect.Value, tags ...string) {

	for i := 0; i < val.NumField(); i++ {
		eachVal := val.Field(i)
		eachType := val.Type().Field(i)
		eachTag := eachType.Tag.Get("json")

		if !eachVal.CanSet() || !m.HasKey(eachTag) {
			continue
		}

		if len(tags) > 0 && !inTags(eachTag, tags...) {
			continue
		}

		eachKind := eachType.Type.Kind()

		if eachKind == reflect.Struct {

			switch eachType.Type.String() {
			case "null.String":
				eachVal.FieldByName("String").SetString(m.GetAsString(eachTag))
				eachVal.FieldByName("Valid").SetBool(true)
			case "null.Bool":
				eachVal.FieldByName("Bool").SetBool(m.GetAsBool(eachTag))
				eachVal.FieldByName("Valid").SetBool(true)
			case "null.Float32":
				eachVal.FieldByName("Float32").SetFloat(m.GetAsFloat(eachTag))
				eachVal.FieldByName("Valid").SetBool(true)
			case "null.Float64":
				eachVal.FieldByName("Float64").SetFloat(m.GetAsFloat(eachTag))
				eachVal.FieldByName("Valid").SetBool(true)
			case "null.Int":
				eachVal.FieldByName("Int").SetInt(m.GetAsInt(eachTag))
				eachVal.FieldByName("Valid").SetBool(true)
			case "null.Int8":
				eachVal.FieldByName("Int8").SetInt(m.GetAsInt(eachTag))
				eachVal.FieldByName("Valid").SetBool(true)
			case "null.Int16":
				eachVal.FieldByName("Int16").SetInt(m.GetAsInt(eachTag))
				eachVal.FieldByName("Valid").SetBool(true)
			case "null.Int32":
				eachVal.FieldByName("Int32").SetInt(m.GetAsInt(eachTag))
				eachVal.FieldByName("Valid").SetBool(true)
			case "null.Int64":
				eachVal.FieldByName("Int64").SetInt(m.GetAsInt(eachTag))
				eachVal.FieldByName("Valid").SetBool(true)
			case "null.Uint":
				eachVal.FieldByName("Uint").SetUint(uint64(m.GetAsInt(eachTag)))
				eachVal.FieldByName("Valid").SetBool(true)
			case "null.Uint8":
				eachVal.FieldByName("Uint8").SetUint(uint64(m.GetAsInt(eachTag)))
				eachVal.FieldByName("Valid").SetBool(true)
			case "null.Uint16":
				eachVal.FieldByName("Uint16").SetUint(uint64(m.GetAsInt(eachTag)))
				eachVal.FieldByName("Valid").SetBool(true)
			case "null.Uint32":
				eachVal.FieldByName("Uint32").SetUint(uint64(m.GetAsInt(eachTag)))
				eachVal.FieldByName("Valid").SetBool(true)
			case "null.Uint64":
				eachVal.FieldByName("Uint64").SetUint(uint64(m.GetAsInt(eachTag)))
				eachVal.FieldByName("Valid").SetBool(true)
			default:

				if oJson, ok := m.GetAsObject(eachTag); ok {
					oJson.toFieldsValue(eachVal, downDepthWW(tags)...)
				}

			}

		} else {

			switch eachType.Type.String() {
			case "int", "int8", "int16", "int32", "int64":
				eachVal.SetInt(m.GetAsInt(eachTag))
			case "uint", "uint8", "uint16", "uint32", "uint64":
				eachVal.SetUint(uint64(m.GetAsInt(eachTag)))
			case "float32", "float64":
				eachVal.SetFloat(m.GetAsFloat(eachTag))
			case "string":
				eachVal.SetString(m.GetAsString(eachTag))
			case "bool":
				eachVal.SetBool(m.GetAsBool(eachTag))
			}
		}
	}
}

func (m *DJSON) ToFields(st interface{}, tags ...string) {
	target := reflect.ValueOf(st)
	elements := target.Elem()
	m.toFieldsValue(elements, tags...)
}

func (m *DJSON) fromFieldsValue(val reflect.Value, tags ...string) {

	kind := val.Type().Kind()

	if kind == reflect.Array || kind == reflect.Slice {

		for i := 0; i < val.Len(); i++ {
			eachVal := val.Index(i)
			eachType := eachVal.Type()

			switch eachVal.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				m.PutAsArray(eachVal.Int())
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				m.PutAsArray(eachVal.Uint())
			case reflect.Bool:
				m.PutAsArray(eachVal.Bool())
			case reflect.String:
				m.PutAsArray(eachVal.String())
			case reflect.Float32, reflect.Float64:
				m.PutAsArray(eachVal.Float())
			case reflect.Array, reflect.Slice:
				sJson := NewDJSON()
				sJson.SetAsArray()
				sJson.fromFieldsValue(eachVal, downDepthWW(tags)...)
				m.PutAsArray(sJson)
			case reflect.Struct, reflect.Map:
				switch eachType.String() {
				case "null.String":
					m.PutAsArray(eachVal.FieldByName("String").String())
				case "null.Bool":
					m.PutAsArray(eachVal.FieldByName("Bool").Bool())
				case "null.Float32":
					m.PutAsArray(eachVal.FieldByName("Float32").Float())
				case "null.Float64":
					m.PutAsArray(eachVal.FieldByName("Float64").Float())
				case "null.Int":
					m.PutAsArray(eachVal.FieldByName("Int").Int())
				case "null.Int8":
					m.PutAsArray(eachVal.FieldByName("Int8").Int())
				case "null.Int16":
					m.PutAsArray(eachVal.FieldByName("Int16").Int())
				case "null.Int32":
					m.PutAsArray(eachVal.FieldByName("Int32").Int())
				case "null.Int64":
					m.PutAsArray(eachVal.FieldByName("Int64").Int())
				case "null.Uint":
					m.PutAsArray(eachVal.FieldByName("Uint").Uint())
				case "null.Uint8":
					m.PutAsArray(eachVal.FieldByName("Uint8").Uint())
				case "null.Uint16":
					m.PutAsArray(eachVal.FieldByName("Uint16").Uint())
				case "null.Uint32":
					m.PutAsArray(eachVal.FieldByName("Uint32").Uint())
				case "null.Uint64":
					m.PutAsArray(eachVal.FieldByName("Uint64").Uint())
				default:
					sJson := NewDJSON()
					sJson.SetAsObject()
					sJson.fromFieldsValue(eachVal, downDepthWW(tags)...)
					m.PutAsArray(sJson)
				}
			default:
				m.PutAsArray(nil)
			}

		}

	} else if kind == reflect.Struct {

		for i := 0; i < val.NumField(); i++ {
			eachVal := val.Field(i)
			eachType := val.Type().Field(i)
			eachTag := eachType.Tag.Get("json")

			if len(tags) > 0 && !inTags(eachTag, tags...) {
				continue
			}

			eachKind := eachType.Type.Kind()

			if eachKind == reflect.Struct || eachKind == reflect.Map {

				switch eachType.Type.String() {
				case "null.String":
					m.Put(eachTag, eachVal.FieldByName("String").String())
				case "null.Bool":
					m.Put(eachTag, eachVal.FieldByName("Bool").Bool())
				case "null.Float32":
					m.Put(eachTag, eachVal.FieldByName("Float32").Float())
				case "null.Float64":
					m.Put(eachTag, eachVal.FieldByName("Float64").Float())
				case "null.Int":
					m.Put(eachTag, eachVal.FieldByName("Int").Int())
				case "null.Int8":
					m.Put(eachTag, eachVal.FieldByName("Int8").Int())
				case "null.Int16":
					m.Put(eachTag, eachVal.FieldByName("Int16").Int())
				case "null.Int32":
					m.Put(eachTag, eachVal.FieldByName("Int32").Int())
				case "null.Int64":
					m.Put(eachTag, eachVal.FieldByName("Int64").Int())
				case "null.Uint":
					m.Put(eachTag, eachVal.FieldByName("Uint").Uint())
				case "null.Uint8":
					m.Put(eachTag, eachVal.FieldByName("Uint8").Uint())
				case "null.Uint16":
					m.Put(eachTag, eachVal.FieldByName("Uint16").Uint())
				case "null.Uint32":
					m.Put(eachTag, eachVal.FieldByName("Uint32").Uint())
				case "null.Uint64":
					m.Put(eachTag, eachVal.FieldByName("Uint64").Uint())
				default:
					sJson := NewDJSON()
					sJson.SetAsObject()
					sJson.fromFieldsValue(eachVal, downDepthWW(tags)...)
					m.Put(eachTag, sJson)
				}
			} else if eachKind == reflect.Array || eachKind == reflect.Slice {

				sJson := NewDJSON()
				sJson.SetAsArray()
				sJson.fromFieldsValue(eachVal, downDepthWW(tags)...)
				m.Put(eachTag, sJson)

			} else {

				switch eachType.Type.String() {
				case "int", "int8", "int16", "int32", "int64":
					m.Put(eachTag, eachVal.Int())
				case "uint", "uint8", "uint16", "uint32", "uint64":
					m.Put(eachTag, eachVal.Uint())
				case "float32", "float64":
					m.Put(eachTag, eachVal.Float())
				case "string":
					m.Put(eachTag, eachVal.String())
				case "bool":
					m.Put(eachTag, eachVal.Bool())
				}
			}
		}
	} else if kind == reflect.Map {

		for _, e := range val.MapKeys() {
			eachKey, ok := e.Interface().(string)
			if !ok {
				continue
			}

			if len(tags) > 0 && !inTags(eachKey, tags...) {
				continue
			}

			eachVal := val.MapIndex(e)

			switch t := eachVal.Interface().(type) {
			case int:
				m.Put(eachKey, t)
			case int8:
				m.Put(eachKey, t)
			case int16:
				m.Put(eachKey, t)
			case int32:
				m.Put(eachKey, t)
			case int64:
				m.Put(eachKey, t)
			case uint:
				m.Put(eachKey, t)
			case uint8:
				m.Put(eachKey, t)
			case uint16:
				m.Put(eachKey, t)
			case uint32:
				m.Put(eachKey, t)
			case uint64:
				m.Put(eachKey, t)
			case float32:
				m.Put(eachKey, t)
			case float64:
				m.Put(eachKey, t)
			case string:
				m.Put(eachKey, t)
			case bool:
				m.Put(eachKey, t)
			case nil:
				m.Put(eachKey, t)
			case null.String:
				m.Put(eachKey, t.String)
			case null.Bool:
				m.Put(eachKey, t.Bool)
			case null.Int:
				m.Put(eachKey, t.Int)
			case null.Int8:
				m.Put(eachKey, t.Int8)
			case null.Int16:
				m.Put(eachKey, t.Int16)
			case null.Int32:
				m.Put(eachKey, t.Int32)
			case null.Int64:
				m.Put(eachKey, t.Int64)
			case null.Uint:
				m.Put(eachKey, t.Uint)
			case null.Uint8:
				m.Put(eachKey, t.Uint8)
			case null.Uint16:
				m.Put(eachKey, t.Uint16)
			case null.Uint32:
				m.Put(eachKey, t.Uint32)
			case null.Uint64:
				m.Put(eachKey, t.Uint64)
			case null.Float32:
				m.Put(eachKey, t.Float32)
			case null.Float64:
				m.Put(eachKey, t.Float64)
			default:

				skind := reflect.ValueOf(t).Type().Kind()

				if skind == reflect.Struct || skind == reflect.Map {
					sJson := NewDJSON()
					sJson.SetAsObject()
					sJson.FromFields(t, downDepthWW(tags)...)
					m.Put(eachKey, sJson)
				}

			}

		}
	}
}

func (m *DJSON) FromFields(st interface{}, tags ...string) *DJSON {
	baseValue := reflect.ValueOf(st)

	kind := baseValue.Type().Kind()

	if kind == reflect.Array || kind == reflect.Slice {

		m.SetAsArray()
		m.fromFieldsValue(baseValue, tags...)

	} else if kind == reflect.Struct || kind == reflect.Map {

		m.SetAsObject()
		m.fromFieldsValue(baseValue, tags...)

	}

	return m
}

func downDepthWW(tags []string) []string {
	tags2 := make([]string, 0)
	for idx := range tags {
		tmp := strings.Split(tags[idx], ".")
		tmp2 := strings.Join(tmp[1:], ".")
		if tmp2 != "" {
			tags2 = append(tags2, tmp2)
		} else {
			tags2 = append(tags2, "")
		}
	}

	return tags2
}

func inTags(idv string, tags ...string) bool {
	for idx := range tags {
		tmp := strings.Split(tags[idx], ".")
		if tmp[0] == idv {
			return true
		}
	}

	return false
}

func (m *DJSON) doSort(isAsc bool, k ...interface{}) bool {
	var tArray *DA

	if len(k) == 0 {
		if m.JsonType == JSON_ARRAY {
			tArray = m.Array
		}
	}

	if len(k) > 0 {

		if m.JsonType == JSON_OBJECT {
			if key, ok := k[0].(string); ok {
				if da, ok := m.Object.GetAsArray(key); ok {
					tArray = da
				}
			}
		} else if m.JsonType == JSON_ARRAY {
			if idx, ok := k[0].(int); ok {
				if da, ok := m.Array.GetAsArray(idx); ok {
					tArray = da
				}
			}
		}
	}

	if tArray != nil {
		return tArray.Sort(isAsc)
	} else {
		return false
	}
}

func (m *DJSON) SortAsc(k ...interface{}) bool {
	return m.doSort(true, k...)
}

func (m *DJSON) SortDesc(k ...interface{}) bool {
	return m.doSort(false, k...)
}

func (m *DJSON) SortObjectArray(isAsc bool, key string) bool {
	if m.JsonType != JSON_ARRAY {
		return false
	}

	return m.Array.SortObject(isAsc, key)
}

func (m *DJSON) SortObjectArrayAsc(key string) bool {
	return m.SortObjectArray(true, key)
}

func (m *DJSON) SortObjectArrayDesc(key string) bool {
	return m.SortObjectArray(false, key)
}

func (m *DJSON) Equal(t *DJSON) bool {
	if m.JsonType != t.JsonType {
		return false
	}

	switch m.JsonType {
	case JSON_NULL:
		return true
	case JSON_BOOL:
		return m.Bool == t.Bool
	case JSON_INT:
		return m.Int == t.Int
	case JSON_FLOAT:
		return m.Float == t.Float
	case JSON_STRING:
		return m.String == t.String
	case JSON_OBJECT:
		return m.Object.Equal(t.Object)
	case JSON_ARRAY:
		return m.Array.Equal(t.Array)
	}

	return false
}

func (m *DJSON) Clone() *DJSON {
	t := NewDJSON(m.JsonType)

	switch m.JsonType {
	case JSON_NULL:
	case JSON_BOOL:
		t.Bool = m.Bool
	case JSON_INT:
		t.Int = m.Int
	case JSON_FLOAT:
		t.Float = m.Float
	case JSON_STRING:
		t.String = m.String
	case JSON_OBJECT:
		t.Object = m.Object.Clone()
	case JSON_ARRAY:
		t.Array = m.Array.Clone()
	}

	return t
}

func (m *DJSON) HasKeys(k ...interface{}) bool {
	for i := range k {
		if !m.HasKey(k[i]) {
			return false
		}
	}

	return true
}

func (m *DJSON) GetKeys(k ...interface{}) []string {
	rk := make([]string, 0)

	if IsEmptyArg(k) {
		if m.JsonType != JSON_OBJECT {
			return rk
		}

		for k := range m.Object.Map {
			rk = append(rk, k)
		}

		return rk
	}

	if t, ok := m.GetAsObject(k[0]); ok {
		return t.GetKeys()
	}

	return rk
}

func (m *DJSON) Find(key string, val string) *DJSON {
	if key == "" || m.JsonType != JSON_ARRAY {
		return nil
	}

	for i := 0; i < m.Length(); i++ {
		each, ok := m.GetAsObject(i)
		if !ok {
			continue
		}

		if each.GetAsString(key) == val {
			return each
		}
	}

	return nil
}

func (m *DJSON) Append(arrJson *DJSON) *DJSON {
	if arrJson == nil || m.JsonType != JSON_ARRAY || !arrJson.IsArray() {
		return m
	}

	for i := 0; i < arrJson.Length(); i++ {
		m.PutAsArray(arrJson.Array.Element[i])
	}

	return m
}

func IsEmptyArg(key []interface{}) bool {
	return len(key) == 0 || (len(key) == 1 && key[0] == "")
}
