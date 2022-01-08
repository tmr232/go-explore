package itertools

import "constraints"

func LessThan[T constraints.Ordered](value T) Predicate[T] {
	return func(other T) bool {
		return other < value
	}
}

func GreaterThan[T constraints.Ordered](value T) Predicate[T] {
	return func(other T) bool {
		return other > value
	}
}

func EqualTo[T comparable](value T) Predicate[T] {
	return func(other T) bool {
		return other == value
	}
}

type Predicate[T any] func(T) bool

func Or[T any](predicateA Predicate[T], predicateB Predicate[T]) Predicate[T] {
	return func(value T) bool {
		return predicateA(value) || predicateB(value)
	}
}

func And[T any](a, b Predicate[T]) Predicate[T] {
	return func(value T) bool {
		return a(value) && b(value)
	}
}

func All[T any](preds ...Predicate[T]) Predicate[T] {
	return func(value T) bool {
		for _, pred := range preds {
			if !pred(value) {
				return false
			}
		}
		return true
	}
}

func Any[T any](preds ...Predicate[T]) Predicate[T] {
	return func(value T) bool {
		for _, pred := range preds {
			if pred(value) {
				return true
			}
		}
		return false
	}
}

func Not[T any](predicate Predicate[T]) Predicate[T] {
	return func(value T) bool {
		return !predicate(value)
	}
}

func SliceNotEmpty[T any](slice []T) bool {
	return len(slice) > 0
}
