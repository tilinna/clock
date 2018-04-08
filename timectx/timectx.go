// Package timectx implements a library for mocking time in context.Context.
//
// All methods are safe for concurrent use.
package timectx

import (
	"context"
	"time"

	"github.com/tilinna/clock"
)

type clockKey struct{}

// WithClock returns a copy of parent in which the clock is associated with.
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

func WithDeadline(parent context.Context, deadline time.Time) (context.Context, context.CancelFunc) {
	if mock, ok := Clock(parent).(*clock.Mock); ok {
		return mock.DeadlineContext(parent, deadline)
	}
	return context.WithDeadline(parent, deadline)
}

func WithTimeout(parent context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	if mock, ok := Clock(parent).(*clock.Mock); ok {
		return mock.TimeoutContext(parent, timeout)
	}
	return context.WithTimeout(parent, timeout)
}

func After(ctx context.Context, d time.Duration) <-chan time.Time {
	return Clock(ctx).After(d)
}

func AfterFunc(ctx context.Context, d time.Duration, f func()) *clock.Timer {
	return Clock(ctx).AfterFunc(d, f)
}

func NewTicker(ctx context.Context, d time.Duration) *clock.Ticker {
	return Clock(ctx).NewTicker(d)
}

func NewTimer(ctx context.Context, d time.Duration) *clock.Timer {
	return Clock(ctx).NewTimer(d)
}

func Now(ctx context.Context) time.Time {
	return Clock(ctx).Now()
}

func Since(ctx context.Context, t time.Time) time.Duration {
	return Clock(ctx).Since(t)
}

func Sleep(ctx context.Context, d time.Duration) {
	Clock(ctx).Sleep(d)
}

func Tick(ctx context.Context, d time.Duration) <-chan time.Time {
	return Clock(ctx).Tick(d)
}

func Until(ctx context.Context, t time.Time) time.Duration {
	return Clock(ctx).Until(t)
}
