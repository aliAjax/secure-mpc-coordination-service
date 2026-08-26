package crypto

import (
	"errors"
	"math/big"
	"testing"
)

func TestReconstructThresholdErrorChain(t *testing.T) {
	if _, err := Reconstruct(nil, 2); !errors.Is(err, ErrThresholdNotMet) {
		t.Fatalf("expected ErrThresholdNotMet, got %v", err)
	}
}

func TestReconstructDuplicateErrorChain(t *testing.T) {
	pts := []SharePoint{{Index: 1, Value: "1"}, {Index: 1, Value: "2"}}
	if _, err := Reconstruct(pts, 2); !errors.Is(err, ErrDuplicateShare) {
		t.Fatalf("expected ErrDuplicateShare, got %v", err)
	}
}

func TestFromDecimalErrorChain(t *testing.T) {
	if _, err := FromDecimal("abc"); !errors.Is(err, ErrFieldInvalid) {
		t.Fatalf("expected ErrFieldInvalid, got %v", err)
	}
}

func TestFromDecimalRangeErrorChain(t *testing.T) {
	if _, err := FromDecimal(big.NewInt(-1).String()); !errors.Is(err, ErrFieldRange) {
		t.Fatalf("expected ErrFieldRange, got %v", err)
	}
}

func TestSplitInvalidArgsErrorChain(t *testing.T) {
	if _, err := Split(big.NewInt(42), 1, 3); !errors.Is(err, ErrInvalidSplitArg) {
		t.Fatalf("expected ErrInvalidSplitArg, got %v", err)
	}
}
