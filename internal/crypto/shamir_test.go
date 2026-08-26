package crypto

import (
	"math/big"
	"testing"
)

func TestSplitReconstruct(t *testing.T) {
	secret := bigInt("42")
	shares, e := Split(secret, 3, 5)
	if e != nil {
		t.Fatal(e)
	}
	got, e := Reconstruct(shares[:3], 3)
	if e != nil {
		t.Fatal(e)
	}
	if got.Cmp(secret) != 0 {
		t.Fatalf("got %s", got)
	}
	if _, e = Reconstruct(shares[:2], 3); e == nil {
		t.Fatal("expected threshold error")
	}
}
func TestCommitment(t *testing.T) {
	c := Commit("value", "nonce")
	if !VerifyCommit("value", "nonce", c) {
		t.Fatal("commitment failed")
	}
	if VerifyCommit("other", "nonce", c) {
		t.Fatal("tamper accepted")
	}
}
func bigInt(s string) *big.Int { v, _ := new(big.Int).SetString(s, 10); return v }
