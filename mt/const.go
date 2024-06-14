package mt

type numbers interface {
	uint | uint8 | uint16 | uint32 | uint64 | int | int8 | int16 | int32 | int64 | float32 | float64
}

type integers interface {
	uint | uint8 | uint16 | uint32 | uint64 | int | int8 | int16 | int32 | int64
}

const (
	JS_MAX_INT int64 = 9007199254740991
)

const (
	ONE_WEEK_MS = 604800_000
	ONE_DAY_MS  = 86400_000
	ONE_HOUR_MS = 3600_000

	ONE_WEEK_SEC = 604800
	ONE_DAY_SEC  = 86400
	ONE_HOUR_SEC = 3600
)
