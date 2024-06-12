package event

import (
	"fmt"
	"sync"

	"github.com/lokks307/djson/v2"
	"github.com/sirupsen/logrus"
)

const mNameEvt = "[EvtMan] "

type AEvent struct {
	Type     string
	Data     interface{}
	DataInts []int
	DataStrs []string
	DataJson *djson.JSON
}

type EvtStage uint8

const (
	STAGE_NULL EvtStage = iota
	STAGE_INIT
	STAGE_READY
	STAGE_STOP
)

type EventManager struct {
	ssFuncTable map[string]func(ae AEvent) // skey , function
	ssNameTable map[string][]string        // event type, [skey, skey, ...]
	cancelChan  chan bool
	stage       EvtStage
	skeyCount   int
	combiner    []*Combiner
	ssMutex     sync.RWMutex
}

var Bus chan AEvent
var Manager *EventManager

func init() {

	Bus = make(chan AEvent, 100)

	Manager = &EventManager{
		ssFuncTable: make(map[string]func(ae AEvent)),
		ssNameTable: make(map[string][]string),
		cancelChan:  make(chan bool, 2),
		stage:       STAGE_INIT,
		combiner:    make([]*Combiner, 0),
	}
}

func (m *EventManager) AddCombiner(cm *Combiner) {
	m.combiner = append(m.combiner, cm)
}

func (m *EventManager) Subscribe(funcp func(ae AEvent), etype ...string) string {

	m.skeyCount++

	skey := fmt.Sprintf("%06d", m.skeyCount)

	m.UpdateSubscription(skey, funcp, etype...)

	return skey

}

func (m *EventManager) UpdateSubscription(skey string, funcp func(ae AEvent), etype ...string) {
	m.ssMutex.Lock()
	defer m.ssMutex.Unlock()

	m.ssFuncTable[skey] = funcp

	for _, eachetype := range etype {
		skeyList := m.ssNameTable[eachetype]
		found := false
		for _, eachskey := range skeyList {
			if eachskey == skey {
				found = true
			}
		}

		if !found {
			m.ssNameTable[eachetype] = append(m.ssNameTable[eachetype], skey)
		}
	}
}

func (m *EventManager) RemoveSubscribe(skey string) {
	m.ssMutex.Lock()
	defer m.ssMutex.Unlock()

	delete(m.ssFuncTable, skey)

	for eachetype := range m.ssNameTable {
		skeyList := m.ssNameTable[eachetype]

		foundIdx := -1
		for idx, eachskey := range skeyList {
			if eachskey == skey {
				foundIdx = idx
				break
			}
		}

		if foundIdx >= 0 {

			if len(skeyList) == 0 {
				m.ssNameTable[eachetype] = make([]string, 0)
			} else {
				skeyList[foundIdx] = skeyList[len(skeyList)-1]
				m.ssNameTable[eachetype] = skeyList[0 : len(skeyList)-1]
			}
		}
	}

}

func (m *EventManager) Run() error {

	if m.stage == STAGE_NULL || m.stage == STAGE_READY {
		logrus.Error(mNameEvt, ERROR_EVT_FAIL_RUN)
		return ERROR_EVT_FAIL_RUN
	}

	go func() {

		logrus.Info(mNameEvt, "started")
		m.stage = STAGE_READY

		defer func() {
			for len(m.cancelChan) > 0 {
				<-m.cancelChan
			}
		}()

		for {
			select {
			case <-m.cancelChan:
				return
			case oneEvent := <-Bus:
				logrus.Trace(mNameEvt, "new event type=", oneEvent.Type)

				for idx := range m.combiner {
					m.combiner[idx].Listen(&oneEvent)
				}

				m.ssMutex.RLock()
				if skeyList, ok := m.ssNameTable[oneEvent.Type]; ok {
					for _, skey := range skeyList {
						if funcP, ok := m.ssFuncTable[skey]; ok {
							go funcP(oneEvent)
						}
					}
				}
				m.ssMutex.RUnlock()

			}
		}
	}()

	return nil
}

func (m *EventManager) Stop() {
	m.cancelChan <- true
}

func (m *EventManager) On(etype string, funcp func(ae AEvent)) string {
	return m.Subscribe(funcp, etype)
}

func On(etype string, funcp func(ae AEvent)) string {
	return Manager.Subscribe(funcp, etype)
}
