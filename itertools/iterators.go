package itertools

import (
	"constraints"
	"github.com/edwingeng/deque"
)

type Iterator[T any] interface {
	Next() bool // Advance to next value, return true of one exists.
	Value() T   // Get the current value
}

type SizedIterator[T any] interface {
	Iterator[T]
	Len() int // Number or remaining valid .Next() calls
}

type IteratorClosure[T any] struct {
	next  func() bool
	value func() T
}

func (ic IteratorClosure[T]) Next() bool {
	return ic.next()
}

func (ic IteratorClosure[T]) Value() T {
	return ic.value()
}

func ToSlice[T any](iter Iterator[T]) []T {
	slice := []T{}
	for iter.Next() {
		slice = append(slice, iter.Value())
	}
	return slice
}

type SliceIterator[T any] struct {
	slice []T
	index int
}

func (s *SliceIterator[T]) Next() bool {
	s.index++
	return s.index < len(s.slice)
}

func (s *SliceIterator[T]) Value() T {
	return s.slice[s.index]
}

func (s *SliceIterator[T]) Len() int {
	return len(s.slice) - s.index - 1
}

func FromSlice[T any](slice []T) Iterator[T] {
	return &SliceIterator[T]{slice: slice, index: -1}
}

type FilterInIterator[T any] struct {
	iter     Iterator[T]
	filterIn func(T) bool
}

func (f *FilterInIterator[T]) Next() bool {
	for f.iter.Next() {
		if f.filterIn(f.iter.Value()) {
			return true
		}
	}
	return false
}

func (f *FilterInIterator[T]) Value() T {
	return f.iter.Value()
}

func FilterIn[T any](iter Iterator[T], filterIn func(T) bool) Iterator[T] {
	return &FilterInIterator[T]{iter: iter, filterIn: filterIn}
}

type IntRangeIterator struct {
	stop    int
	current int
}

func (i *IntRangeIterator) Next() bool {
	i.current++
	return i.current < i.stop

}

func (i *IntRangeIterator) Value() int {
	return i.current
}

func IntRange(stop int) Iterator[int] {
	return &IntRangeIterator{stop: stop, current: -1}
}

type CycleIterator[T any] struct {
	iter     Iterator[T]
	slice    []T
	dq       deque.Deque
	size     int
	index    int
	consumed bool
	value    T
}

func (ci *CycleIterator[T]) Next() bool {
	if !ci.consumed {
		if ci.iter.Next() {
			ci.value = ci.iter.Value()
			ci.dq.PushBack(ci.value)
			ci.size++
			return true
		} else {
			ci.consumed = true
			ci.slice = make([]T, ci.size)
			for i, elem := range ci.dq.DequeueMany(0) {
				ci.slice[i] = elem.(T)
			}
		}
	}

	if ci.size == 0 {
		return false
	}

	ci.index++

	if ci.index >= len(ci.slice) {
		ci.index = 0
	}
	ci.value = ci.slice[ci.index]
	return true
}

func (ci *CycleIterator[T]) Value() T {
	return ci.value
}

func Cycle[T any](iter Iterator[T]) Iterator[T] {
	return &CycleIterator[T]{iter: iter, index: -1, dq: deque.NewDeque()}
}

type ISliceIterator[T any] struct {
	iter Iterator[T]
	skip int
	take int
}

func (is *ISliceIterator[T]) Next() bool {
	for ; is.skip > 0; is.skip-- {
		if !is.iter.Next() {
			return false
		}
	}

	if is.take <= 0 {
		return false
	}
	is.take--

	return is.iter.Next()
}

func (is *ISliceIterator[T]) Value() T {
	return is.iter.Value()
}

func Take[T any](n int, iter Iterator[T]) Iterator[T] {
	return &ISliceIterator[T]{iter: iter, take: n}
}

func Drop[T any](n int, iter Iterator[T]) Iterator[T] {
	return &ISliceIterator[T]{iter: iter, skip: n}
}

func ISlice[T any](iter Iterator[T], start, end int) Iterator[T] {
	return &ISliceIterator[T]{iter: iter, skip: start - 1, take: end - start}
}

