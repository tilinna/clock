package clock

import (
	"sync"
	"time"
)

type mockTimers interface {
	start(t *mockTimer)
	stop(t *mockTimer)
	reset(t *mockTimer)
	next() *mockTimer
}

// Mock represents a Clock that only moves with Add() or Set().
//
// The clock can be suspended with Lock and resumed with Unlock.
type Mock struct {
	sync.Mutex
	now time.Time
	mockTimers
}

// New returns a new mocked Clock with current time set to now.
//
// Use clock.Realtime() to get the standard real-time Clock.
func New(now time.Time) *Mock {
	return &Mock{
		now:        now,
		mockTimers: &timerHeap{},
	}
}

func (m *Mock) Add(d time.Duration) {
	m.Lock()
	defer m.Unlock()
	m.set(m.now.Add(d))
}

func (m *Mock) Set(now time.Time) {
	m.Lock()
	defer m.Unlock()
	m.set(now)
}

func (m *Mock) set(now time.Time) {
	for {
		t := m.next()
		if t == nil || t.deadline.After(now) {
			m.now = now
			return
		}
		m.now = t.deadline
		d := t.timeoutFunc()
		if d == 0 {
			// Timer is always stopped
			m.stop(t)
		} else {
			// Ticker's next deadline is set to the first tick after the new now.
			dd := (now.Sub(m.now)/d + 1) * d
			t.deadline = m.now.Add(dd)
			m.reset(t)
		}
	}
}

func (m *Mock) Now() time.Time {
	m.Lock()
	defer m.Unlock()
	return m.now
}

func (m *Mock) Since(t time.Time) time.Duration {
	m.Lock()
	defer m.Unlock()
	return m.now.Sub(t)
}

func (m *Mock) Until(t time.Time) time.Duration {
	m.Lock()
	defer m.Unlock()
	return t.Sub(m.now)
}
