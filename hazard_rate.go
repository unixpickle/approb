package approb

// A HazardDensity approximates the density function of
// a distribution that was given by its hazard rate.
type HazardDensity struct {
	xVals       []float64
	densityVals []float64
}

// NewHazardDensity creates an approximate density function
// from a hazard rate function.
// The density will become 0 after maxTime, thus bounding
// the amount of computation done for the approximation.
func NewHazardDensity(maxTime float64, rate func(t float64) float64) *HazardDensity {
	cumulative := 0.0
	stepWidth := maxTime / numSampleDivisions
	res := &HazardDensity{}
	for i := 0; i < numSampleDivisions && cumulative < 1; i++ {
		t := maxTime * float64(i) / numSampleDivisions
		hr := rate(t)
		density := hr * (1 - cumulative)
		cumulative += density * stepWidth
		res.xVals = append(res.xVals, t)
		res.densityVals = append(res.densityVals, density)
	}

	// Enforce an integral of 1.
	scale := 1 / cumulative
	for i, x := range res.densityVals {
		res.densityVals[i] = x * scale
	}

	return res
}

// Eval evaluates the density at the given time.
func (h *HazardDensity) Eval(t float64) float64 {
	return bilinearEval(h.xVals, h.densityVals, t)
}
