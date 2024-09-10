package mt

import (
	"time"

	"github.com/lokks307/djson/v2"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

// models.M (in sqlboiler) is equal to map[string]any

func OrderByCol(colName string, isAsc bool) qm.QueryMod {
	if isAsc {
		return qm.OrderBy(colName + " ASC")
	}
	return qm.OrderBy(colName + " DESC")
}

func assignWithEval[T any](m map[string]any, key string, v T, ev func(c any) bool) {
	if ev == nil {
		m[key] = v
	} else {
		if ev(v) {
			m[key] = v
		}
	}
}

func AppendBoilCols(m map[string]any, key string, v any, ev ...func(c any) bool) map[string]any {
	var evalfunc func(c any) bool
	if len(ev) == 1 {
		evalfunc = ev[0]
	}

	switch t := v.(type) {
	case int:
		assignWithEval(m, key, t, evalfunc)
	case int8:
		assignWithEval(m, key, t, evalfunc)
	case int16:
		assignWithEval(m, key, t, evalfunc)
	case int32:
		assignWithEval(m, key, t, evalfunc)
	case int64:
		assignWithEval(m, key, t, evalfunc)
	case uint:
		assignWithEval(m, key, t, evalfunc)
	case uint8:
		assignWithEval(m, key, t, evalfunc)
	case uint16:
		assignWithEval(m, key, t, evalfunc)
	case uint32:
		assignWithEval(m, key, t, evalfunc)
	case uint64:
		assignWithEval(m, key, t, evalfunc)
	case *int:
		if t != nil {
			assignWithEval(m, key, *t, evalfunc)
		}
	case *int8:
		if t != nil {
			assignWithEval(m, key, *t, evalfunc)
		}
	case *int16:
		if t != nil {
			assignWithEval(m, key, *t, evalfunc)
		}
	case *int32:
		if t != nil {
			assignWithEval(m, key, *t, evalfunc)
		}
	case *int64:
		if t != nil {
			assignWithEval(m, key, *t, evalfunc)
		}
	case *uint:
		if t != nil {
			assignWithEval(m, key, *t, evalfunc)
		}
	case *uint16:
		if t != nil {
			assignWithEval(m, key, *t, evalfunc)
		}
	case *uint32:
		if t != nil {
			assignWithEval(m, key, *t, evalfunc)
		}
	case *uint64:
		if t != nil {
			assignWithEval(m, key, *t, evalfunc)
		}
	case float32:
		assignWithEval(m, key, t, evalfunc)
	case *float32:
		if t != nil {
			assignWithEval(m, key, *t, evalfunc)
		}
	case float64:
		assignWithEval(m, key, t, evalfunc)
	case *float64:
		if t != nil {
			assignWithEval(m, key, *t, evalfunc)
		}
	case string:
		assignWithEval(m, key, t, evalfunc)
	case *string:
		if t != nil {
			assignWithEval(m, key, *t, evalfunc)
		}
	case bool:
		assignWithEval(m, key, t, evalfunc)
	case *bool:
		if t != nil {
			assignWithEval(m, key, *t, evalfunc)
		}
	case time.Time:
		assignWithEval(m, key, t, evalfunc)
	case *time.Time:
		if t != nil {
			assignWithEval(m, key, *t, evalfunc)
		}
	case null.Int:
		if t.Valid {
			assignWithEval(m, key, t.Int, evalfunc)
		}
	case null.Int8:
		if t.Valid {
			assignWithEval(m, key, t.Int8, evalfunc)
		}
	case null.Int16:
		if t.Valid {
			assignWithEval(m, key, t.Int16, evalfunc)
		}
	case null.Int32:
		if t.Valid {
			assignWithEval(m, key, t.Int32, evalfunc)
		}
	case null.Int64:
		if t.Valid {
			assignWithEval(m, key, t.Int64, evalfunc)
		}
	case null.Uint:
		if t.Valid {
			assignWithEval(m, key, t.Uint, evalfunc)
		}
	case null.Uint8:
		if t.Valid {
			assignWithEval(m, key, t.Uint8, evalfunc)
		}
	case null.Uint16:
		if t.Valid {
			assignWithEval(m, key, t.Uint16, evalfunc)
		}
	case null.Uint32:
		if t.Valid {
			assignWithEval(m, key, t.Uint32, evalfunc)
		}
	case null.Uint64:
		if t.Valid {
			assignWithEval(m, key, t.Uint64, evalfunc)
		}
	case null.Time:
		if t.Valid {
			assignWithEval(m, key, t.Time, evalfunc)
		}
	case *djson.JSON:
		if t != nil {
			assignWithEval(m, key, t.ToString(), evalfunc)
		}
	}

	return m
}

func UpdateValueAllowEmpty(vv any, opt *djson.JSON, key string) bool {
	if vv == nil || opt == nil || key == "" {
		return false
	}

	if !opt.HasKey(key) {
		return false
	}

	switch t := vv.(type) {
	case *string:
		v := opt.String(key)
		if *t != v {
			*t = v
			return true
		}

	case *null.String:
		v := opt.String(key)
		if t.String != v {
			t.SetValid(v)
			return true
		}

	case *int:
		v := int(opt.Int(key))
		if *t != v {
			*t = v
			return true
		}

	case *null.Int:
		v := int(opt.Int(key))
		if t.Int != v {
			t.SetValid(v)
			return true
		}
	case *int8:
		v := int8(opt.Int(key))
		if *t != v {
			*t = v
			return true
		}

	case *null.Int8:
		v := int8(opt.Int(key))
		if t.Int8 != v {
			t.SetValid(v)
			return true
		}
	case *int16:
		v := int16(opt.Int(key))
		if *t != v {
			*t = v
			return true
		}

	case *null.Int16:
		v := int16(opt.Int(key))
		if t.Int16 != v {
			t.SetValid(v)
			return true
		}
	case *int32:
		v := int32(opt.Int(key))
		if *t != v {
			*t = v
			return true
		}

	case *null.Int32:
		v := int32(opt.Int(key))
		if t.Int32 != v {
			t.SetValid(v)
			return true
		}
	case *int64:
		v := int64(opt.Int(key))
		if *t != v {
			*t = v
			return true
		}

	case *null.Int64:
		v := int64(opt.Int(key))
		if t.Int64 != v {
			t.SetValid(v)
			return true
		}

	case *bool:
		v := opt.Bool(key)
		if *t != v {
			*t = v
			return true
		}

	case *null.Bool:
		v := opt.Bool(key)
		if t.Bool != v {
			t.Bool = v
			return true
		}

	default:

	}

	return false
}

func UpdateValueNonEmpty(vv any, opt *djson.JSON, key string) bool {
	if vv == nil || opt == nil || key == "" {
		return false
	}

	if !opt.HasKey(key) {
		return false
	}

	switch vv.(type) {
	case *string, *null.String:
		if opt.String(key) == "" {
			return false
		}

	case *int, *int8, *int16, *int32, *int64, *null.Int, *null.Int8, *null.Int16, *null.Int32, *null.Int64:
		if opt.Int(key) == 0 {
			return false
		}

	case *bool, *null.Bool:
		if !opt.Bool(key) {
			return false
		}
	}

	return UpdateValueAllowEmpty(vv, opt, key)
}
