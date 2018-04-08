package clock

import (
	"context"
	"time"
)

func (m *Mock) DeadlineContext(parent context.Context, d time.Time) (context.Context, context.CancelFunc) {
	m.Lock()
	defer m.Unlock()
	return m.deadlineContext(parent, d)
}

func (m *Mock) TimeoutContext(parent context.Context, d time.Duration) (context.Context, context.CancelFunc) {
	m.Lock()
	defer m.Unlock()
	return m.deadlineContext(parent, m.now.Add(d))
}

func (m *Mock) deadlineContext(parent context.Context, deadline time.Time) (context.Context, context.CancelFunc) {
	m.Lock()
	defer m.Unlock()
	cancelCtx, cancel := context.WithCancel(parent)
	if pd, ok := parent.Deadline(); ok && !deadline.After(pd) {
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
