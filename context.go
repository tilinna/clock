package clock

import (
	"context"
	"time"
)

type clockKey struct{}

// Context returns a copy of parent in which the Clock c is associated with.
//
// If the clock is an instance of *Mock, all subsequent uses of the
// context in this package are automatically mocked, else the standard
// realtime clock is used.
func Context(parent context.Context, c Clock) context.Context {
	return context.WithValue(parent, clockKey{}, c)
}

// FromContext returns the clock associated with the context, or Realtime().
func FromContext(ctx context.Context) Clock {
	if c, ok := ctx.Value(clockKey{}).(Clock); ok {
		return c
	}
	return Realtime()
}

// WithDeadline returns a copy of the parent context with the deadline adjusted
// to be no later than d.
//
// If the FromContext(parent) returns a *Mock, it is used to mock the deadline,
// else context.WithDeadline is called directly.
func WithDeadline(parent context.Context, d time.Time) (context.Context, context.CancelFunc) {
	if m, ok := FromContext(parent).(*Mock); ok {
		m.Lock()
		defer m.Unlock()
		return m.deadlineContext(parent, d)
	}
	return context.WithDeadline(parent, d)
}

// WithTimeout returns WithDeadline(parent, Clock(parent).Now().Add(timeout)).
//
// If the FromContext(parent) returns a *Mock, it is used to mock the deadline,
// else context.WithTimeout is called directly.
func WithTimeout(parent context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	if m, ok := FromContext(parent).(*Mock); ok {
		m.Lock()
		defer m.Unlock()
		return m.deadlineContext(parent, m.now.Add(timeout))
	}
	return context.WithTimeout(parent, timeout)
}

// After is a convenience wrapper for FromContext(ctx).After.
func After(ctx context.Context, d time.Duration) <-chan time.Time {
	return FromContext(ctx).After(d)
}

// AfterFunc is a convenience wrapper for FromContext(ctx).AfterFunc.
func AfterFunc(ctx context.Context, d time.Duration, f func()) *Timer {
	return FromContext(ctx).AfterFunc(d, f)
}

// NewTicker is a convenience wrapper for FromContext(ctx).NewTicker.
func NewTicker(ctx context.Context, d time.Duration) *Ticker {
	return FromContext(ctx).NewTicker(d)
}

// NewTimer is a convenience wrapper for FromContext(ctx).NewTimer.
func NewTimer(ctx context.Context, d time.Duration) *Timer {
	return FromContext(ctx).NewTimer(d)
}

// Now is a convenience wrapper for FromContext(ctx).Now.
func Now(ctx context.Context) time.Time {
	return FromContext(ctx).Now()
}

// Since is a convenience wrapper for FromContext(ctx).Since.
func Since(ctx context.Context, t time.Time) time.Duration {
	return FromContext(ctx).Since(t)
}

// Sleep is a convenience wrapper for FromContext(ctx).Sleep.
func Sleep(ctx context.Context, d time.Duration) {
	FromContext(ctx).Sleep(d)
}

// Tick is a convenience wrapper for FromContext(ctx).Tick.
func Tick(ctx context.Context, d time.Duration) <-chan time.Time {
	return FromContext(ctx).Tick(d)
}

// Until is a convenience wrapper for FromContext(ctx).Until.
func Until(ctx context.Context, t time.Time) time.Duration {
	return FromContext(ctx).Until(t)
}

func (m *Mock) deadlineContext(parent context.Context, deadline time.Time) (context.Context, context.CancelFunc) {
	cancelCtx, cancel := context.WithCancel(parent)
	if pd, ok := parent.Deadline(); ok && !pd.After(deadline) {
		return cancelCtx, cancel
	}
	ctx := &mockCtx{
		Context:  cancelCtx,
		done:     make(chan struct{}),
		deadline: deadline,
	}
	t := m.newTimerFunc(deadline, nil)
	go func() {
		select {
		case <-t.C:
			ctx.err = context.DeadlineExceeded
		case <-cancelCtx.Done():
			ctx.err = cancelCtx.Err()
			defer t.Stop()
		}
		close(ctx.done)
	}()
	return ctx, cancel
}

type mockCtx struct {
	context.Context
	deadline time.Time
	done     chan struct{}
	err      error
}

func (ctx *mockCtx) Deadline() (time.Time, bool) {
	return ctx.deadline, true
}

func (ctx *mockCtx) Done() <-chan struct{} {
	return ctx.done
}

func (ctx *mockCtx) Err() error {
	select {
	case <-ctx.done:
		return ctx.err
	default:
		return nil
	}
}
