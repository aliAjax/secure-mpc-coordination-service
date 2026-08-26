package application

import (
	"context"
	"testing"
)

func TestAuditChainVerify(t *testing.T) {
	a := NewAuditLog()
	a.Append(context.Background(), "u", "create", "c1")
	a.Append(context.Background(), "u", "start", "c1")
	if !a.Verify() {
		t.Fatal("expected valid audit chain")
	}
	// tamper with an event in place
	a.events[0].Digest = "tampered"
	if a.Verify() {
		t.Fatal("expected tampered audit chain to fail verification")
	}
}
