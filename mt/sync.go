package mt

import (
	"sync"

	"github.com/lokks307/micro-saas-phr/lib/flag"
)

type WaitGroup struct {
	wg sync.WaitGroup
	g  map[string]bool
	l  sync.Mutex
}

func (m *WaitGroup) Add(kk ...string) {
	m.l.Lock()
	defer m.l.Unlock()

	if m.g == nil {
		m.g = make(map[string]bool)
	}

	if len(kk) == 0 {
		m.wg.Add(1)
		return
	}

	for _, k := range kk {
		if _, ok := m.g[k]; !ok {
			m.g[k] = false
			m.wg.Add(1)
		}
	}
}

func (m *WaitGroup) IsDone(kk ...string) bool {
	m.l.Lock()
	defer m.l.Unlock()

	for _, k := range kk {
		if _, ok := m.g[k]; ok && !m.g[k] {
			if !m.g[k] {
				return false
			}
		}
	}

	return true
}

func (m *WaitGroup) Done(kk ...string) {
	m.l.Lock()
	defer m.l.Unlock()

	if m.g == nil {
		m.g = make(map[string]bool)
	}

	if len(kk) == 0 {
		m.wg.Done()
		return
	}

	for _, k := range kk {
		if _, ok := m.g[k]; ok && !m.g[k] {
			m.g[k] = true
			m.wg.Done()
		}
	}
}

func (m *WaitGroup) Wait() {
	m.wg.Wait()
}

type WaitMeet struct {
	wc chan bool
	f  *flag.AtomicBool
}

func NewWaitMeet() *WaitMeet {
	return &WaitMeet{
		wc: make(chan bool, 2),
		f:  flag.New(),
	}
}

func (m *WaitMeet) Meet() {
	if m.f.IsSet() {
		return
	}

	m.f.Set()
	select {
	case m.wc <- true:
	default:
	}
}

func (m *WaitMeet) Wait() {
	<-m.wc
}

func (m *WaitMeet) IsMeet() bool {
	return m.f.IsSet()
}
