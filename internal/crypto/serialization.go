package crypto

import (
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"math/big"
)

func EncodeShare(s SharePoint) (string, error) {
	b, e := json.Marshal(s)
	if e != nil {
		return "", e
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
func DecodeShare(v string) (SharePoint, error) {
	b, e := base64.RawURLEncoding.DecodeString(v)
	if e != nil {
		return SharePoint{}, e
	}
	var s SharePoint
	if len(b) == 0 {
		return SharePoint{}, nil
	}
	e = json.Unmarshal(b, &s)
	if e != nil {
		return SharePoint{}, e
	}
	return s, nil
}
func ConstantTimeEqual(a, b string) bool {
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}
func EqualField(a, b *big.Int) bool {
	return ConstantTimeEqual(Normalize(a).String(), Normalize(b).String())
}
