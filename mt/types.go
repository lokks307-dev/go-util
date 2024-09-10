package mt

import (
	"fmt"
	"reflect"
	"strings"

	gov "github.com/asaskevich/govalidator"
)

func GetStringBase(v any) (string, bool) {
	if v == nil {
		return "nil", true
	}

	if IsInTypes(v, "string", "bool", "float32", "float64") {
		return fmt.Sprintf("%v", v), true
	}

	if IsInTypes(v, "int", "uint", "int8", "uint8", "int16", "uint16", "int32", "uint32", "int64", "uint64") {
		return fmt.Sprintf("%d", v), true
	}

	return "", false
}

func GetBoolBase(v any) (bool, bool) {
	if IsInTypes(v, "int", "uint", "int8", "uint8", "int16", "uint16", "int32", "uint32", "int64", "uint64") {
		intVal, _ := gov.ToInt(v)
		if intVal == 0 {
			return false, true
		}
	}

	if IsInTypes(v, "string") {
		if strVal, ok := v.(string); ok {
			if strings.EqualFold(strVal, "true") {
				return true, true
			} else if strings.EqualFold(strVal, "false") {
				return false, true
			}
		}
	}

	if IsInTypes(v, "bool") {
		if boolVal, ok := v.(bool); ok {
			return boolVal, true
		}
	}

	return false, false
}

func GetFloatBase(v any) (float64, bool) {
	if floatVal, err := gov.ToFloat(v); err != nil {
		return 0, false
	} else {
		return floatVal, true
	}
}

func GetIntBase(v any) (int64, bool) {
	if intVal, err := gov.ToInt(v); err != nil {
		return 0, false
	} else {
		return intVal, true
	}
}

func IsBaseType(v any) bool {
	return IsInTypes(v, "string", "bool", "int", "uint", "int8", "uint8", "int16", "uint16", "int32", "uint32", "int64", "uint64", "float32", "float64")
}

func IsIntType(v any) bool {
	return IsInTypes(v, "int", "uint", "int8", "uint8", "int16", "uint16", "int32", "uint32", "int64", "uint64")
}

func IsFloatType(v any) bool {
	return IsInTypes(v, "float32", "float64")
}

func IsBoolType(v any) bool {
	return IsInTypes(v, "bool")
}

func IsStringType(v any) bool {
	return IsInTypes(v, "string")
}

func IsSliceType(v any) bool {
	return IsInTypes(v,
		"[]string", "[]bool", "[]int", "[]uint", "[]int8", "[]uint8", "[]int16", "[]uint16", "[]int32", "[]uint32", "[]int64", "[]uint64", "[]float32", "[]float64",
		"[]null.String", "[]null.Bool", "[]null.Int", "[]null.Uint", "[]null.Int8", "[]null.Uint8", "[]null.Int16", "[]null.Uint16", "[]null.Int32", "[]null.Uint32", "[]null.Int64", "[]null.Uint64", "[]null.Float32", "[]null.Float64")
}

func IsInTypes(v any, types ...string) bool {
	var vTypeStr string
	if v == nil {
		vTypeStr = "nil"
	} else {
		vTypeStr = reflect.TypeOf(v).String()
	}

	for idx := range types {
		if vTypeStr == types[idx] {
			return true
		}
	}

	return false
}
