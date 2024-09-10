package mt

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

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

func ToSqlValStr(sa ...any) string {
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

func DoInTransaction(dbCon *sql.DB, fn func(tx *sql.Tx) error) error {
	ctx := context.Background()

	tx, err := dbCon.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if err = fn(tx); err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func DoInTransactionWithErrorCode(dbCon *sql.DB, fn func(tx *sql.Tx) (error, string)) (error, string) {
	ctx := context.Background()

	tx, err := dbCon.BeginTx(ctx, nil)
	if err != nil {
		return err, ""
	}

	if err, code := fn(tx); err != nil {
		tx.Rollback()
		return err, code
	}

	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return err, ""
	}

	return nil, ""
}

func GetStringCols(size int) string {
	// ex. size = 3
	// => stringsCols = "(?, ?, ?)"

	stringsCols := ""
	for i := 0; i < size; i++ {
		if i != size-1 {
			stringsCols += "?,"
			continue
		}
		stringsCols += "?"
	}

	return "(" + stringsCols + ")"
}

func GetStringColsWithNoParenthesis(columns ...string) string {
	// ex. columns = a, b, c
	// => stringsCols = "a, b, c"

	stringsCols := ""
	for idx, col := range columns {
		if idx == len(columns)-1 {
			stringsCols += col
			break
		}
		stringsCols += col + ","
	}

	return stringsCols
}

func OrderByString(columnName, order string) string {
	query := fmt.Sprintf("%s ASC", columnName)
	switch strings.ToUpper(strings.Trim(order, " ")) {
	case "DESC":
		query = fmt.Sprintf("%s DESC", columnName)
	case "ASC":
		query = fmt.Sprintf("%s ASC", columnName)
	}

	return query
}

func CancatParams(in []any) string {
	ret := []string{}

	for _, v := range in {
		switch t := v.(type) {
		case int64, int32, int, int8, uint16:
			ret = append(ret, fmt.Sprintf("%d", t))
		}
	}

	return strings.Join(ret, ",")
}

func GetBulkInsertQuery(tableName string, cols []string, vals [][]any) string {
	if tableName == "" || len(cols) == 0 || len(vals) == 0 {
		return ""
	}

	colsStmt := "(" + strings.Join(cols, ",") + ")"

	valStrSlice := make([]string, 0)
	for i := range vals {
		valStrSlice = append(valStrSlice, ToSqlValStr(vals[i]...))
	}

	return fmt.Sprintf("INSERT IGNORE INTO %s %s VALUES %s", tableName, colsStmt, strings.Join(valStrSlice, ","))
}

func GetJoinQueryStmt(tableName string, cols ...string) string {
	if tableName == "" || (len(cols)%2) == 1 {
		return ""
	}

	ret := tableName
	if len(cols) > 0 {
		colStmt := make([]string, 0)

		for i := 0; i < len(cols); i += 2 {
			if cols[i+1] == "" || strings.EqualFold(cols[i+1], "NULL") {
				colStmt = append(colStmt, fmt.Sprintf("%s IS %s", cols[i], cols[i+1]))
			} else {
				colStmt = append(colStmt, fmt.Sprintf("%s = %s", cols[i], cols[i+1]))
			}
		}

		ret = fmt.Sprintf("%s ON %s", tableName, strings.Join(colStmt, " AND "))
	}

	return ret
}
