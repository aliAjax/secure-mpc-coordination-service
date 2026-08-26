package crypto

import (
	"testing"
)

func TestKeyIDShortNoPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("KeyID panicked: %v", r)
		}
	}()
	_ = KeyID([]byte{1, 2, 3})
}

func TestDecodeShareEmptyNoPanic(t *testing.T) {
	if _, err := DecodeShare(""); err == nil {
		t.Fatal("expected error for empty share")
	}
}

func TestDecodeShareEmptyValueNoPanic(t *testing.T) {
	enc, err := EncodeShare(SharePoint{Index: 1, Value: ""})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := DecodeShare(enc); err == nil {
		t.Fatal("expected error for share with empty value")
	}
}

func TestFixedPointDecodeNil(t *testing.T) {
	f := FixedPoint{Scale: 100}
	if _, err := f.Decode(nil); err == nil {
		t.Fatal("expected error for nil fixed point value")
	}
}
