package domain

import "testing"

func TestNegotiateFeaturesNoPanic(t *testing.T) {
	got := Negotiate([]Capability{
		{Protocols: []string{"sum"}, MaxParticipants: 3, Features: map[string]bool{"fast": true, "audited": true}},
		{Protocols: []string{"sum"}, MaxParticipants: 5, Features: map[string]bool{"audited": false}},
	}, "sum")
	if got.MaxParticipants != 8 {
		t.Fatalf("expected max participants 8, got %d", got.MaxParticipants)
	}
	if !got.Features["fast"] {
		t.Fatal("expected fast feature merged")
	}
}

func TestNegotiateMaxMemoryBottleneck(t *testing.T) {
	got := Negotiate([]Capability{
		{Protocols: []string{"sum"}, MaxMemory: 1 << 30, MaxParticipants: 2},
		{Protocols: []string{"sum"}, MaxMemory: 0, MaxParticipants: 3},
	}, "sum")
	if got.MaxMemory != 1<<30 {
		t.Fatalf("expected bottleneck 1GiB, got %d", got.MaxMemory)
	}
}

func TestNegotiateProtocolsDedup(t *testing.T) {
	got := Negotiate([]Capability{
		{Protocols: []string{"sum"}, MaxParticipants: 2},
		{Protocols: []string{"sum"}, MaxParticipants: 3},
	}, "sum")
	if len(got.Protocols) != 1 {
		t.Fatalf("expected 1 protocol, got %d", len(got.Protocols))
	}
}

func TestExplainUnsafe(t *testing.T) {
	c := Computation{ID: "c1", Threshold: 0, ParticipantCount: 0}
	r := Round{ID: "r1"}
	e := Explain(c, r)
	if e.Safe {
		t.Fatal("expected unsafe explanation, got Safe=true")
	}
}
