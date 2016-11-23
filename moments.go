package approb

import (
	"github.com/unixpickle/num-analysis/kahan"
	"github.com/unixpickle/num-analysis/linalg"
)

// Moments calls the function f repeatedly to sample a
// distribution n times.
// It then returns the mean and variance of the samples.
func Moments(n int, f func() float64) (mean, variance float64) {
	var sum, squareSum kahan.Summer64
	for i := 0; i < n; i++ {
		val := f()
		sum.Add(val)
		squareSum.Add(val * val)
	}
	mean = sum.Sum() / float64(n)
	variance = squareSum.Sum()/float64(n) - mean*mean
	return
}

// Indicator is like Moments, but f returns a bool
// which is treated like a 1 if true and 0 if false.
func Indicator(n int, f func() bool) (mean, variance float64) {
	return Moments(n, func() float64 {
		if f() {
			return 1
		} else {
			return 0
		}
	})
}

// Mean computes the mean of a multivariate random
// variable by sampling it n times.
func Mean(n int, f func() linalg.Vector) linalg.Vector {
	var sum linalg.Vector
	for i := 0; i < n; i++ {
		if sum == nil {
			sum = f().Copy()
		} else {
			sum.Add(f())
		}
	}
	sum.Scale(1 / float64(n))
	return sum
}

// Covariances computes the covariance matrix for a
// multivariate distribution (sampled via f).
func Covariances(n int, f func() linalg.Vector) *linalg.Matrix {
	size := len(f())
	res := linalg.NewMatrix(size, size)
	mean := Mean(n, f).Scale(-1)
	for i := 0; i < n; i++ {
		sampled := mean.Copy().Add(f())
		for j := 0; j < size; j++ {
			for k := 0; k <= j; k++ {
				prod := sampled[j] * sampled[k]
				res.Set(j, k, res.Get(j, k)+prod)
			}
		}
	}
	for j := 0; j < size; j++ {
		for k := 0; k <= j; k++ {
			val := res.Get(j, k)
			res.Set(j, k, val/float64(n))
			res.Set(k, j, res.Get(j, k))
		}
	}
	return res
}
