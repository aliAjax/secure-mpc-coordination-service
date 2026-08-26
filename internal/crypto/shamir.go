package crypto

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"math/big"
	"sort"
)

type SharePoint struct {
	Index int    `json:"index"`
	Value string `json:"value"`
}

func Split(secret *big.Int, threshold, count int) ([]SharePoint, error) {
	if threshold < 2 || count < threshold || count > 255 {
		return nil, errors.New("invalid threshold/count")
	}
	coeff := make([]*big.Int, threshold)
	coeff[0] = Normalize(secret)
	for i := 1; i < threshold; i++ {
		v, e := Random()
		if e != nil {
			return nil, e
		}
		coeff[i] = v
	}
	out := make([]SharePoint, 0, count)
	for i := 1; i <= count; i++ {
		x := big.NewInt(int64(i))
		y := new(big.Int)
		power := big.NewInt(1)
		for _, c := range coeff {
			y.Add(y, new(big.Int).Mul(c, power))
			y.Mod(y, Prime)
			power.Mul(power, x)
			power.Mod(power, Prime)
		}
		out = append(out, SharePoint{Index: i, Value: y.String()})
	}
	return out, nil
}

func Reconstruct(points []SharePoint, threshold int) (*big.Int, error) {
	if len(points) < threshold || threshold < 2 {
		return nil, errors.New("threshold not met")
	}
	uniq := map[int]bool{}
	selected := points[:threshold]
	for _, p := range selected {
		if p.Index <= 0 || uniq[p.Index] {
			return nil, errors.New("duplicate/invalid share index")
		}
		uniq[p.Index] = true
	}
	for _, p := range selected {
		if _, e := FromDecimal(p.Value); e != nil {
			return nil, e
		}
	}
	result := big.NewInt(0)
	for i, pi := range selected {
		yi, _ := FromDecimal(pi.Value)
		num := big.NewInt(1)
		den := big.NewInt(1)
		xi := big.NewInt(int64(pi.Index))
		for j, pj := range selected {
			if i == j {
				continue
			}
			xj := big.NewInt(int64(pj.Index))
			num = Mul(num, Sub(big.NewInt(0), xj))
			den = Mul(den, Sub(xi, xj))
		}
		inv, e := Inv(den)
		if e != nil {
			return nil, e
		}
		result = Add(result, Mul(yi, Mul(num, inv)))
	}
	return result, nil
}

func SortShares(s []SharePoint) {
	sort.Slice(s, func(i, j int) bool { return s[i].Index < s[j].Index })
}
func ShareDigest(p SharePoint) string {
	h := sha256.Sum256([]byte(p.Value))
	return hex.EncodeToString(h[:])
}
