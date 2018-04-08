// Package timectx associates the clock with context.Context.
//
// All methods are safe for concurrent use.
package timectx

import (
	"context"
	"time"

	"github.com/tilinna/clock"
)

type clockKey struct{}

// WithClock returns a copy of parent in which the clock c is associated with.
//
// If the clock is an instance of *clock.Mock, all subsequent uses of the
// context in this package are automatically mocked, else the standard
// realtime clock is used.
func WithClock(parent context.Context, c clock.Clock) context.Context {
	return context.WithValue(parent, clockKey{}, c)
}

// Clock returns the clock associated with this context, or clock.Realtime().
func Clock(ctx context.Context) clock.Clock {
	if c, ok := ctx.Value(clockKey{}).(clock.Clock); ok {
		return c
	}
	return clock.Realtime()
}

// WithDeadline returns a copy of the parent context with the deadline adjusted
// to be no later than d.
//
// If the Clock(parent) returns a *clock.Mock, it is used to mock the deadline,
// else context.WithDeadline is called directly.
func WithDeadline(parent context.Context, d time.Time) (context.Context, context.CancelFunc) {
	if mock, ok := Clock(parent).(*clock.Mock); ok {
		return mock.DeadlineContext(parent, d)
	}
	return context.WithDeadline(parent, d)
}

// WithTimeout returns WithDeadline(parent, Clock(parent).Now().Add(timeout)).
//
// If the Clock(parent) returns a *clock.Mock, it is used to mock the deadline,
// else context.WithTimeout is called directly.
func WithTimeout(parent context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	if mock, ok := Clock(parent).(*clock.Mock); ok {
		return mock.TimeoutContext(parent, timeout)
	}
	return context.WithTimeout(parent, timeout)
}

// After is a convenience wrapper for Clock(ctx).After.
func After(ctx context.Context, d time.Duration) <-chan time.Time {
	return Clock(ctx).After(d)
}

// AfterFunc is a convenience wrapper for Clock(ctx).AfterFunc.
func AfterFunc(ctx context.Context, d time.Duration, f func()) *clock.Timer {
	return Clock(ctx).AfterFunc(d, f)
}

// NewTicker is a convenience wrapper for Clock(ctx).NewTicker.
func NewTicker(ctx context.Context, d time.Duration) *clock.Ticker {
	return Clock(ctx).NewTicker(d)
}

// NewTimer is a convenience wrapper for Clock(ctx).NewTimer.
func NewTimer(ctx context.Context, d time.Duration) *clock.Timer {
	return Clock(ctx).NewTimer(d)
}

// Now is a convenience wrapper for Clock(ctx).Now.
func Now(ctx context.Context) time.Time {
	return Clock(ctx).Now()
}

// Since is a convenience wrapper for Clock(ctx).Since.
func Since(ctx context.Context, t time.Time) time.Duration {
	return Clock(ctx).Since(t)
}

// Sleep is a convenience wrapper for Clock(ctx).Sleep.
func Sleep(ctx context.Context, d time.Duration) {
	Clock(ctx).Sleep(d)
}

// Tick is a convenience wrapper for Clock(ctx).Tick.
func Tick(ctx context.Context, d time.Duration) <-chan time.Time {
	return Clock(ctx).Tick(d)
}

// Until is a convenience wrapper for Clock(ctx).Until.
func Until(ctx context.Context, t time.Time) time.Duration {
	return Clock(ctx).Until(t)
}
