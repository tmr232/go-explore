package itertools

import (
	"github.com/edwingeng/deque"
)

// TODO it seems that for completeness we need Iterable[T] as well, that implements `Iter()`.
// TODO for Iterator[T], it should just return itself.
// TODO this should allow easier interop with custom slice and container types.

type Iterator[T any] interface {
	// Next tries to advance to the next value.
	// Returns true if a value exists, false if not.
	// Once an iterator returns false to indicate exhaustion,
	// it should continue returning false.
	Next() bool
	// Value returns the current value of the iterator.
	// Next() must be called and return true before every call to Value().
	Value() T
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

// FromAdvance returns an iterator built using the provided advance function.
// The advance function is called once per iteration, returning a flag indicating
// the existence of a value, and the value itself.
// When implementing your own, return (true, value) as long as values are available,
// and (false, *new(T)) when the iterator is exhausted.
func FromAdvance[T any](advance func() (hasValue bool, value T)) Iterator[T] {
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

// FromAdvanceSafe is like FromAdvance, but ensure that once (false, _) is returned,
// the provided advance function will not get called, and all consecutive calls to `.Next()`
// will return false.
func FromAdvanceSafe[T any](advance func() (hasValue bool, value T)) Iterator[T] {
	var value T
	var ic IteratorClosure[T]
	ic = IteratorClosure[T]{
		next: func() bool {
			hasNext, newValue := advance()
			if !hasNext {
				ic.next = func() bool { return false }
			}
			value = newValue
			return hasNext
		},
		value: func() T {
			return value
		},
	}
	return &ic
}

type filterInIterator[T any] struct {
	iter     Iterator[T]
	filterIn func(T) bool
}

func (f *filterInIterator[T]) Next() bool {
	for f.iter.Next() {
		if f.filterIn(f.iter.Value()) {
			return true
		}
	}
	return false
}

func (f *filterInIterator[T]) Value() T {
	return f.iter.Value()
}

// FilterIn keeps only elements matched by a predicate
func FilterIn[T any](iter Iterator[T], filterIn func(T) bool) Iterator[T] {
	return &filterInIterator[T]{iter: iter, filterIn: filterIn}
}

// FilterOut skips all elements matched by a predicate
func FilterOut[T any](iter Iterator[T], filterOut func(T) bool) Iterator[T] {
	return FilterIn(iter, Not(filterOut))
}

type intRangeIterator struct {
	stop    int
	current int
}

func (i *intRangeIterator) Next() bool {
	i.current++
	return i.current < i.stop
}

func (i *intRangeIterator) Value() int {
	return i.current
}

// IntRange iterates over the range [0, stop)
func IntRange(stop int) Iterator[int] {
	return &intRangeIterator{stop: stop, current: -1}
}

type cycleIterator[T any] struct {
	iter     Iterator[T]
	slice    []T
	dq       deque.Deque
	size     int
	index    int
	consumed bool
	value    T
}

func (ci *cycleIterator[T]) Next() bool {
	if !ci.consumed {
		if ci.iter.Next() {
			ci.value = ci.iter.Value()
			ci.dq.PushBack(ci.value)
			ci.size++
			return true
		} else {
			ci.consumed = true
			ci.slice = dequeToSlice[T](ci.dq)
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

func (ci *cycleIterator[T]) Value() T {
	return ci.value
}

// Cycle iterates over the elements of iter, repeating when it reaches the end.
// Cycle will store an internal copy of all iterated elements.
func Cycle[T any](iter Iterator[T]) Iterator[T] {
	return &cycleIterator[T]{iter: iter, index: -1, dq: deque.NewDeque()}
}

// Consume iterates through n items from the input iterator.
// Consume returns true if it succeeded, false if the iterator was consumed before completion.
func Consume[T any](n int, iter Iterator[T]) bool {
	for i := 0; i < n; i++ {
		if !iter.Next() {
			return false
		}
	}
	return true
}

func ISlice[T any](iter Iterator[T], options ...RangeOption) Iterator[T] {
	start, stop, step := getConfig(options).Unpack()
	if stop <= start || step <= 0 || start < 0 || stop < 0 {
		return EmptyIterator[T]()
	}

	var resultIterator IteratorClosure[T]

	i := start

	next := func() bool {
		i += step
		return i < stop && Consume(step, iter)
	}

	firstNext := func() bool {
		if Consume(start, iter) && iter.Next() {
			resultIterator.next = next
			return true
		}
		resultIterator.next = func() bool { return false }
		return false
	}

	resultIterator.next = firstNext
	resultIterator.value = iter.Value

	return &resultIterator
}

// Take yields the first n elements of the input iterator.
func Take[T any](n int, iter Iterator[T]) Iterator[T] {
	return ISlice(iter, Stop(n))
}

// Drop skips the first n elements of the input iterator
func Drop[T any](n int, iter Iterator[T]) Iterator[T] {
	return ISlice(iter, Start(n))
}

type zipIterator[T any, U any] struct {
	first  Iterator[T]
	second Iterator[U]
}

func (zip *zipIterator[T, U]) Next() bool {
	hasFirst := zip.first.Next()
	hasSecond := zip.second.Next()

	return hasFirst && hasSecond
}

func (zip *zipIterator[T, U]) Value() Pair[T, U] {
	return Pair[T, U]{zip.first.Value(), zip.second.Value()}
}

type zipLongestIterator[T any, U any] struct {
	zipIterator[T, U]
	fill      Pair[T, U]
	hasFirst  bool
	hasSecond bool
}

func (zip *zipLongestIterator[T, U]) Next() bool {
	zip.hasFirst = zip.first.Next()
	zip.hasSecond = zip.second.Next()

	return zip.hasFirst || zip.hasSecond
}

func (zip *zipLongestIterator[T, U]) Value() Pair[T, U] {
	if zip.hasFirst && zip.hasSecond {
		return Pair[T, U]{zip.first.Value(), zip.second.Value()}
	}
	if zip.hasFirst {
		return Pair[T, U]{zip.first.Value(), zip.fill.second}
	}
	return Pair[T, U]{zip.fill.first, zip.second.Value()}
}

func Zip[T any, U any](first Iterator[T], second Iterator[U]) Iterator[Pair[T, U]] {
	return &zipIterator[T, U]{first: first, second: second}
}

func ZipLongest[T any, U any](first Iterator[T], second Iterator[U], fill Pair[T, U]) Iterator[Pair[T, U]] {
	return &zipLongestIterator[T, U]{zipIterator: zipIterator[T, U]{first: first, second: second}, fill: fill}
}

type repeatIterator[T any] struct {
	item T
}

func (r repeatIterator[T]) Next() bool {
	return true
}

func (r repeatIterator[T]) Value() T {
	return r.item
}

// Repeat yields the input value infinite times.
func Repeat[T any](item T) Iterator[T] {
	return &repeatIterator[T]{item: item}
}

// RepeatN repeats the input value n times.
func RepeatN[T any](item T, n int) Iterator[T] {
	return Take(n, Repeat(item))
}

type countIterator struct {
	i    int
	step int
}

func (c *countIterator) Next() bool {
	c.i += c.step
	return true
}

func (c *countIterator) Value() int {
	return c.i - c.step
}

// Count yields consecutive increasing numbers from start
func Count(start int) Iterator[int] {
	return &countIterator{i: start, step: 1}
}

// CountBy yields numbers starting with start and incrementing by step
func CountBy(start int, step int) Iterator[int] {
	return &countIterator{i: start, step: step}
}

type scanIterator[T any] struct {
	source  Iterator[T]
	binOp   func(T, T) T
	value   T
	isFirst bool
}

func (s *scanIterator[T]) Next() bool {
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

func (s *scanIterator[T]) Value() T {
	return s.value
}

// Scan performce a stepwise reduction over the input iterator, yielding all intermediate values.
func Scan[T any](iter Iterator[T], binOp func(T, T) T) Iterator[T] {
	return &scanIterator[T]{source: iter, binOp: binOp}
}

type chainIterator[T any] struct {
	iterators []Iterator[T]
	current   int
}

func (c *chainIterator[T]) Next() bool {
	for ; ; c.current++ {
		if c.current >= len(c.iterators) {
			return false
		}

		if c.iterators[c.current].Next() {
			return true
		}
	}
}

func (c *chainIterator[T]) Value() T {
	return c.iterators[c.current].Value()
}

// Chain yields elements from the first input iterator until it is exhausted, then from the
// second input iterator, the third, and so forth.
func Chain[T any](iterators ...Iterator[T]) Iterator[T] {
	return &chainIterator[T]{iterators: iterators}
}

// Literal yields the values it is created with.
func Literal[T any](elem ...T) Iterator[T] {
	return FromSlice(elem)
}

// Compress yields data elements corresponding to true selector elements.
// Forms a shorter iterator from selected data elements using the selectors to
// choose the data elements.
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

// Flatten removes one level of nesting in an iterator of iterators.
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

// FlattenSlices removes one level of nesting in an iterator of iterators.
func FlattenSlices[T any](iter Iterator[[]T]) Iterator[T] {
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
				current = FromSlice(iter.Value())
			}
		},
		value: func() T {
			return current.Value()
		},
	}
}

// Tee returns a slice of n independent iterators with the same content as the input iterator.
func Tee[T any](iter Iterator[T], n int) []Iterator[T] {
	deques := makeSlice[deque.Deque](n, deque.NewDeque)
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

// Tee2 returns two independent iterators with the same content as the input iterator.
func Tee2[T any](iter Iterator[T]) (Iterator[T], Iterator[T]) {
	iterators := Tee(iter, 2)
	return iterators[0], iterators[1]
}

// Map returns an iterator whose elements are the result of calling op on the elements
// of the input iterator.
// Note that op is only guaranteed to be called when getting a value, not in pure iteration.
func Map[A any, B any](op func(A) B, iter Iterator[A]) Iterator[B] {
	return &IteratorClosure[B]{
		next: iter.Next,
		value: func() B {
			return op(iter.Value())
		},
	}
}

/* TODO use codegen for generating Map functions that take more iterables of different element types.
Map2 takes 2 iterators, Map3 takes 3...
*/

func Tabulate[T any](f func(int) T, start int) Iterator[T] {
	return Map(f, Count(start))
}

// Tail returns an iterator with the last n elements of the input iterator.
func Tail[T any](n int, iter Iterator[T]) Iterator[T] {
	ring := NewRingSlice[T](n)
	for iter.Next() {
		ring.Push(iter.Value())
	}
	return ring.Iter()
}

// Nth gets the n'th element if an iterator, if one exists.
func Nth[T any](iter Iterator[T], n int) (value T, exists bool) {
	islice := ISlice(iter, Start(n), Stop(n+1))
	if islice.Next() {
		return islice.Value(), true
	}
	return *new(T), false
}

// Reduce performs a reduction over the values of the input iterator, if there are such values.
func Reduce[T any](iter Iterator[T], binOp func(a, b T) T) (value T, exists bool) {
	if !iter.Next() {
		return *new(T), false
	}
	value = iter.Value()

	for iter.Next() {
		value = binOp(value, iter.Value())
	}

	return value, true
}

func Fold[Acc any, Elem any](init Acc, iter Iterator[Elem], op func(acc Acc, elem Elem) Acc) Acc {
	acc := init
	for iter.Next() {
		acc = op(acc, iter.Value())
	}
	return acc
}

func Prefix[T any](iter Iterator[T], prefix ...T) Iterator[T] {
	return Chain(Literal(prefix...), iter)
}

func Suffix[T any](iter Iterator[T], suffix ...T) Iterator[T] {
	return Chain(iter, Literal(suffix...))
}

// ILen returns the number of elements in an iterator. Consumes the iterator.
func ILen[T any](iter Iterator[T]) int {
	length := 0
	for iter.Next() {
		length++
	}
	return length
}

type Group[T any, K any] struct {
	slice []T
	key   K
}

type Key[T any, K any] struct {
	create func(T) K
	equal  func(K, K) bool
}

func MakeKey[T any, K any](key func(T) K, compare func(K, K) bool) Key[T, K] {
	return Key[T, K]{create: key, equal: compare}
}

func GroupByKey[T any, K any](iter Iterator[T], key Key[T, K]) Iterator[Group[T, K]] {
	if !iter.Next() {
		return EmptyIterator[Group[T, K]]()
	}

	currentValue := iter.Value()
	currentKey := key.create(currentValue)
	dq := deque.NewDeque()
	dq.PushBack(currentValue)

	fromDeque := func() (bool, Group[T, K]) {
		if dq.Empty() {
			var zero Group[T, K]
			return false, zero
		}
		slice := dequeToSlice[T](dq)
		return true, Group[T, K]{slice, currentKey}
	}

	advance := func() (bool, Group[T, K]) {
		for iter.Next() {
			newValue := iter.Value()
			newKey := key.create(newValue)
			if key.equal(newKey, currentKey) {
				dq.PushBack(newValue)
				currentKey = newKey
			} else {
				next, group := fromDeque()
				currentKey = newKey
				dq.PushBack(newValue)
				return next, group
			}
		}
		return fromDeque()
	}

	return FromAdvance(advance)
}

func GroupByValue[T comparable](iter Iterator[T]) Iterator[Group[T, T]] {
	return GroupByKey(
		iter,
		MakeKey(Identity[T], IsSameValue[T]),
	)
}

func GroupByFunc[T any, K comparable](iter Iterator[T], key func(T) K) Iterator[Group[T, K]] {
	return GroupByKey(iter, MakeKey(key, IsSameValue[K]))
}

func AllEqualValue[T comparable](iter Iterator[T]) bool {
	groupBy := GroupByValue(iter)
	groupBy.Next()
	return !groupBy.Next()
}

func AllEqualByKey[T any, K any](iter Iterator[T], key Key[T, K]) bool {
	groupBy := GroupByKey(iter, key)
	groupBy.Next()
	return !groupBy.Next()
}

// EmptyIterator returns an empty iterator.
// It's .Next method will always return false.
func EmptyIterator[T any]() Iterator[T] {
	return FromAdvance(func() (bool, T) {
		var value T
		return false, value
	})
}

// Pairwise yields pairs of consecutive values from the input iterator.
func Pairwise[T any](iter Iterator[T]) Iterator[Pair[T, T]] {
	a, b := Tee2(iter)
	if !b.Next() {
		return EmptyIterator[Pair[T, T]]()
	}
	return Zip(a, b)
}

// Product is the Cartesian product of the input iterators.
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

	return FromAdvance(advance)
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

func FromCallable[T any](f func() T) Iterator[T] {
	return FromAdvance(func() (bool, T) {
		return true, f()
	})
}

func FromCallableWhile[T any](f func() T, predicate Predicate[T]) Iterator[T] {
	return TakeWhile(predicate, FromCallable(f))
}

func FromCallableUntil[T any](f func() T, predicate Predicate[T]) Iterator[T] {
	return TakeWhile(Not(predicate), FromCallable(f))
}

func WindowedWithFiller[T any](iter Iterator[T], n int, filler T) Iterator[[]T] {
	ring := NewRingSlice[T](n)
	first := true
	advance := func() (bool, []T) {
		if first {
			ForEach(Take(n, Chain(iter, Repeat(filler))), ring.Push)
			first = false
			return true, ring.ToSlice()
		}
		if !iter.Next() {
			return false, nil
		}
		ring.Push(iter.Value())
		return true, ring.ToSlice()
	}

	return FromAdvance(advance)
}
func Windowed[T any](iter Iterator[T], n int) Iterator[[]T] {
	ring := NewRingSlice[T](n)
	first := true
	advance := func() (bool, []T) {
		if first {
			ForEach(Take(n, Chain(iter)), ring.Push)
			first = false
			result := ring.ToSlice()
			if len(result) < n {
				return false, nil
			}
			return true, result
		}
		if !iter.Next() {
			return false, nil
		}
		ring.Push(iter.Value())
		return true, ring.ToSlice()
	}

	return FromAdvance(advance)
}

func Chunked[T any](iter Iterator[T], n int) Iterator[[]T] {
	return FromCallableWhile(func() []T { return ToSlice(Take(n, iter)) }, SliceNotEmpty[T])
}

func ChunkBy[T any, K comparable](iter Iterator[T], key func(T) K) Iterator[[]T] {
	return Map(
		func(g Group[T, K]) []T {
			return g.slice
		},
		GroupByFunc(iter, key),
	)
}

// EnumerateFrom yields pairs containing a count (from start)
// and a value yielded by the iterator argument.
func EnumerateFrom[T any](iter Iterator[T], start int) Iterator[Pair[int, T]] {
	return Zip(Count(start), iter)
}

// Enumerate  yields pairs containing a count (from zero)
// and a value yielded by the iterator argument.
func Enumerate[T any](iter Iterator[T]) Iterator[Pair[int, T]] {
	return EnumerateFrom(iter, 0)
}

// ForEach executes f on every element of the input iterator.
func ForEach[T any](iter Iterator[T], f func(T)) {
	for iter.Next() {
		f(iter.Value())
	}
}

// RoundRobin cycles through all input iterators yielding an element from each, until all are exhausted.
func RoundRobin[T any](iterators ...Iterator[T]) Iterator[T] {
	indices := Cycle(IntRange(len(iterators)))
	active := ToSlice(Take(len(iterators), Repeat(true)))
	activeCount := len(iterators)
	advance := func() (bool, T) {
		for indices.Next() && activeCount > 0 {
			index := indices.Value()
			if active[index] {
				if !iterators[index].Next() {
					active[index] = false
					activeCount--
				} else {
					return true, iterators[index].Value()
				}
			}
		}
		var zero T
		return false, zero
	}
	return FromAdvance(advance)
}

// InterleaveFlat returns
func InterleaveFlat[T any](iterators ...Iterator[T]) Iterator[T] {
	size := len(iterators)
	slice := make([]T, size)
	var zero T
	k := 0
	advance := func() (bool, T) {
		if k == 0 {
			for i, iter := range iterators {
				if !iter.Next() {
					return false, zero
				}
				slice[i] = iter.Value()
			}
		}

		value := slice[k]
		k = (k + 1) % size
		return true, value
	}

	return FromAdvance(advance)
}

func Interleave[T any](iterators ...Iterator[T]) Iterator[[]T] {
	size := len(iterators)
	slice := make([]T, size)

	advance := func() (bool, []T) {
		for i, iter := range iterators {
			if !iter.Next() {
				return false, nil
			}
			slice[i] = iter.Value()
		}
		return true, slice
	}

	return FromAdvance(advance)
}

func InterleaveLongest[T any](filler T, iterators ...Iterator[T]) Iterator[[]T] {
	size := len(iterators)
	slice := make([]T, size)

	advance := func() (bool, []T) {
		hasValues := false
		for i, iter := range iterators {
			if iter.Next() {
				slice[i] = iter.Value()
				hasValues = true
			} else {
				slice[i] = filler
			}
		}
		return hasValues, slice
	}

	return FromAdvance(advance)
}

// AllButLast splits the input iterator into two iterators.
// The first iterator will yield all values except the last n.
// The second will yield the last n items.
// The first iterator must be fully consumed before iterating over the second iterator.
func AllButLast[T any](iter Iterator[T], n int) (first Iterator[T], last Iterator[T]) {
	lastItems := make([]T, n)
	i := 0
	var zero T
	firstAdvance := func() (bool, T) {
		// Load the n last items
		for ; i < n; i++ {
			if !iter.Next() {
				return false, zero
			}
			lastItems[i] = iter.Value()
		}
		// Store new one, return oldest
		if !iter.Next() {
			return false, zero
		}
		value := lastItems[i%n]
		lastItems[i%n] = iter.Value()
		i++
		return true, value
	}
	k := 0
	lastAdvance := func() (bool, T) {
		length := Min(i, n)

		if k >= length {
			return false, zero
		}
		value := lastItems[(i+k)%n]
		k++
		return true, value
	}

	return FromAdvance(firstAdvance), FromAdvance(lastAdvance)
}

// Intersperse yields the input value between every 2 elements of the input iterator.
func Intersperse[T any](iter Iterator[T], value T) Iterator[T] {
	first, last := AllButLast(iter, 1)
	return Chain(InterleaveFlat(first, Repeat(value)), last)
}

// Join yields the elements of iter separated by the separator
func Join[T any](iter Iterator[T], separator T) Iterator[T] {
	return Intersperse(iter, separator)
}

func Unzip[A, B any](iter Iterator[Pair[A, B]]) (Iterator[A], Iterator[B]) {
	iterA, iterB := Tee2(iter)
	return Map(Pair[A, B].First, iterA), Map(Pair[A, B].Second, iterB)
}

func Deinterleave[T any](iter Iterator[[]T]) []Iterator[T] {
	if !iter.Next() {
		return nil
	}
	first := iter.Value()
	size := len(first)
	iters := make([]Iterator[T], size)
	for i, it := range Tee(iter, size) {
		iters[i] = Prefix(Map(ItemGetter[T](i), it), first[i])
	}
	return iters
}

/*
Play with things like getting back Len, copying underlying slice directly, avoiding copying out values,
incrementing without value calculation...
*/

/*
Once this is proper lib - experiment with C++ style iterators & ranges.
*/
