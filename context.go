package clock

import (
	"context"
	"time"
)

type clockKey struct{}

// FromContext returns the Clock associated with the context, or Realtime().
func FromContext(ctx context.Context) Clock {
	if c, ok := ctx.Value(clockKey{}).(Clock); ok {
		return c
	}
	return Realtime()
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

// DeadlineContext is a convenience wrapper for FromContext(ctx).DeadlineContext.
func DeadlineContext(ctx context.Context, d time.Time) (context.Context, context.CancelFunc) {
	return FromContext(ctx).DeadlineContext(ctx, d)
}

// TimeoutContext is a convenience wrapper for FromContext(ctx).TimeoutContext.
func TimeoutContext(ctx context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	return FromContext(ctx).TimeoutContext(ctx, timeout)
}

// Context implements Clock.
func (m *Mock) Context(parent context.Context) context.Context {
	return context.WithValue(parent, clockKey{}, m)
}

// DeadlineContext implements Clock.
func (m *Mock) DeadlineContext(parent context.Context, d time.Time) (context.Context, context.CancelFunc) {
	m.Lock()
	defer m.Unlock()
	return m.deadlineContext(parent, d)
}

// TimeoutContext implements Clock.
func (m *Mock) TimeoutContext(parent context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	m.Lock()
	defer m.Unlock()
	return m.deadlineContext(parent, m.now.Add(timeout))
}

func (m *Mock) deadlineContext(parent context.Context, deadline time.Time) (context.Context, context.CancelFunc) {
	cancelCtx, cancel := context.WithCancel(m.Context(parent))
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
