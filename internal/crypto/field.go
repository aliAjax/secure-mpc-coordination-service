package crypto

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
)

var Prime, _ = new(big.Int).SetString("170141183460469231731687303715884105727", 10)

var (
	ErrFieldInvalid = errors.New("invalid field integer")
	ErrFieldRange   = errors.New("field integer out of range")
)

func Normalize(v *big.Int) *big.Int {
	x := new(big.Int).Mod(new(big.Int).Set(v), Prime)
	if x.Sign() < 0 {
		x.Add(x, Prime)
	}
	return x
}
func Add(a, b *big.Int) *big.Int { return Normalize(new(big.Int).Add(a, b)) }
func Mul(a, b *big.Int) *big.Int { return Normalize(new(big.Int).Mul(a, b)) }
func Sub(a, b *big.Int) *big.Int { return Normalize(new(big.Int).Sub(a, b)) }
func Inv(a *big.Int) (*big.Int, error) {
	x := Normalize(a)
	if x.Sign() == 0 {
		return nil, errors.New("zero has no inverse")
	}
	return new(big.Int).ModInverse(x, Prime), nil
}
func Random() (*big.Int, error) { return rand.Int(rand.Reader, Prime) }
func FromDecimal(s string) (*big.Int, error) {
	x, ok := new(big.Int).SetString(s, 10)
	if !ok {
		return nil, fmt.Errorf("%w: %q", ErrFieldInvalid, s)
	}
	if x.Sign() < 0 || x.Cmp(Prime) >= 0 {
		return nil, fmt.Errorf("%w: %q", ErrFieldRange, s)
	}
	return x, nil
}
