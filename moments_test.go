package approb

import (
	"math/rand"
	"testing"

	"github.com/unixpickle/num-analysis/linalg"
)

func TestCovariances(t *testing.T) {
	cov := Covariances(5000000, func() linalg.Vector {
		val1 := rand.NormFloat64()
		val2 := rand.NormFloat64()
		return []float64{val1 + 2*val2, 3*val1 - val2}
	})
	expected := linalg.Vector([]float64{5, 1, 1, 10})
	if expected.Scale(-1).Add(cov.Data).MaxAbs() > 1e-2 {
		t.Errorf("expected %v but got %v", expected, cov)
	}
}
