package approb

import (
	"math"
	"math/rand"
)

// Poisson samples a Poisson distribution with intensity
// equal to p.
func Poisson(p float64) float64 {
	scaler := 1.0
	goal := rand.Float64() * math.Exp(p)

	// Don't let us run more than 1000 times past
	// the standard deviation.
	maxRet := 1000 * p

	for i := 0.0; i < maxRet; i++ {
		goal -= scaler
		scaler *= p / (i + 1)
		if goal <= 0 {
			return i
		}
	}

	return maxRet
}
