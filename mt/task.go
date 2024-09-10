package mt

import (
	"sync"
	"time"
)

type TaskTimeItem struct {
	Timestamp int64
	Result    any
}

type TaskTimeMap struct {
	TimeMapMutex sync.RWMutex
	TaskMap      map[string]TaskTimeItem
}

func NewTaskTimeMap() *TaskTimeMap {
	m := &TaskTimeMap{
		TaskMap: make(map[string]TaskTimeItem),
	}
	return m
}

func (m *TaskTimeMap) GetTaskWithinTime(taskKey, prefix string, wt int64) (any, bool) {
	m.TimeMapMutex.RLock()
	defer m.TimeMapMutex.RUnlock()

	rkey := prefix + "-" + taskKey
	bt, ok := m.TaskMap[rkey]

	isWt := false
	if ok {
		isWt = AbsInt64(time.Now().Unix()-bt.Timestamp) <= wt
	}

	if isWt {
		return bt.Result, true
	} else {
		return nil, false
	}
}

func (m *TaskTimeMap) SetTask(taskKey, prefix string, taskResult any) {
	m.TimeMapMutex.Lock()
	defer m.TimeMapMutex.Unlock()

	rkey := prefix + "-" + taskKey
	m.TaskMap[rkey] = TaskTimeItem{
		Timestamp: time.Now().Unix(),
		Result:    taskResult,
	}
}

func (m *TaskTimeMap) Remove(taskKey, prefix string) {
	m.TimeMapMutex.Lock()
	defer m.TimeMapMutex.Unlock()

	rkey := prefix + "-" + taskKey
	delete(m.TaskMap, rkey)
}
