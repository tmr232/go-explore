package itertools

import "container/list"

func FromList[T any](l *list.List) Iterator[T] {
	element := l.Front()
	var value T
	var zero T
	advance := func() (bool, T) {
		if element == nil {
			return false, zero
		}
		value = element.Value.(T)
		element = element.Next()
		return true, value
	}

	return FromAdvance(advance)
}

func ToList[T any](iter Iterator[T]) *list.List {
	l := list.New()
	for iter.Next() {
		l.PushBack(iter.Value())
	}
	return l
}
