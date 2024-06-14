package mt

import "github.com/lokks307/djson/v2"

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
