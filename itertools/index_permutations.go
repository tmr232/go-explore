package itertools

func simpleRange(n int) []int {
	r := make([]int, n)
	for i := range r {
		r[i] = i
	}
	return r
}

func reverseRange(start, stop int) []int {
	r := make([]int, start-stop)
	v := start
	for i := range r {
		r[i] = v
		v--
	}
	return r
}

func rotate(slice []int, toStart int) {
	if toStart == 0 {
		return
	}

	read := toStart
	write := 0
	nextRead := 0
	last := len(slice)
	for read != last {
		if write == nextRead {
			nextRead = read
		}
		slice[write], slice[read] = slice[read], slice[write]
		read += 1
		write += 1
	}
	rotate(slice[write:], nextRead-write)
}

type indexPermutationState struct {
	r       int
	n       int
	first   bool
	indices []int
	cycles  []int
	valid   bool
}

func (state *indexPermutationState) Value() []int {
	return state.indices[:state.r]
}

func (state *indexPermutationState) Next() bool {
	r := state.r
	indices := state.indices
	cycles := state.cycles
	n := state.n

	if !state.valid {
		return false
	}

	if state.first {
		state.first = false
		return true
	}

	for i := r - 1; i >= 0; i-- {
		cycles[i] -= 1

		if cycles[i] == 0 {
			rotate(indices[i:], 1)
			cycles[i] = n - i
		} else {
			j := cycles[i]
			indices[i], indices[len(indices)-j] = indices[len(indices)-j], indices[i]
			return true
		}
	}
	state.valid = false
	return false
}

func IndexPermutations(n int, r int) *indexPermutationState {
	return &indexPermutationState{
		r:       r,
		first:   true,
		indices: simpleRange(n),
		cycles:  reverseRange(n, n-r),
		valid:   r <= n,
		n:       n,
	}
}
