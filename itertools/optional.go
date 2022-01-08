package itertools

type Optional[T any] struct {
	isPresent bool
	value     T
}

func OptionalWith[T any](value T) Optional[T] {
	return Optional[T]{isPresent: true, value: value}
}

func EmptyOptional[T any]() Optional[T] {
	return Optional[T]{}
}

func (o Optional[T]) Get() (T, bool) {
	return o.value, o.isPresent
}

func (o Optional[T]) Or(alt T) T {
	if o.isPresent {
		return o.value
	}
	return alt
}
