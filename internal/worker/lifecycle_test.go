package worker

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/example/027-mpc-coordinator/internal/repository"
)

func TestLeaseWorkerDoubleStopNoPanic(t *testing.T) {
	s := repository.NewMemoryStore("")
	w := NewLeaseWorker(s, nil)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	w.Start(ctx)
	time.Sleep(30 * time.Millisecond)
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("double Stop panicked: %v", r)
		}
	}()
	w.Stop()
	w.Stop()
}

func TestRetryHonorsContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	calls := 0
	err := Retry(ctx, RetryPolicy{Attempts: 3, Base: time.Millisecond, Max: time.Millisecond}, func(context.Context) error {
		calls++
		return errors.New("boom")
	})
	if calls != 0 {
		t.Fatalf("expected 0 calls for canceled context, got %d", calls)
	}
	if err == nil {
		t.Fatal("expected non-nil error")
	}
}

func TestRetryDelayCapped(t *testing.T) {
	p := RetryPolicy{Base: time.Second, Max: 2 * time.Second}
	d := p.Delay(2)
	if d > 3*time.Second {
		t.Fatalf("expected delay capped near 2s, got %v", d)
	}
}


func TestRetryConcurrentCancel(t *testing.T) {
	const n = 8
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	start := make(chan struct{})
	var wg sync.WaitGroup
	var calls int32
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			_ = Retry(ctx, RetryPolicy{Attempts: 2, Base: time.Millisecond, Max: time.Millisecond}, func(context.Context) error {
				atomic.AddInt32(&calls, 1)
				return errors.New("boom")
			})
		}()
	}
	close(start)
	wg.Wait()
	if atomic.LoadInt32(&calls) != 0 {
		t.Fatalf("expected 0 calls for canceled context, got %d", calls)
	}
}