type Pair[T any, U any] struct {
	first  T
	second U
}

func PairToSlice[T any](pair Pair[T, T]) []T {
	return []T{pair.first, pair.second}
}

type ZipIterator[T any, U any] struct {
	first  Iterator[T]
	second Iterator[U]
}

func (zip *ZipIterator[T, U]) Next() bool {
	hasFirst := zip.first.Next()
	hasSecond := zip.second.Next()

	return hasFirst && hasSecond

}

func (zip *ZipIterator[T, U]) Value() Pair[T, U] {
	return Pair[T, U]{zip.first.Value(), zip.second.Value()}
}

type ZipLongestIterator[T any, U any] struct {
	ZipIterator[T, U]
	fill      Pair[T, U]
	hasFirst  bool
	hasSecond bool
}

func (zip *ZipLongestIterator[T, U]) Next() bool {
	zip.hasFirst = zip.first.Next()
	zip.hasSecond = zip.second.Next()

	return zip.hasFirst || zip.hasSecond
}

func (zip *ZipLongestIterator[T, U]) Value() Pair[T, U] {
	if zip.hasFirst && zip.hasSecond {
		return Pair[T, U]{zip.first.Value(), zip.second.Value()}
	}
	if zip.hasFirst {
		return Pair[T, U]{zip.first.Value(), zip.fill.second}
	}
	return Pair[T, U]{zip.fill.first, zip.second.Value()}
}

func Zip[T any, U any](first Iterator[T], second Iterator[U]) Iterator[Pair[T, U]] {
	return &ZipIterator[T, U]{first: first, second: second}
}

func ZipLongest[T any, U any](first Iterator[T], second Iterator[U], fill Pair[T, U]) Iterator[Pair[T, U]] {
	return &ZipLongestIterator[T, U]{ZipIterator: ZipIterator[T, U]{first: first, second: second}, fill: fill}
}

type RepeatIterator[T any] struct {
	item T
}

func (r RepeatIterator[T]) Next() bool {
	return true
}

func (r RepeatIterator[T]) Value() T {
	return r.item
}

func Repeat[T any](item T) Iterator[T] {
	return &RepeatIterator[T]{item: item}
}

func RepeatN[T any](item T, n int) Iterator[T] {
	return Take(n, Repeat(item))
}

type CountIterator struct {
	i    int
	step int
}

func (c *CountIterator) Next() bool {
	c.i += c.step
	return true
}

func (c *CountIterator) Value() int {
	return c.i - c.step
}

func Count(start int) Iterator[int] {
	return &CountIterator{i: start, step: 1}
}

func CountBy(start int, step int) Iterator[int] {
	return &CountIterator{i: start, step: step}
}

type ScanIterator[T any] struct {
	source  Iterator[T]
	binOp   func(T, T) T
	value   T
	isFirst bool
}

func (s *ScanIterator[T]) Next() bool {
	if !s.source.Next() {
		return false
	}
	if s.isFirst {
		s.value = s.source.Value()
		s.isFirst = false
	} else {
		s.value = s.binOp(s.value, s.source.Value())
	}

	return true
}

func (s *ScanIterator[T]) Value() T {
	return s.value
}

func Scan[T any](iter Iterator[T], binOp func(T, T) T) Iterator[T] {
	return &ScanIterator[T]{source: iter, binOp: binOp}
}

type ChainIterator[T any] struct {
	iterators []Iterator[T]
	current   int
}

func (c *ChainIterator[T]) Next() bool {
	for ; ; c.current++ {
		if c.current >= len(c.iterators) {
			return false
		}

		if c.iterators[c.current].Next() {
			return true
		}
	}
}

func (c *ChainIterator[T]) Value() T {
	return c.iterators[c.current].Value()
}

func Chain[T any](iterators ...Iterator[T]) Iterator[T] {
	return &ChainIterator[T]{iterators: iterators}
}

func Literal[T any](elem ...T) Iterator[T] {
	return FromSlice(elem)
}

