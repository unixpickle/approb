package approb

import (
	"math"
	"math/rand"
	"testing"
)

func TestSampleDist(t *testing.T) {
	sampler := NewSampler(-20, 20, func(x float64) float64 {
		return 1 / math.Sqrt(2*math.Pi) * math.Exp(-x*x/2)
	})
	avg, v := Moments(10000, sampler.Sample)
	if math.Abs(avg) > 0.01 {
		t.Errorf("expected avg 0 got %f", avg)
	}
	if math.Abs(v-1) > 0.01 {
		t.Errorf("expected var 1 got %f", v)
	}
}

func TestCorrelation(t *testing.T) {
	corr := Correlation(100000, 0.1, func() float64 {
		return rand.Float64()
	}, func() float64 {
		return rand.Float64()
	})
	if math.Abs(corr-1) > 1e-2 {
		t.Errorf("expected correlation 1 but got %f", corr)
	}
	corr = Correlation(500000, 0.1, func() float64 {
		return rand.NormFloat64()
	}, func() float64 {
		return rand.Float64()
	})
	if math.Abs(corr - 0.287) > 1e-2 {
		t.Errorf("expected correlation 0.287 but got %f", corr)
	}
}
