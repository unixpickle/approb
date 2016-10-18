package approb

import (
	"math"
	"testing"
)

func TestPoisson(t *testing.T) {
	var sum, count float64
	for i := 0; i < 500000; i++ {
		count++
		sum += Poisson(3.7)
	}
	diff := math.Abs(sum/count - 3.7)
	if diff > 1e-2 {
		t.Errorf("expectation should be %f but got %f", 3.7, sum/count)
	}
}