func Compress[T any](data Iterator[T], selectors Iterator[bool]) Iterator[T] {
	zip := Zip(data, selectors)
	var pair Pair[T, bool]
	return &IteratorClosure[T]{
		next: func() bool {
			for {
				if !zip.Next() {
					return false
				}
				pair = zip.Value()
				if pair.second {
					return true
				}
			}
		},
		value: func() T {
			return pair.first
		},
	}
}

func Flatten[T any](iter Iterator[Iterator[T]]) Iterator[T] {
	var current Iterator[T]
	return &IteratorClosure[T]{
		next: func() bool {
			for {
				if current != nil && current.Next() {
					return true
				}
				if !iter.Next() {
					return false
				}
				current = iter.Value()
			}
		},
		value: func() T {
			return current.Value()
		},
	}
}

func MakeSlice[T any](n int, gen func() T) []T {
	slice := make([]T, n)
	for i := range slice {
		slice[i] = gen()
	}
	return slice
}

func Tee[T any](iter Iterator[T], n int) []Iterator[T] {
	deques := MakeSlice[deque.Deque](n, deque.NewDeque)
	makeIterator := func(dq deque.Deque) Iterator[T] {
		var current T
		return &IteratorClosure[T]{
			next: func() bool {
				if dq.Empty() {
					if !iter.Next() {
						return false
					}
					value := iter.Value()
					for _, d := range deques {
						d.PushBack(value)
					}
				}
				current = dq.PopFront().(T)
				return true
			},
			value: func() T {
				return current
			},
		}
	}

	iterators := make([]Iterator[T], n)
	for i, dq := range deques {
		iterators[i] = makeIterator(dq)
	}

	return iterators
}

func Tee2[T any](iter Iterator[T]) (Iterator[T], Iterator[T]) {
	iterators := Tee(iter, 2)
	return iterators[0], iterators[1]
}

func ClosureFromSingle[T any](advance func() (bool, T)) Iterator[T] {
	var value T
	return &IteratorClosure[T]{
		next: func() bool {
			hasNext, newValue := advance()
			value = newValue
			return hasNext
		},
		value: func() T {
			return value
		},
	}
}

func Map[A any, B any](op func(A) B, iter Iterator[A]) Iterator[B] {
	return &IteratorClosure[B]{
		next: iter.Next,
		value: func() B {
			return op(iter.Value())
		},
	}
}

func Tabulate[T any](f func(int) T, start int) Iterator[T] {
	return Map(f, Count(start))
}

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

func Tail[T any](n int, iter Iterator[T]) Iterator[T] {
	ring := NewRingSlice[T](n)
	for iter.Next() {
		ring.Push(iter.Value())
	}
	return ring.Iter()
}

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

func Nth[T any](iter Iterator[T], n int) Optional[T] {
	islice := Drop(n-1, iter)
	if islice.Next() {
		return OptionalWith(islice.Value())
	}
	return EmptyOptional[T]()
}
func Reduce[T any](iter Iterator[T], binOp func(a, b T) T) Optional[T] {
	if !iter.Next() {
		return EmptyOptional[T]()
	}
	value := iter.Value()

	for iter.Next() {
		value = binOp(value, iter.Value())
	}

	return OptionalWith(value)
}

func Prefix[T any](iter Iterator[T], prefix ...T) Iterator[T] {
	return Chain(Literal(prefix...), iter)
}
func Suffix[T any](iter Iterator[T], suffix ...T) Iterator[T] {
	return Chain(iter, Literal(suffix...))
}

func ILen[T any](iter Iterator[T]) int {
	length := 0
	for iter.Next() {
		length++
	}
	return length
}

type Group[T any, K any] struct {
	iter Iterator[T]
	key  K
}

type Key[T any, K any] struct {
	create func(T) K
	equal  func(K, K) bool
}

func MakeKey[T any, K any](key func(T) K, compare func(K, K) bool) Key[T, K] {
	return Key[T, K]{create: key, equal: compare}
}

