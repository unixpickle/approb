package approb

import (
	"math"
	"math/rand"
	"sort"

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
	step := (max[varIdx] - min[varIdx]) / float64(steps-1)

	var res mvSamplerNode
	var totalProb kahan.Summer64
	for i := 0; i < steps; i++ {
		coordVal := min[varIdx] + float64(i)*step
		givens = append(givens, coordVal)
		subNode := newMVSamplerNode(steps, min, max, givens, f)
		givens = givens[:len(givens)-1]
		res.Probabilities = append(res.Probabilities, totalProb.Sum())
		res.Values = append(res.Values, coordVal)
		if len(subNode.Values) > 0 {
			res.Children = append(res.Children, subNode)
		}

		// We have (steps-1) rectangles but (steps) children.
		// Thus, we don't want to count the last probability
		// as part of a rectangle.
		if i+1 < steps {
			totalProb.Add(subNode.TotalProb)
		}
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

func (m *mvSamplerNode) SampleCond(v linalg.Vector) linalg.Vector {
	if len(v) == 0 {
		return m.Sample()
	}
	idx := sort.SearchFloat64s(m.Values, v[0])
	if idx == len(m.Values) {
		idx = len(m.Values) - 1
	}
	if idx > 0 && math.Abs(m.Values[idx-1]-v[0]) < math.Abs(m.Values[idx]-v[0]) {
		idx--
	}
	return append([]float64{v[0]}, m.Children[idx].SampleCond(v[1:])...)
}

func rectProb(steps int, min, max, coord linalg.Vector, f func(linalg.Vector) float64) float64 {
	rectSize := 1.0
	for i, minVal := range min {
		maxVal := max[i]
		rectSize *= (maxVal - minVal) / float64(steps-1)
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
	return NewMVSamplerPrec(numSampleDivisions, min, max, f)
}

// NewMVSamplerPrec creates an MVSampler.
//
// Building the distribution will result in f being called
// approximately n times.
func NewMVSamplerPrec(n int, min, max linalg.Vector, f func(linalg.Vector) float64) *MVSampler {
	if len(min) != len(max) {
		panic("mismatching min and max sizes")
	}
	count := 1 + int(math.Ceil(math.Pow(float64(n), 1/float64(len(min)))))
	return &MVSampler{
		rootNode: newMVSamplerNode(count, min, max, nil, f),
	}
}

// Sample samples the distribution.
func (m *MVSampler) Sample() linalg.Vector {
	return m.rootNode.Sample()
}

// SampleCond samples the distribution conditioned on the
// values for a subset of the variables.
// More specifically, the first len(v) variables will be
// set to the values in v.
//
// The entire set of sampled variables will be returned,
// including the pre-determined values in v.
func (m *MVSampler) SampleCond(v linalg.Vector) linalg.Vector {
	return m.rootNode.SampleCond(v)
}
