package approb

import (
	"math"
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
