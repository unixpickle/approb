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
	return bilinearEval(s.integrals, s.xValues, rand.Float64())
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

// Correlation evaluates the correlation between two
// sampleable distributions.
// The binSize argument specifies the amount of space
// to use for a bin in the generated histograms.
// The count argument specifies how many samples to make
// before drawing a conclusion.
func Correlation(count int, binSize float64, dist1, dist2 func() float64) float64 {
	histograms := make([]map[int]int, 2)
	for i, f := range []func() float64{dist1, dist2} {
		histograms[i] = map[int]int{}
		for j := 0; j < count; j++ {
			raw := f()
			bin := int(math.Floor(raw / binSize))
			histograms[i][bin]++
		}
	}
	var correlation float64
	var mag1, mag2 float64
	for bin, x := range histograms[0] {
		correlation += float64(x) * float64(histograms[1][bin])
		mag1 += float64(x) * float64(x)
	}
	for _, y := range histograms[1] {
		mag2 += float64(y) * float64(y)
	}
	return correlation / math.Sqrt(mag1*mag2)
}

func bilinearEval(xs, ys []float64, x float64) float64 {
	min := -1
	max := len(xs)
	for min+1 < max {
		mid := (min + max) / 2
		intVal := xs[mid]
		if intVal > x {
			max = mid
		} else if intVal < x {
			min = mid
		} else if intVal == x {
			return ys[mid]
		}
	}
	if min < 0 {
		return ys[max]
	} else if max >= len(xs) {
		return ys[min]
	}
	int1 := xs[min]
	int2 := xs[max]
	progress := (x - int1) / (int2 - int1)
	return progress*ys[max] + (1-progress)*ys[min]
}
