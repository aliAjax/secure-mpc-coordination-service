package repository

import (
	"testing"
	"time"
)

func TestTimeIndexBeforeIsolation(t *testing.T) {
	now := time.Now()
	idx := NewTimeIndex[int]()
	idx.Add(now, 1)
	if got := idx.Before(now); len(got) != 0 {
		t.Fatalf("expected 0 items strictly before now, got %d", len(got))
	}
}

func TestTimeIndexAddAscending(t *testing.T) {
	now := time.Now()
	idx := NewTimeIndex[int]()
	idx.Add(now.Add(time.Second), 2)
	idx.Add(now, 1)
	got := idx.Before(now.Add(time.Hour))
	if len(got) != 2 || got[0] != 1 || got[1] != 2 {
		t.Fatalf("expected ascending order [1 2], got %v", got)
	}
}
