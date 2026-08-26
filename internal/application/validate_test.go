package application

import (
	"context"
	"testing"

	"github.com/example/027-mpc-coordinator/internal/crypto"
	"github.com/example/027-mpc-coordinator/internal/domain"
	"github.com/example/027-mpc-coordinator/internal/repository"
)

func TestValidateComputationErrNotFound(t *testing.T) {
	svc := NewService(repository.NewMemoryStore(""), crypto.NewKeyProvider())
	c, err := svc.Create(context.Background(), CreateRequest{
		TenantID: "acme", Protocol: "nope", ProtocolVersion: "1.0",
		Threshold: 2, ParticipantCount: 3,
	}, "create-1")
	if err != nil {
		t.Fatal(err)
	}
	if err := svc.ValidateComputation(context.Background(), c.ID); !domain.IsNotFound(err) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}
