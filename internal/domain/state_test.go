package domain

import "testing"

func TestTransitionIntermediateToSucceeded(t *testing.T) {
	c := &Computation{Status: StatusWaitingShares}
	if err := TransitionComputation(c, StatusReconstructing); err != nil {
		t.Fatalf("expected WaitingShares -> Reconstructing, got %v", err)
	}
	if err := TransitionComputation(c, StatusSucceeded); err != nil {
		t.Fatalf("expected Reconstructing -> Succeeded, got %v", err)
	}
}

func TestRoundCollectingToComplete(t *testing.T) {
	r := &Round{Status: RoundCollecting}
	if err := TransitionRound(r, RoundComplete); err != nil {
		t.Fatalf("expected RoundCollecting -> RoundComplete, got %v", err)
	}
}
