package mt

import (
	"time"

	"github.com/lokks307/djson/v2"
	"github.com/volatiletech/null/v8"
)

type NullMap struct {
	Valid bool
	Map   map[string]string
}

func (m *NullMap) SetValid(mp map[string]string) {
	if mp != nil {
		m.Valid = true
		m.Map = mp
	}
}

func (m *NullMap) ToString() string {
	r := djson.New()

	if m.Valid {
		for k, v := range m.Map {
			r.Put(k, v)
		}
	}

	return r.ToString()
}

func NewNullMap(mp map[string]string, valid bool) NullMap {
	if mp != nil {
		return NullMap{
			Valid: valid,
			Map:   mp,
		}
	} else {
		return NullMap{
			Valid: false,
			Map:   nil,
		}
	}
}

func NullMapFrom(mp map[string]string) NullMap {
	return NewNullMap(mp, true)
}

func NullMapFromJson(mp *djson.JSON) NullMap {
	m := make(map[string]string)

	if mp != nil && mp.IsObject() {
		klist := mp.GetKeys()
		for _, k := range klist {
			m[k] = mp.String(k)
		}
	}

	return NewNullMap(m, true)
}

type NullStringSlice struct {
	Valid bool
	Slice []string
}

func (m *NullStringSlice) SetValid(ss []string) {
	if ss != nil {
		m.Valid = true
		m.Slice = ss
	}
}

func NewNullStringSlice(ss []string, valid bool) NullStringSlice {
	if ss != nil {
		return NullStringSlice{
			Valid: valid,
			Slice: ss,
		}
	} else {
		return NullStringSlice{
			Valid: false,
			Slice: nil,
		}
	}
}

func NullStringSliceFrom(ss []string) NullStringSlice {
	return NewNullStringSlice(ss, true)
}

func NullStringSliceFromJson(ss *djson.JSON) NullStringSlice {
	s := make([]string, 0)
	if ss != nil && ss.IsArray() {
		ss.Seek()
		for ss.Next() {
			se := ss.Scan()
			s = append(s, se.String())
		}
	}

	return NewNullStringSlice(s, true)
}

func GetTimeUnix(s null.Time) int64 {
	if s.Valid {
		return s.Time.Unix()
	}

	return 0
}

func GetBool(s null.Bool) bool {
	if s.Valid {
		return s.Bool
	}

	return false
}

func GetString(s null.String) string {
	if s.Valid {
		return s.String
	}

	return ""
}

func GetFloat(s interface{}) float64 {
	switch t := s.(type) {
	case null.Float32:
		if t.Valid {
			return float64(t.Float32)
		}
	case null.Float64:
		if t.Valid {
			return t.Float64
		}
	}

	return 0
}

func GetInt(s interface{}) int64 {
	switch t := s.(type) {
	case null.Uint:
		if t.Valid {
			return int64(t.Uint)
		}
	case null.Uint8:
		if t.Valid {
			return int64(t.Uint8)
		}
	case null.Uint16:
		if t.Valid {
			return int64(t.Uint16)
		}
	case null.Uint32:
		if t.Valid {
			return int64(t.Uint32)
		}
	case null.Uint64:
		if t.Valid {
			return int64(t.Uint64)
		}
	case null.Int:
		if t.Valid {
			return int64(t.Int)
		}
	case null.Int8:
		if t.Valid {
			return int64(t.Int8)
		}
	case null.Int16:
		if t.Valid {
			return int64(t.Int16)
		}
	case null.Int32:
		if t.Valid {
			return int64(t.Int32)
		}
	case null.Int64:
		if t.Valid {
			return t.Int64
		}
	}

	return 0
}

func GetUpdatedUnix(createdAt time.Time, updatedAt null.Time) int64 {
	if updatedAt.Valid {
		return updatedAt.Time.Unix()
	}

	return createdAt.Unix()
}
