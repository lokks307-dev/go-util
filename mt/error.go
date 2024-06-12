package mt

import (
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

func buildLogString(fileName string, line int, ok bool, msgIn []interface{}) string {
	msgList := make([]string, 0)
	for _, msg := range msgIn {
		switch v := msg.(type) {
		case error:
			msgList = append(msgList, v.Error())
		case []string:
			msgList = append(msgList, fmt.Sprintf("[%s]", strings.Join(v, ",")))
		case []interface{}:
			vx := make([]string, 0)
			for _, xx := range v {
				vx = append(vx, fmt.Sprintf("%v", xx))
			}
			msgList = append(msgList, fmt.Sprintf("[%s]", strings.Join(vx, ",")))
		default:
			msgList = append(msgList, fmt.Sprintf("%v", v))
		}
	}

	if ok {
		return fmt.Sprintf("%s (%s:%d)", strings.Join(msgList, " "), filepath.Base(fileName), line)
	}

	return strings.Join(msgList, " ")
}

func Error(msgIn ...interface{}) error {
	if len(msgIn) == 1 && msgIn[0] == nil {
		return nil
	}

	_, filename, line, ok := runtime.Caller(1)
	return errors.New(buildLogString(filename, line, ok, msgIn))
}

func ErrorOut(msgIn ...interface{}) {
	_, filename, line, ok := runtime.Caller(1)
	msg := buildLogString(filename, line, ok, msgIn)

	logrus.Error(msg)
}

func Trace(msgIn ...interface{}) string {
	_, filename, line, ok := runtime.Caller(1)
	return buildLogString(filename, line, ok, msgIn)
}

type LogMap struct {
	Log map[string]string
}

func (m *LogMap) Error(msgIn ...interface{}) error {
	prelog := m.ToString()
	if prelog == "" {
		return Error(msgIn...)
	}

	vv := make([]interface{}, 0)
	vv = append(vv, prelog)
	vv = append(vv, msgIn...)

	return Error(vv...)
}

func (m *LogMap) Trace(msgIn ...interface{}) string {
	prelog := m.ToString()
	if prelog == "" {
		return Trace(msgIn...)
	}

	vv := make([]interface{}, 0)
	vv = append(vv, prelog)
	vv = append(vv, msgIn...)

	return Trace(vv...)
}

func (m *LogMap) With(vv ...interface{}) *LogMap {
	if m.Log == nil {
		m.Log = make(map[string]string)
	}

	aLogMap := With(vv...)

	for k, v := range aLogMap.Log {
		m.Log[k] = v
	}

	return m
}

func (m *LogMap) ToString() string {
	logList := make([]string, 0)
	for k, v := range m.Log {
		logList = append(logList, fmt.Sprintf("%s=%s", k, v))
	}
	return strings.Join(logList, ",")
}

func With(vv ...interface{}) *LogMap {

	var ret LogMap
	ret.Log = make(map[string]string)

	lenVV := len(vv)

	if lenVV%2 == 1 {
		lenVV--
	}

	if lenVV == 0 {
		return &ret
	}

	for i := 0; i < lenVV; i += 2 {
		k := fmt.Sprintf("%s", vv[i])
		if k != "" {
			switch x := vv[i+1].(type) {
			case error:
				ret.Log[k] = x.Error()
			case []string:
				ret.Log[k] = fmt.Sprintf("[%s]", strings.Join(x, ","))
			case []interface{}:
				vx := make([]string, 0)
				for _, xx := range x {
					vx = append(vx, fmt.Sprintf("%v", xx))
				}
				ret.Log[k] = fmt.Sprintf("[%s]", strings.Join(vx, ","))
			default:
				ret.Log[k] = fmt.Sprintf("%v", vv[i+1])
			}

		}
	}

	return &ret
}
