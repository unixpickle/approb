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

func TestMVSamplerCond(t *testing.T) {
	min := []float64{-20, 0}
	max := []float64{20, 20}
	doubleDist := NewMVSampler(min, max, func(v linalg.Vector) float64 {
		x, y := v[0], v[1]
		return 1 / math.Sqrt(2*math.Pi) * math.Exp(-y) * math.Exp(-math.Pow(x-y, 2)/2)
	})
	fixedX := 0.3
	condDist := NewSampler(0, 20, Normalize(0, 20, func(y float64) float64 {
		return 1 / math.Sqrt(2*math.Pi) * math.Exp(-y) * math.Exp(-math.Pow(fixedX-y, 2)/2)
	}))
	corr := Correlation(10000, 0.1, condDist.Sample, func() float64 {
		out := doubleDist.SampleCond([]float64{fixedX})
		if out[0] != fixedX {
			t.Fatalf("expected out[0]=%f but got %f", fixedX, out[0])
		}
		return out[1]
	})
	if math.Abs(corr-1) > 1e-2 {
		t.Errorf("expected correlation 1 but got %f", corr)
	}
}
