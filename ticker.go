package clock

import (
	"errors"
	"time"
)

type Ticker struct {
	C      <-chan time.Time
	ticker *time.Ticker
	*mockTimer
}

func (m *Mock) NewTicker(d time.Duration) *Ticker {
	m.Lock()
	defer m.Unlock()
	if d <= 0 {
		panic(errors.New("non-positive interval for NewTicker"))
	}
	return m.newTicker(d)
}

func (m *Mock) Tick(d time.Duration) <-chan time.Time {
	m.Lock()
	defer m.Unlock()
	if d <= 0 {
		return nil
	}
	return m.newTicker(d).C
}

func (m *Mock) newTicker(d time.Duration) *Ticker {
	c := make(chan time.Time, 1)
	t := &Ticker{
		C: c,
		mockTimer: &mockTimer{
			deadline: m.now.Add(d),
			mock:     m,
		},
	}
	t.timeoutFunc = func() time.Duration {
		select {
		case c <- m.now:
		default:
		}
		return d
	}
	m.start(t.mockTimer)
	return t
}

func (t *Ticker) Stop() {
	if t.ticker != nil {
		t.ticker.Stop()
		return
	}
	t.mock.Lock()
	defer t.mock.Unlock()
	t.mock.stop(t.mockTimer)
}
