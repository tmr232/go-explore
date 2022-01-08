package itertools

type RingSlice[T any] struct {
	slice    []T
	last     int
	overflow bool
}

func NewRingSlice[T any](size int) RingSlice[T] {
	return RingSlice[T]{slice: make([]T, size)}
}

func (rs *RingSlice[T]) Push(item T) {
	if rs.last >= len(rs.slice) {
		rs.last = 0
		rs.overflow = true
	}

	rs.slice[rs.last] = item

	rs.last++
}

func (rs *RingSlice[T]) Iter() Iterator[T] {
	if !rs.overflow {
		return FromSlice(rs.slice[:rs.last])
	} else {
		return Chain(FromSlice(rs.slice[rs.last:]), FromSlice(rs.slice[:rs.last]))
	}
}

func (rs *RingSlice[T]) ToSlice() []T {
	if !rs.overflow {
		return append([]T{}, rs.slice[:rs.last]...)
	} else {
		return append(append([]T{}, rs.slice[rs.last:]...), rs.slice[:rs.last]...)
	}
}
