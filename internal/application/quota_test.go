package application

import (
	"context"
	"testing"
	"time"

	"github.com/example/027-mpc-coordinator/internal/domain"
)

func TestReserveHonorsCancel(t *testing.T) {
	q := NewQuotaManager()
	q.SetLimit("acme", domain.ResourceBudget{CPUSeconds: 100})
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	q.Reserve(ctx, "acme", domain.ResourceBudget{CPUSeconds: 1})
	time.Sleep(50 * time.Millisecond)
	if used, _ := q.Snapshot("acme"); used.CPUSeconds != 0 {
		t.Fatalf("expected no quota used, got %d", used.CPUSeconds)
	}
}

func TestReleaseHonorsCancel(t *testing.T) {
	q := NewQuotaManager()
	q.SetLimit("acme", domain.ResourceBudget{CPUSeconds: 100})
	q.Reserve(context.Background(), "acme", domain.ResourceBudget{CPUSeconds: 10})
	canceled, cancel := context.WithCancel(context.Background())
	cancel()
	q.Release(canceled, "acme", domain.ResourceBudget{CPUSeconds: 10})
	if used, _ := q.Snapshot("acme"); used.CPUSeconds != 10 {
		t.Fatalf("expected release skipped, got %d", used.CPUSeconds)
	}
}
