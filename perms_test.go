package approb

import "testing"

func TestPerms(t *testing.T) {
	p := Perms(3)
	actual := map[int]bool{}
	for x := range p {
		if len(x) != 3 {
			t.Fatal("invalid length", len(x))
		}
		nums := map[int]bool{}
		nums[x[0]] = true
		nums[x[1]] = true
		nums[x[2]] = true
		if !nums[0] || !nums[1] || !nums[2] {
			t.Fatal("bad perm", x)
		}
		actual[x[0]*9+x[1]*3+x[2]] = true
	}
	if len(actual) != 6 {
		t.Error("unexpected perm count", len(actual))
	}
}
