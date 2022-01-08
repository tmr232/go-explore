package itertools

import "constraints"

type Predicate[T any] func(T) bool

// LessThan creates a predicate that returns true for values less than the input value
func LessThan[T constraints.Ordered](value T) Predicate[T] {
	return func(other T) bool {
		return other < value
	}
}

// GreaterThan creates a predicate that returns true for values greater than the input value
func GreaterThan[T constraints.Ordered](value T) Predicate[T] {
	return func(other T) bool {
		return other > value
	}
}

// EqualTo creates a predicate that returns true for values equal to the input value
func EqualTo[T comparable](value T) Predicate[T] {
	return func(other T) bool {
		return other == value
	}
}

// Or creates a predicate that returns true if any of the base predicates are true.
func Or[T any](predicateA Predicate[T], predicateB Predicate[T]) Predicate[T] {
	return func(value T) bool {
		return predicateA(value) || predicateB(value)
	}
}

// And creates a predicate that returns true if all the base predicates are true.
func And[T any](a, b Predicate[T]) Predicate[T] {
	return func(value T) bool {
		return a(value) && b(value)
	}
}

// All creates a predicate that returns true if all the base predicates are true.
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

// Any creates a predicate that returns true if any of the base predicates are true.
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

// Not applies ! to the result of the input predicate
func Not[T any](predicate Predicate[T]) Predicate[T] {
	return func(value T) bool {
		return !predicate(value)
	}
}

// SliceNotEmpty returns true if the input slice is not empty
func SliceNotEmpty[T any](slice []T) bool {
	return len(slice) > 0
}
