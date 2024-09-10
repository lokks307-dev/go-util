package mt

import (
	"time"

	"github.com/volatiletech/null/v8"
)

type NullInt interface {
	null.Int | null.Int8 | null.Int16 | null.Int32 | null.Int64 | null.Uint | null.Uint8 | null.Uint16 | null.Uint32 | null.Uint64
}

type NullFloat interface {
	null.Float32 | null.Float64
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

func GetFloat[T NullFloat](s T) float64 {
	r, ok := AnyToFloat64(s)
	if ok {
		return r
	}
	return 0.0
}

func GetInt[T NullInt](s T) int64 {
	r, ok := AnyToInt64(s)
	if ok {
		return r
	}
	return 0
}

func GetUpdatedUnix(createdAt time.Time, updatedAt null.Time) int64 {
	if updatedAt.Valid {
		return updatedAt.Time.Unix()
	}

	return createdAt.Unix()
}
