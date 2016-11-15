package approb

import (
	"math"
	"math/rand"

	"github.com/unixpickle/num-analysis/kahan"
	"github.com/unixpickle/num-analysis/linalg"
)

// An mvSamplerNode contains information about a variable
// and its sub-variables in a joint distribution, holding
// a set of parent-variables fixed.
type mvSamplerNode struct {
	TotalProb     float64
	Probabilities []float64
	Values        []float64
	Children      []*mvSamplerNode
}

func newMVSamplerNode(steps int, min, max, givens linalg.Vector,
	f func(linalg.Vector) float64) *mvSamplerNode {
	if len(givens) == len(min) {
		prob := rectProb(steps, min, max, givens, f)
		return &mvSamplerNode{TotalProb: prob}
	}

	varIdx := len(givens)
	step := (max[varIdx] - min[varIdx]) / float64(steps)

	var res mvSamplerNode
	var totalProb kahan.Summer64
	for i := 0; i < steps; i++ {
		coordVal := (float64(i) + 0.5) * step
		givens = append(givens, coordVal)
		subNode := newMVSamplerNode(steps, min, max, givens, f)
		givens = givens[:len(givens)-1]
		res.Probabilities = append(res.Probabilities, subNode.TotalProb)
		res.Values = append(res.Values, coordVal)
		if len(subNode.Values) > 0 {
			res.Children = append(res.Children, subNode)
		}
		totalProb.Add(res.TotalProb)
	}
	res.TotalProb = totalProb.Sum()
	return &res
}

func (m *mvSamplerNode) Sample() linalg.Vector {
	val, idx := bilinearEval(m.Probabilities, m.Values, rand.Float64()*m.TotalProb)
	if len(m.Children) == 0 {
		return []float64{val}
	}
	sub := m.Children[idx].Sample()
	return append([]float64{val}, sub...)
}

func rectProb(steps int, min, max, coord linalg.Vector, f func(linalg.Vector) float64) float64 {
	var rectSize float64
	for i, minVal := range min {
		maxVal := max[i]
		rectSize *= (maxVal - minVal) / float64(steps)
	}
	return f(coord) * rectSize
}

// MVSampler samples from a multivariate distribution.
type MVSampler struct {
	rootNode *mvSamplerNode
}

// NewMVSampler creates an MVSampler from a multivariate
// density function f.
//
// The min and max parameters specify bounds on the
// parameters.
func NewMVSampler(min, max linalg.Vector, f func(linalg.Vector) float64) *MVSampler {
	if len(min) != len(max) {
		panic("mismatching min and max sizes")
	}
	count := 1 + int(math.Ceil(math.Pow(numSampleDivisions, 1/float64(len(min)))))
	return &MVSampler{
		rootNode: newMVSamplerNode(count, min, max, nil, f),
	}
}

// Sample samples the distribution.
func (m *MVSampler) Sample() linalg.Vector {
	return m.rootNode.Sample()
}
