package repository

import (
	"context"
	"testing"
	"time"

	"github.com/example/027-mpc-coordinator/internal/domain"
)

func TestGetComputationIsolation(t *testing.T) {
	s := NewMemoryStore("")
	c := &domain.Computation{ID: "c1", TenantID: "acme", Status: domain.StatusDraft}
	if err := s.CreateComputation(context.Background(), c); err != nil {
		t.Fatal(err)
	}
	got, err := s.GetComputation(context.Background(), "c1")
	if err != nil {
		t.Fatal(err)
	}
	got.Status = domain.StatusSucceeded
	again, err := s.GetComputation(context.Background(), "c1")
	if err != nil {
		t.Fatal(err)
	}
	if again.Status != domain.StatusDraft {
		t.Fatalf("expected isolated copy (draft), got %s", again.Status)
	}
}

func TestGetRoundIsolation(t *testing.T) {
	s := NewMemoryStore("")
	r := &domain.Round{
		ID: "r1", ComputationID: "c1", Status: domain.RoundOpen,
		Deadline: time.Now().Add(time.Hour), Shares: map[string]domain.Share{},
	}
	if err := s.CreateRound(context.Background(), r); err != nil {
		t.Fatal(err)
	}
	got, err := s.GetRound(context.Background(), "r1")
	if err != nil {
		t.Fatal(err)
	}
	got.Shares["p1"] = domain.Share{ParticipantID: "p1"}
	again, err := s.GetRound(context.Background(), "r1")
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := again.Shares["p1"]; ok {
		t.Fatalf("expected isolated round (no p1), got mutated map")
	}
}
