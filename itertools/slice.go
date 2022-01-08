package itertools

import "github.com/edwingeng/deque"

func dequeToSlice[T any](dq deque.Deque) []T {
	slice := make([]T, dq.Len())
	for i, elem := range dq.DequeueMany(0) {
		slice[i] = elem.(T)
	}
	return slice
}

// ToSlice consumes an iterator, returning a slice of all of its values.
// ToSlice should not be called on infinite iterators!
func ToSlice[T any](iter Iterator[T]) []T {
	dq := deque.NewDeque()
	ForEach(iter, func(t T) { dq.PushBack(t) })
	return dequeToSlice[T](dq)
}

type sliceIterator[T any] struct {
	slice []T
	index int
}

func (s *sliceIterator[T]) Next() bool {
	s.index++
	return s.index < len(s.slice)
}

func (s *sliceIterator[T]) Value() T {
	return s.slice[s.index]
}

func (s *sliceIterator[T]) Len() int {
	return len(s.slice) - s.index - 1
}

// FromSlice wraps a slice in an iterator.
// The input slice should not be mutated!
func FromSlice[T any](slice []T) Iterator[T] {
	return &sliceIterator[T]{slice: slice, index: -1}
}
