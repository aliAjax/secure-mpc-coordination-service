package application

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/example/027-mpc-coordinator/internal/crypto"
	"github.com/example/027-mpc-coordinator/internal/repository"
)

func TestConcurrentIdempotentCreate(t *testing.T) {
	svc := NewService(repository.NewMemoryStore(""), crypto.NewKeyProvider())
	const n = 20
	start := make(chan struct{})
	var wg sync.WaitGroup
	ids := make([]string, n)
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			<-start
			c, err := svc.Create(context.Background(), CreateRequest{
				TenantID: "acme", Protocol: "sum", ProtocolVersion: "1",
				Threshold: 2, ParticipantCount: 3,
			}, "same-key")
			if err != nil {
				t.Errorf("create error: %v", err)
				return
			}
			ids[i] = c.ID
		}(i)
	}
	close(start)
	wg.Wait()

	uniq := map[string]bool{}
	for _, id := range ids {
		uniq[id] = true
	}
	if len(uniq) != 1 {
		t.Fatalf("expected exactly 1 computation for same idempotency key, got %d", len(uniq))
	}
}

func TestExpiredLeaseRejected(t *testing.T) {
	svc := NewService(repository.NewMemoryStore(""), crypto.NewKeyProvider())
	c, err := svc.Create(context.Background(), CreateRequest{
		TenantID: "acme", Protocol: "sum", ProtocolVersion: "1",
		Threshold: 2, ParticipantCount: 3,
	}, "create-1")
	if err != nil {
		t.Fatal(err)
	}
	r, err := svc.Start(context.Background(), c.ID)
	if err != nil {
		t.Fatal(err)
	}
	r.LeaseOwner = "old-owner"
	r.LeaseUntil = time.Now().UTC().Add(-time.Minute)
	if err := svc.store.UpdateRound(context.Background(), r); err != nil {
		t.Fatal(err)
	}

	got, err := svc.AcquireLease(context.Background(), r.ID, "new-owner", 30*time.Second)
	if err != nil {
		t.Fatalf("expected expired lease to be acquirable, got %v", err)
	}
	if got.LeaseOwner != "new-owner" {
		t.Fatalf("expected lease owner new-owner, got %s", got.LeaseOwner)
	}
}
