package approb

import (
	"math"
	"testing"

	"github.com/unixpickle/num-analysis/linalg"
)

func TestMVSampler(t *testing.T) {
	distributions := []func(linalg.Vector) float64{
		func(arg3 linalg.Vector) float64 {
			return (9.0 / 4) * math.Pow(arg3[0], 2) * math.Pow(arg3[1], 2)
		},
		func(arg3 linalg.Vector) float64 {
			return (9.0 / 4) * math.Pow(arg3[0]-1, 2) * math.Pow(arg3[1], 2)
		},
	}
	mins := []linalg.Vector{{-1, -1}, {0, -1}}
	maxes := []linalg.Vector{{1, 1}, {2, 1}}
	expected := []linalg.Vector{{0, 0}, {1, 0}}

	for i, f := range distributions {
		s := NewMVSamplerPrec(1e6, mins[i], maxes[i], f)
		var mean linalg.Vector = []float64{0, 0}
		var count float64
		for i := 0; i < 100000; i++ {
			mean.Add(s.Sample())
			count++
		}
		mean.Scale(1 / count)
		if mean.Copy().Scale(-1).Add(expected[i]).MaxAbs() > 1e-2 {
			t.Errorf("task %d: expected %v but got %v", i, expected[i], mean)
		}
	}
}
