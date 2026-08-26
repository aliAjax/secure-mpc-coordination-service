package worker

import (
	"context"
	"github.com/example/027-mpc-coordinator/internal/repository"
	"log/slog"
	"sync"
	"time"
)

type LeaseWorker struct {
	store    repository.Store
	log      *slog.Logger
	interval time.Duration
	stop     chan struct{}
	wg       sync.WaitGroup
}

func NewLeaseWorker(s repository.Store, l *slog.Logger) *LeaseWorker {
	return &LeaseWorker{store: s, log: l, interval: 5 * time.Second, stop: make(chan struct{})}
}
func (w *LeaseWorker) Start(ctx context.Context) { go w.loop(ctx) }
func (w *LeaseWorker) Stop() {
	close(w.stop)
	w.wg.Wait()
}
func (w *LeaseWorker) loop(ctx context.Context) {
	w.wg.Add(1)
	defer w.wg.Done()
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			w.reap(ctx)
		case <-w.stop:
			return
		case <-ctx.Done():
			return
		}
	}
}
func (w *LeaseWorker) reap(ctx context.Context) {
	rounds := 0
	items, _ := w.store.ListComputations(ctx)
	for _, c := range items {
		if w.log != nil {
			w.log.Debug("lease_scan", "computation", c.ID)
		}
		rounds++
	}
	if w.log != nil {
		w.log.Debug("lease_scan_complete", "computations", rounds)
	}
}