func GroupByKey[T any, K any](iter Iterator[T], key Key[T, K]) Iterator[Group[T, K]] {
	var currentKey K
	var targetKey K
	var currentGrouper int
	var currentValue T
	isFirst := true
	grouper := func(targetKey K, grouperId int) Iterator[T] {
		doneIteration := false
		next := func() bool {
			if doneIteration || currentGrouper != grouperId || !key.equal(currentKey, targetKey) {
				return false
			}
			if !iter.Next() {
				doneIteration = true
				return true
			}
			currentValue = iter.Value()
			currentKey = key.create(currentValue)
			return true
		}
		value := func() T {
			return currentValue
		}
		return &IteratorClosure[T]{next: next, value: value}
	}

	next := func() bool {
		currentGrouper++

		if isFirst {
			isFirst = false

			if !iter.Next() {
				return false
			}
			currentValue = iter.Value()
			currentKey = key.create(currentValue)
			targetKey = currentKey
			return true
		}

		for key.equal(currentKey, targetKey) {
			if !iter.Next() {
				return false
			}
			currentValue = iter.Value()
			currentKey = key.create(currentValue)
		}
		targetKey = currentKey
		return true
	}
	value := func() Group[T, K] {
		return Group[T, K]{iter: grouper(targetKey, currentGrouper), key: currentKey}
	}

	return &IteratorClosure[Group[T, K]]{next: next, value: value}
}

func GroupByValue[T constraints.Ordered](iter Iterator[T]) Iterator[Group[T, T]] {
	return GroupByKey(
		iter,
		MakeKey(Identity[T], IsSameValue[T]),
	)
}

func Identity[T any](t T) T {
	return t
}

func IsSameValue[T constraints.Ordered](a, b T) bool {
	return a <= b && a >= b
}

func AllEqualValue[T constraints.Ordered](iter Iterator[T]) bool {
	groupBy := GroupByValue(iter)
	groupBy.Next()
	return !groupBy.Next()
}
func AllEqualByKey[T any, K any](iter Iterator[T], key Key[T, K]) bool {
	groupBy := GroupByKey(iter, key)
	groupBy.Next()
	return !groupBy.Next()
}

func EmptyIterator[T any]() Iterator[T] {
	return ClosureFromSingle(func() (bool, T) {
		var value T
		return false, value
	})
}

func Pairwise[T any](iter Iterator[T]) Iterator[Pair[T, T]] {
	a, b := Tee2(iter)
	if !b.Next() {
		return EmptyIterator[Pair[T, T]]()
	}
	return Zip(a, b)
}

func Product[T any](iterators ...Iterator[T]) Iterator[[]T] {
	pools := ToSlice(Map(ToSlice[T], FromSlice(iterators)))
	indices := make([]int, len(iterators))
	done := false
	advance := func() (bool, []T) {
		if done {
			return false, nil
		}
		value := make([]T, len(iterators))
		for slot, index := range indices {
			value[slot] = pools[slot][index]
		}

		carry := true
		for i := len(indices) - 1; i >= 0 && carry; i-- {
			indices[i]++
			if indices[i] >= len(pools[i]) {
				indices[i] = 0
				continue
			}
			carry = false
		}
		done = carry
		return true, value
	}

	return ClosureFromSingle(advance)

}

func TakeWhile[T any](predicate func(T) bool, iter Iterator[T]) Iterator[T] {
	next := func() bool {
		if !iter.Next() {
			return false
		}
		return predicate(iter.Value())
	}

	return IteratorClosure[T]{next: next, value: iter.Value}
}

func DropWhile[T any](predicate func(T) bool, iter Iterator[T]) Iterator[T] {
	next := func() bool {
		for {
			if !iter.Next() {
				return false
			}
			if !predicate(iter.Value()) {
				return true
			}
		}
	}

	return IteratorClosure[T]{next: next, value: iter.Value}
}

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

func EqualTo[T constraints.Ordered](value T) Predicate[T] {
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

/*
Play with things like getting back Len, copying underlying slice directly, avoiding copying out values,
incrementing without value calculation...
*/

/*
Once this is proper lib - experiment with C++ style iterators & ranges.
*/
