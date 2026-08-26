package domain

import (
	"errors"
	"testing"
)

func TestResolveKnownProtocol(t *testing.T) {
	r := NewProtocolRegistry()
	if _, err := r.Resolve("sum", "1"); err != nil {
		t.Fatalf("expected sum/1 to resolve, got %v", err)
	}
}

func TestResolveUnknownProtocolChain(t *testing.T) {
	r := NewProtocolRegistry()
	if _, err := r.Resolve("nope", "1"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}
