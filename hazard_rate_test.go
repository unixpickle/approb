package approb

import (
	"math"
	"math/rand"
	"testing"
)

func TestHazardDensity(t *testing.T) {
	rate := NewHazardDensity(20, func(x float64) float64 {
		return 2
	})
	for i := 0; i < 1000; i++ {
		arg := rand.Float64() * 20
		actual := rate.Eval(arg)
		expected := 2 * math.Exp(-2*arg)
		if math.Abs(actual-expected) > 1e-4 {
			t.Fatalf("argument %f should give %f but gave %f",
				arg, expected, actual)
		}
	}
}
