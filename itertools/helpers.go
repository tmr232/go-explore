package itertools

import "constraints"

func Identity[T any](t T) T {
	return t
}

func IsSameValue[T comparable](a, b T) bool {
	return a == b
}

func PairAdaptor[T any, R any](f func(T, T) R) func(Pair[T, T]) R {
	return func(p Pair[T, T]) R {
		return f(p.first, p.second)
	}
}

func Max[T constraints.Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}

func Min[T constraints.Ordered](a, b T) T {
	if a > b {
		return b
	}
	return a
}

func ItemGetter[T any](i int) func([]T) T {
	return func(slice []T) T {
		return slice[i]
	}
}

func makeSlice[T any](n int, gen func() T) []T {
	slice := make([]T, n)
	for i := range slice {
		slice[i] = gen()
	}
	return slice
}
