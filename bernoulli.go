package approb

import "math/rand"

// Bernoulli samples a Bernoulli distribution with a
// success probability p and a count n.
func Bernoulli(n int, p float64) float64 {
	var res float64
	for i := 0; i < n; i++ {
		if rand.Float64() < p {
			res++
		}
	}
	return res
}
