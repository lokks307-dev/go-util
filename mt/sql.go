package mt

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

func ToSqlValStr(sa ...interface{}) string {
	contents := make([]string, 0)

	for _, ea := range sa {
		switch v := ea.(type) {
		case string:
			contents = append(contents, fmt.Sprintf("'%s'", MysqlRealEscapeString(v)))
		case int, int8, int16, int64, uint, uint8, uint16, uint32, uint64:
			contents = append(contents, fmt.Sprintf("%d", v))
		case float32, float64:
			contents = append(contents, fmt.Sprintf("%.4f", v))
		case bool:
			contents = append(contents, fmt.Sprintf("%t", v))
		case time.Time:
			contents = append(contents, fmt.Sprintf("'%s'", v.In(LocalLoc).Format("2006-01-02 15:04:05")))
		default:
			contents = append(contents, "NULL")
		}
	}

	return "(" + strings.Join(contents, ",") + ")"
}

func MysqlRealEscapeString(value string) string {
	var sb strings.Builder
	for i := 0; i < len(value); i++ {
		c := value[i]
		switch c {
		case '\\', 0, '\n', '\r', '\'', '"':
			sb.WriteByte('\\')
			sb.WriteByte(c)
		case '\032':
			sb.WriteByte('\\')
			sb.WriteByte('Z')
		default:
			sb.WriteByte(c)
		}
	}
	return sb.String()
}

func BuildParamSqlite(num int) string {
	paramRaw := make([]string, 0)

	for i := 0; i < num; i++ {
		paramRaw = append(paramRaw, "?")
	}

	return strings.Join(paramRaw, ",")
}

func BatchSliceTask[T integers | string](numOnce int, ss []T, taskFunc func([]T) error, abortOnError bool) {

	lenSlice := len(ss)
	for i := 0; i < lenSlice; i += numOnce {
		sliceBegin := i
		sliceEnd := i + numOnce
		if sliceEnd > lenSlice {
			sliceBegin = lenSlice
		}

		err := taskFunc(ss[sliceBegin:sliceEnd])
		if err != nil {
			if abortOnError {
				break
			}
		}
	}
}

func BatchSliceTaskWithTx[T integers | string](numOnce int, dbconn *sql.DB, ss []T, taskFunc func(*sql.Tx, []T) error) error {
	tx, err := dbconn.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}

	lenSlice := len(ss)
	for i := 0; i < lenSlice; i += numOnce {
		sliceBegin := i
		sliceEnd := i + numOnce
		if sliceEnd > lenSlice {
			sliceBegin = lenSlice
		}

		err := taskFunc(tx, ss[sliceBegin:sliceEnd])
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
