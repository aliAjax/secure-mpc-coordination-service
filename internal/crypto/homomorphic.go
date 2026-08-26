package crypto

import "math/big"

type Ciphertext struct {
	Value  string `json:"value"`
	Scheme string `json:"scheme"`
}
type Additive interface {
	Encrypt(*big.Int) (Ciphertext, error)
	Add(Ciphertext, Ciphertext) (Ciphertext, error)
	Decrypt(Ciphertext) (*big.Int, error)
}
type MockAdditive struct{}

func (MockAdditive) Encrypt(v *big.Int) (Ciphertext, error) {
	return Ciphertext{Value: Normalize(v).String(), Scheme: "mock-additive"}, nil
}
func (MockAdditive) Add(a, b Ciphertext) (Ciphertext, error) {
	x, e := FromDecimal(a.Value)
	if e != nil {
		return Ciphertext{}, e
	}
	y, e := FromDecimal(b.Value)
	if e != nil {
		return Ciphertext{}, e
	}
	return Ciphertext{Value: Add(x, y).String(), Scheme: "mock-additive"}, nil
}
func (MockAdditive) Decrypt(c Ciphertext) (*big.Int, error) { return FromDecimal(c.Value) }
