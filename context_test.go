package clock_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/tilinna/clock"
)

var mockTime = time.Date(2018, 1, 1, 10, 0, 0, 0, time.UTC)

func TestContext(t *testing.T) {
	ctx := context.Background()

	c := clock.FromContext(ctx)
	if c != clock.Realtime() {
		t.Fatalf("want realtime clock, got %T", c)
	}

	ctx = clock.Context(ctx, clock.New(mockTime))
	c = clock.FromContext(ctx)
	m, ok := clock.FromContext(ctx).(*clock.Mock)
	if !ok {
		t.Fatalf("want *clock.Mock, got %T", m)
	}

	tm := clock.NewTimer(ctx, 5*time.Second)
	ctx1, cfn1 := clock.WithTimeout(ctx, 10*time.Second)
	defer cfn1()
	ctx2, cfn2 := clock.WithDeadline(ctx, mockTime.Add(15*time.Second))
	defer cfn2()
	ctx3, cfn3 := clock.WithTimeout(ctx, 10*time.Second)
	cfn3()
	<-ctx3.Done()

	if got, want := ctx3.Err(), context.Canceled; want != got {
		t.Fatalf("want ctx3.Err(): %q, got: %q", want, got)
	}

	if d, ok := ctx2.Deadline(); !ok || !d.Equal(mockTime.Add(15*time.Second)) {
		t.Fatalf("want ctx2.Deadline(): %q, got: %q", mockTime.Add(15*time.Second), d)
	}

	var wg sync.WaitGroup
	wg.Add(3)

	var timeout time.Time
	go func() {
		timeout = <-tm.C
		wg.Done()
	}()

	go func() {
		<-ctx1.Done()
		wg.Done()
	}()

	go func() {
		<-ctx2.Done()
		wg.Done()
	}()

	m.Add(20 * time.Second) // fires all timers simultaneously
	wg.Wait()

	if !timeout.Equal(mockTime.Add(5 * time.Second)) {
		t.Fatalf("want tm timer to expire after 5 seconds, got %q", timeout)
	}
	if got, want := ctx1.Err(), context.DeadlineExceeded; want != got {
		t.Fatalf("want ctx1.Err(): %q, got: %q", want, got)
	}
	if got, want := ctx2.Err(), context.DeadlineExceeded; want != got {
		t.Fatalf("want ctx2.Err(): %q, got: %q", want, got)
	}

	<-ctx3.Done()
	if got, want := ctx3.Err(), context.Canceled; want != got {
		t.Fatalf("want ctx3.Err(): %q, got: %q", want, got)
	}
}
