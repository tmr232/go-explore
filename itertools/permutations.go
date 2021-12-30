package itertools

import (
	"fmt"
	"reflect"
)

func SimpleRange(n int) []int {
	r := make([]int, n)
	for i := range r {
		r[i] = i
	}
	return r
}

func ReverseRange(start, stop int) []int {
	r := make([]int, start-stop)
	v := start
	for i := range r {
		r[i] = v
		v--
	}
	return r
}

func Rotate(slice []int, toStart int) {
	if toStart == 0 {
		return
	}

	read := toStart
	write := 0
	nextRead := 0
	last := len(slice)
	for read != last {
		if write == nextRead {
			nextRead = read
		}
		slice[write], slice[read] = slice[read], slice[write]
		read += 1
		write += 1
	}
	Rotate(slice[write:], nextRead-write)
}

type PermutationState struct {
	pool        []int
	r           int
	first       bool
	indices     []int
	cycles      []int
	valid       bool
	currentPerm []int
}

func (state *PermutationState) Value() []int {
	return state.currentPerm
}

func (state *PermutationState) Next() bool {
	r := state.r
	pool := state.pool
	indices := state.indices
	cycles := state.cycles
	n := len(pool)

	if !state.valid {
		return false
	}

	if state.first {
		state.first = false
		copy(state.currentPerm, pool[:r])
		return true
	}

	for i := r - 1; i >= 0; i-- {
		cycles[i] -= 1

		if cycles[i] == 0 {
			Rotate(indices[i:], 1)
			cycles[i] = n - i
		} else {
			j := cycles[i]
			indices[i], indices[len(indices)-j] = indices[len(indices)-j], indices[i]
			perm := state.currentPerm
			for pi, k := range indices[:r] {
				perm[pi] = pool[k]
			}
			return true
		}
	}
	state.valid = false
	return false
}

func Permutations(slice []int, r int) *PermutationState {
	n := len(slice)
	return &PermutationState{
		pool:        append([]int{}, slice...),
		r:           r,
		first:       true,
		indices:     SimpleRange(n),
		cycles:      ReverseRange(n, n-r),
		valid:       r <= n,
		currentPerm: make([]int, r),
	}
}

func GetPermutation(slice interface{}, permutation []int) interface{} {
	rawSlice := reflect.ValueOf(slice)

	if rawSlice.Len() < len(permutation) {
		return nil
	}

	result := reflect.MakeSlice(reflect.TypeOf(slice), len(permutation), len(permutation))
	for write, read := range permutation {
		result.Index(write).Set(rawSlice.Index(read))
	}

	return result.Interface()
}

func GenericGetPermutation[T any](slice []T, permutation []int) []T {
	if len(slice) < len(permutation) {
		return nil
	}
	result := make([]T, len(permutation))
	for write, read := range permutation {
		result[write] = slice[read]
	}
	return result
}

type PermutationsOf struct {
	indexPermutations *PermutationState
	rawSlice          reflect.Value
	output            reflect.Value
}

type GenericPermutationsOf[T any] struct {
	indexPermutations *PermutationState
	slice             []T
	output            []T
}

func (state *PermutationsOf) Next() bool {
	return state.indexPermutations.Next()
}

func (state *GenericPermutationsOf[T]) Next() bool {
	return state.indexPermutations.Next()
}

func (state *PermutationsOf) Value() interface{} {
	indices := state.indexPermutations.Value()
	if indices == nil {
		return nil
	}
	state.applyPermutation()
	return state.output.Interface()
}

func (state *PermutationsOf) applyPermutation() {
	for write, read := range state.indexPermutations.Value() {
		state.output.Index(write).Set(state.rawSlice.Index(read))
	}
}

func (state *GenericPermutationsOf[T]) applyPermutation() {
	for write, read := range state.indexPermutations.Value() {
		state.output[write] = state.slice[read]
	}
}

func (state *GenericPermutationsOf[T]) Value() []T {
	indices := state.indexPermutations.Value()
	if indices == nil {
		return nil
	}
	state.applyPermutation()
	return state.output
}

func GetPermutationsOf(slice interface{}, r int) *PermutationsOf {
	rawSlice := reflect.ValueOf(slice)
	length := rawSlice.Len()

	indexPermutations := Permutations(SimpleRange(length), r)

	return &PermutationsOf{
		indexPermutations: indexPermutations,
		rawSlice:          reflect.ValueOf(slice),
		output:            reflect.MakeSlice(reflect.TypeOf(slice), r, r),
	}
}

func GenericGetPermutationsOf[T any](slice []T, r int) *GenericPermutationsOf[T] {
	indexPermutations := Permutations(SimpleRange(len(slice)), r)
	return &GenericPermutationsOf[T]{indexPermutations: indexPermutations, slice: slice, output: make([]T, r)}
}

func main() {
	permutations := GetPermutationsOf([]string{"a", "b", "c"}, 3)
	for permutations.Next() {
		fmt.Println(permutations.Value().([]string))
	}
	fmt.Println("Yoohoo!")
	genericPermutations := GenericGetPermutationsOf([]string{"a", "b", "c"}, 3)
	for genericPermutations.Next() {
		fmt.Println(genericPermutations.Value())
	}
}
