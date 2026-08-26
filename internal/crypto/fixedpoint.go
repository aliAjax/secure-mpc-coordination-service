package crypto

import (
	"errors"
	"math/big"
)

type FixedPoint struct{ Scale int64 }

func (f FixedPoint) Encode(v int64) *big.Int {
	return Normalize(new(big.Int).Mul(big.NewInt(v), big.NewInt(f.Scale)))
}
func (f FixedPoint) Decode(v *big.Int) (int64, error) {
	if f.Scale <= 0 {
		return 0, errors.New("invalid scale")
	}
	q := new(big.Int).Quo(v, big.NewInt(f.Scale))
	if !q.IsInt64() {
		return 0, errors.New("decoded overflow")
	}
	return q.Int64(), nil
}
func AddEncoded(values []*big.Int) *big.Int {
	sum := big.NewInt(0)
	for _, v := range values {
		sum = Add(sum, v)
	}
	return sum
}
