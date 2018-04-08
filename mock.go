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
// While suspended, all attempts to use the API will block.
//
// To increase predictability, all Mock methods acquire
// and release the Mutex only once during their execution.
type Mock struct {
	sync.Mutex
	now time.Time
	mockTimers
}

// New returns a new mocked Clock with current time set to now.
//
// Use Realtime to get the standard real-time Clock.
func New(now time.Time) *Mock {
	return &Mock{
		now:        now,
		mockTimers: &timerHeap{},
	}
}

// Add advances the current time by d and fires all expires timers.
//
// To increase predictability and speed, Tickers are ticked only once per call.
func (m *Mock) Add(d time.Duration) {
	m.Lock()
	defer m.Unlock()
	m.set(m.now.Add(d))
}

// Set sets the current time to now and fires all expired timers.
//
// To increase predictability and speed, Tickers are ticked only once per call.
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
		d := t.fire()
		if d == 0 {
			// Timers are always stopped.
			m.stop(t)
		} else {
			// Ticker's next deadline is set to the first tick after the new now.
			dd := (now.Sub(m.now)/d + 1) * d
			t.deadline = m.now.Add(dd)
			m.reset(t)
		}
	}
}

// Now returns the current mocked time.
func (m *Mock) Now() time.Time {
	m.Lock()
	defer m.Unlock()
	return m.now
}

// Since returns the time elapsed since t.
func (m *Mock) Since(t time.Time) time.Duration {
	m.Lock()
	defer m.Unlock()
	return m.now.Sub(t)
}

// Until returns the duration until t.
func (m *Mock) Until(t time.Time) time.Duration {
	m.Lock()
	defer m.Unlock()
	return t.Sub(m.now)
}
