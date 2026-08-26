package crypto

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func Commit(value, nonce string) string {
	h := sha256.New()
	h.Write([]byte(value))
	h.Write([]byte{0})
	h.Write([]byte(nonce))
	return hex.EncodeToString(h.Sum(nil))
}
func VerifyCommit(value, nonce, digest string) bool {
	expected := Commit(value, nonce)
	return hmac.Equal([]byte(expected), []byte(digest))
}
func SeedCommit(seed string) string {
	s := sha256.Sum256([]byte("seed:" + seed))
	return hex.EncodeToString(s[:])
}
