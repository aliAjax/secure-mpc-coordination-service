package crypto

import (
	"errors"
	"math/big"
)

func ValidateSecret(v *big.Int) error {
	if v == nil || v.Sign() < 0 || v.Cmp(Prime) >= 0 {
		return errors.New("secret outside finite field")
	}
	return nil
}
func ValidatePoints(points []SharePoint) error {
	seen := map[int]bool{}
	for _, p := range points {
		if p.Index <= 0 || p.Index > 255 || seen[p.Index] {
			return errors.New("invalid share points")
		}
		seen[p.Index] = true
		if _, e := FromDecimal(p.Value); e != nil {
			return e
		}
	}
	return nil
}
func InterpolateAt(points []SharePoint, x int) *big.Int {
	if len(points) == 0 {
		return big.NewInt(0)
	}
	result := big.NewInt(0)
	for i, pi := range points {
		yi, _ := FromDecimal(pi.Value)
		num, den := big.NewInt(1), big.NewInt(1)
		for j, pj := range points {
			if i == j {
				continue
			}
			num = Mul(num, Sub(big.NewInt(int64(x)), big.NewInt(int64(pj.Index))))
			den = Mul(den, Sub(big.NewInt(int64(pi.Index)), big.NewInt(int64(pj.Index))))
		}
		inv, e := Inv(den)
		if e != nil {
			continue
		}
		result = Add(result, Mul(yi, Mul(num, inv)))
	}
	return result
}
