package clock

import "time"

type Timer struct {
	C     <-chan time.Time
	timer *time.Timer
	*mockTimer
}

func (m *Mock) After(d time.Duration) <-chan time.Time {
	return m.NewTimer(d).C
}

func (m *Mock) AfterFunc(d time.Duration, f func()) *Timer {
	m.Lock()
	defer m.Unlock()
	return m.newTimerFunc(m.now.Add(d), f)
}

func (m *Mock) NewTimer(d time.Duration) *Timer {
	m.Lock()
	defer m.Unlock()
	return m.newTimerFunc(m.now.Add(d), nil)
}

func (m *Mock) Sleep(d time.Duration) {
	<-m.After(d)
}

func (m *Mock) newTimerFunc(deadline time.Time, afterFunc func()) *Timer {
	t := &Timer{
		mockTimer: &mockTimer{
			deadline: deadline,
			mock:     m,
		},
	}
	if afterFunc != nil {
		t.timeoutFunc = func() time.Duration {
			go afterFunc()
			return 0
		}
	} else {
		c := make(chan time.Time, 1)
		t.C = c
		t.timeoutFunc = func() time.Duration {
			select {
			case c <- m.now:
			default:
			}
			return 0
		}
	}
	if !t.deadline.After(m.now) {
		t.timeoutFunc()
	} else {
		m.start(t.mockTimer)
	}
	return t
}

func (t *Timer) Stop() bool {
	if t.timer != nil {
		return t.timer.Stop()
	}
	t.mock.Lock()
	defer t.mock.Unlock()
	wasActive := !t.mockTimer.stopped()
	t.mock.stop(t.mockTimer)
	return wasActive
}

func (t *Timer) Reset(d time.Duration) bool {
	if t.timer != nil {
		return t.timer.Reset(d)
	}
	t.mock.Lock()
	defer t.mock.Unlock()
	wasActive := !t.mockTimer.stopped()
	t.deadline = t.mock.now.Add(d)
	if !t.deadline.After(t.mock.now) {
		t.timeoutFunc()
		t.mock.stop(t.mockTimer)
	} else {
		t.mock.reset(t.mockTimer)
	}
	return wasActive
}
