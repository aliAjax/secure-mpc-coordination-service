package domain

import "testing"

func TestTransitions(t *testing.T) {
	c := Computation{Status: StatusDraft}
	if e := TransitionComputation(&c, StatusCommitted); e != nil {
		t.Fatal(e)
	}
	if e := TransitionComputation(&c, StatusSucceeded); e == nil {
		t.Fatal("invalid transition accepted")
	}
	r := Round{Status: RoundOpen}
	if e := TransitionRound(&r, RoundCollecting); e != nil {
		t.Fatal(e)
	}
}
