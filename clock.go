// Package clock implements a library for mocking time.
//
// All methods are safe for concurrent use.
package clock

import (
	"time"
)

// Clock represents an interface to the functions in the standard time package.
type Clock interface {
	After(d time.Duration) <-chan time.Time
	AfterFunc(d time.Duration, f func()) *Timer
	NewTicker(d time.Duration) *Ticker
	NewTimer(d time.Duration) *Timer
	Now() time.Time
	Since(t time.Time) time.Duration
	Sleep(d time.Duration)
	Tick(d time.Duration) <-chan time.Time
	Until(t time.Time) time.Duration
}

type clock struct{}

var realtime = clock{}

// Realtime returns the standard real-time Clock.
func Realtime() Clock {
	return realtime
}

func (clock) After(d time.Duration) <-chan time.Time {
	return time.After(d)
}

func (clock) AfterFunc(d time.Duration, f func()) *Timer {
	return &Timer{timer: time.AfterFunc(d, f)}
}

func (clock) NewTicker(d time.Duration) *Ticker {
	t := time.NewTicker(d)
	return &Ticker{
		C:      t.C,
		ticker: t,
	}
}

func (clock) NewTimer(d time.Duration) *Timer {
	t := time.NewTimer(d)
	return &Timer{
		C:     t.C,
		timer: t,
	}
}

func (clock) Now() time.Time {
	return time.Now()
}

func (clock) Since(t time.Time) time.Duration {
	return time.Since(t)
}

func (clock) Sleep(d time.Duration) {
	time.Sleep(d)
}

func (clock) Tick(d time.Duration) <-chan time.Time {
	// Using time.Tick would trigger a vet tool warning.
	if d <= 0 {
		return nil
	}
	return time.NewTicker(d).C
}

func (clock) Until(t time.Time) time.Duration {
	return time.Until(t)
}
