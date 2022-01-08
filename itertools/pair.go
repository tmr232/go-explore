package itertools

type Pair[T any, U any] struct {
	first  T
	second U
}

func MakePair[T any, U any](first T, second U) Pair[T, U] {
	return Pair[T, U]{first, second}
}

func (p Pair[T, U]) First() T {
	return p.first
}

func (p Pair[T, U]) Second() U {
	return p.second
}

func PairToSlice[T any](pair Pair[T, T]) []T {
	return []T{pair.first, pair.second}
}
