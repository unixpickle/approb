package approb

// Combs generates the ways to choose m numbers out
// of the range [0, n).
// Each combination is a sorted (ascending) list of
// distinct numbers from [0, n).
func Combs(n, m int) <-chan []int {
	res := make(chan []int)
	go func() {
		combs([]int{}, 0, n, m, res)
		close(res)
	}()
	return res
}

func combs(substr []int, start, end, rem int, res chan<- []int) {
	if rem == 0 {
		c := make([]int, len(substr))
		copy(c, substr)
		res <- c
		return
	}
	for i := start; i <= end-rem; i++ {
		prefix := append(substr, i)
		combs(prefix, i+1, end, rem-1, res)
	}
}
