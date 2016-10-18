package approb

import "github.com/unixpickle/num-analysis/kahan"

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
