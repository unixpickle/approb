package approb

import (
	"math"
	"math/rand"

	"github.com/unixpickle/num-analysis/integration"
)

const numSampleDivisions = 1e6

// A Sampler samples a probability distribution.
type Sampler struct {
	xValues   []float64
	integrals []float64
}

// NewSampler creates a Sampler from a probability
// distribution.
//
// The start and end arguments specify the interval.
//
// The distribution function f should be normalized.
func NewSampler(start, end float64, f func(float64) float64) *Sampler {
	var res Sampler
	res.xValues = append(res.xValues, start)
	res.integrals = append(res.integrals, 0)
	spacing := (end - start) / numSampleDivisions
	total := 0.0
	for x := start; x < end; x += spacing {
		total += f(x) * spacing
		res.xValues = append(res.xValues, math.Min(x+spacing, end))
		res.integrals = append(res.integrals, total)
	}
	return &res
}

// Sample samples from the sampler.
func (s *Sampler) Sample() float64 {
	x := rand.Float64()
	min := -1
	max := len(s.integrals)
	for min+1 < max {
		mid := (min + max) / 2
		intVal := s.integrals[mid]
		if intVal > x {
			max = mid
		} else if intVal < x {
			min = mid
		} else if intVal == x {
			return s.xValues[mid]
		}
	}
	if min < 0 {
		return s.xValues[max]
	} else if max >= len(s.xValues) {
		return s.xValues[min]
	}
	int1 := s.integrals[min]
	int2 := s.integrals[max]
	progress := (x - int1) / (int2 - int1)
	return progress*s.xValues[max] + (1-progress)*s.xValues[min]
}

// Normalize creates a probability distribution out of
// a non-normalized distribution function by scaling
// the result.
//
// The start and end parameters specify the bounds
// of the random variable.
func Normalize(start, end float64, f func(float64) float64) func(float64) float64 {
	iv := integration.Interval{Start: start, End: end}
	area := integration.Integrate(f, iv)
	return func(x float64) float64 {
		return f(x) / area
	}
}
