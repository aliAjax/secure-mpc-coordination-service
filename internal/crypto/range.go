package crypto

import "math/big"

type RangeProof struct {
	Min    int64  `json:"min"`
	Max    int64  `json:"max"`
	Digest string `json:"digest"`
}

func ProveRange(value int64, min, max int64) RangeProof {
	d := Commit(big.NewInt(value).String(), big.NewInt(min).String()+":"+big.NewInt(max).String())
	return RangeProof{Min: min, Max: max, Digest: d}
}
func VerifyRange(value int64, p RangeProof) bool {
	if value < p.Min || value > p.Max {
		return false
	}
	return p.Digest == ProveRange(value, p.Min, p.Max).Digest
}
