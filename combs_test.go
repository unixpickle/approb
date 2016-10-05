package approb

import (
	"reflect"
	"testing"
)

func TestCombs(t *testing.T) {
	expected := [][]int{{0, 1}, {0, 2}, {1, 2}}
	actual := Combs(3, 2)
	for i, exp := range expected {
		act, ok := <-actual
		if !ok {
			t.Fatal("not enough items", i, "expected", len(expected))
		}
		if !reflect.DeepEqual(exp, act) {
			t.Fatal("element", i, "should be", exp, "but got", act)
		}
	}
	if _, ok := <-actual; ok {
		t.Error("too many items")
	}
}
