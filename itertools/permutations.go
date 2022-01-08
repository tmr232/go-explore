package itertools

func ApplyPermutation[T any](in []T, out []T, permutation []int) []T {
	if len(permutation) > len(in) {
		return nil
	}
	if out == nil {
		out = make([]T, len(permutation))
	}
	for write, read := range permutation {
		out[write] = in[read]
	}
	return out
}

func PermutationsOf[T any](slice []T, r int) Iterator[[]T] {
	indexPermutations := IndexPermutations(len(slice), r)
	perm := make([]T, r)
	pool := append([]T{}, slice...)
	return IteratorClosure[[]T]{
		next: indexPermutations.Next,
		value: func() []T {
			return ApplyPermutation(pool, perm, indexPermutations.Value())
		},
	}
}

func SafePermutationsOf[T any](slice []T, r int) Iterator[[]T] {
	indexPermutations := IndexPermutations(len(slice), r)
	pool := append([]T{}, slice...)
	return IteratorClosure[[]T]{
		next: indexPermutations.Next,
		value: func() []T {
			return ApplyPermutation(pool, nil, indexPermutations.Value())
		},
	}
}
