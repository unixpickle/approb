package approb

import "math/rand"

// An Urn contains an immutable set of integers.
type Urn []int

// NewUrn creates an urn with integers 0 through
// len(s)-1, where there are s[i] instances of i.
func NewUrn(s []int) Urn {
	var res Urn
	for i, c := range s {
		for j := 0; j < c; j++ {
			res = append(res, i)
		}
	}
	return res
}

// Draw chooses a random element from the urn.
// It returns a new urn with the element removed.
// The urn must not be empty.
func (u Urn) Draw() (int, Urn) {
	if len(u) == 0 {
		panic("urn is empty")
	}
	idx := rand.Intn(len(u))
	val := u[idx]
	newUrn := make(Urn, len(u)-1)
	copy(newUrn, u[:idx])
	copy(newUrn[idx:], u[idx+1:])
	return val, newUrn
}

// Choose randomly selects n elements from the urn
// without replacement.
// The urn must contain at least n elements.
func (u Urn) Choose(n int) (choice, remainder Urn) {
	remainder = u
	for i := 0; i < n; i++ {
		var el int
		el, remainder = remainder.Draw()
		choice = append(choice, el)
	}
	return
}

// AllEqual returns whether all of the elements of the urn
// are equal to the value x.
func (u Urn) AllEqual(x int) bool {
	for _, k := range u {
		if k != x {
			return false
		}
	}
	return true
}
