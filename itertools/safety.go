package itertools

func CopySlices[T any](iter Iterator[[]T]) Iterator[[]T] {
	value := func() []T {
		return append([]T{}, iter.Value()...)
	}
	return IteratorClosure[[]T]{next: iter.Next, value: value}
}
