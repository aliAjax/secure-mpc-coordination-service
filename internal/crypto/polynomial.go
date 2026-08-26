package crypto

import "math/big"

type Polynomial struct{ Coefficients []*big.Int }

func NewPolynomial(secret *big.Int, degree int) (Polynomial, error) {
	if err := ValidateSecret(secret); err != nil {
		return Polynomial{}, err
	}
	if degree < 0 || degree > 255 {
		return Polynomial{}, ErrDegree
	}
	coeff := make([]*big.Int, degree+1)
	coeff[0] = Normalize(secret)
	for i := 1; i <= degree; i++ {
		v, e := Random()
		if e != nil {
			return Polynomial{}, e
		}
		coeff[i] = v
	}
	return Polynomial{coeff}, nil
}

var ErrDegree = fieldError("polynomial degree out of range")

type fieldError string

func (e fieldError) Error() string { return string(e) }
func (p Polynomial) Evaluate(x int64) *big.Int {
	result := big.NewInt(0)
	power := big.NewInt(1)
	bx := big.NewInt(x)
	for _, c := range p.Coefficients {
		result = Add(result, Mul(c, power))
		power = Mul(power, bx)
	}
	return result
}
func (p Polynomial) Degree() int {
	if len(p.Coefficients) == 0 {
		return -1
	}
	return len(p.Coefficients) - 1
}
func (p Polynomial) Commitments() []string {
	out := make([]string, len(p.Coefficients))
	for i, c := range p.Coefficients {
		out[i] = Commit(c.String(), string(rune(i)))
	}
	return out
}
func VerifyPolynomial(p Polynomial, commitments []string) bool {
	if len(p.Coefficients) != len(commitments) {
		return false
	}
	for i, c := range p.Coefficients {
		if !VerifyCommit(c.String(), string(rune(i)), commitments[i]) {
			return false
		}
	}
	return true
}
func AddPolynomials(a, b Polynomial) Polynomial {
	n := len(a.Coefficients)
	if len(b.Coefficients) > n {
		n = len(b.Coefficients)
	}
	out := Polynomial{Coefficients: make([]*big.Int, n)}
	for i := 0; i < n; i++ {
		x, y := big.NewInt(0), big.NewInt(0)
		if i < len(a.Coefficients) {
			x = a.Coefficients[i]
		}
		if i < len(b.Coefficients) {
			y = b.Coefficients[i]
		}
		out.Coefficients[i] = Add(x, y)
	}
	return out
}
func ScalePolynomial(p Polynomial, k *big.Int) Polynomial {
	out := Polynomial{Coefficients: make([]*big.Int, len(p.Coefficients))}
	for i, c := range p.Coefficients {
		out.Coefficients[i] = Mul(c, k)
	}
	return out
}
