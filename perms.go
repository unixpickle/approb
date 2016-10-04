package approb

// Perms generates a channel of all permutations of the
// numbers in the range [0,n).
// The channel is closed after all permutations have
// been sent.
func Perms(n int) <-chan []int {
	res := make(chan []int)
	remaining := map[int]bool{}
	for i := 0; i < n; i++ {
		remaining[i] = true
	}
	go func() {
		generatePerms(remaining, []int{}, res)
		close(res)
	}()
	return res
}

func generatePerms(remaining map[int]bool, perm []int, res chan<- []int) {
	if len(remaining) == 0 {
		p := make([]int, len(perm))
		copy(p, perm)
		res <- p
		return
	}
	options := make([]int, 0, len(remaining))
	for x := range remaining {
		options = append(options, x)
	}
	for _, x := range options {
		delete(remaining, x)
		newPerm := append(perm, x)
		generatePerms(remaining, newPerm, res)
		remaining[x] = true
	}
}
