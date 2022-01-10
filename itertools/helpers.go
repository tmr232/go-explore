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

func fold[T any](init T, slice []T, binOp func(a, b T) T) T {
	value := init
	for _, v := range slice {
		value = binOp(value, v)
	}
	return value
}

func MaxOf[T constraints.Ordered](value T, values ...T) T {
	return fold(value, values, Max[T])
}

func Min[T constraints.Ordered](a, b T) T {
	if a > b {
		return b
	}
	return a
}

func MinOf[T constraints.Ordered](value T, values ...T) T {
	return fold(value, values, Min[T])
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
