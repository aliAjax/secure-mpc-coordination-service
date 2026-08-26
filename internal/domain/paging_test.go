package domain

import "testing"

func TestPageStringsDoesNotMutate(t *testing.T) {
	orig := []string{"b", "a", "c"}
	PageStrings(orig, Cursor{Limit: 10})
	if orig[0] != "b" || orig[1] != "a" || orig[2] != "c" {
		t.Fatalf("PageStrings mutated input: %v", orig)
	}
}

func TestNormalizeCursorClampsZero(t *testing.T) {
	c := NormalizeCursor(Cursor{Limit: 0})
	if c.Limit != 100 {
		t.Fatalf("expected limit 100, got %d", c.Limit)
	}
}
