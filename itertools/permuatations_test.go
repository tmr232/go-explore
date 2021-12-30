package itertools

import (
	"fmt"
	"testing"
)

var result []int

func benchmarkGenericGetPermutationsOf(input []int, r int, b *testing.B) {
	for i := 0; i < b.N; i++ {
		genericPermutations := GenericGetPermutationsOf(input, r)
		for genericPermutations.Next() {
			result = genericPermutations.Value()
		}
	}
}

func BenchmarkGenericGetPermutationsOf(b *testing.B) {
	for i := 1; i <= 5; i++ {
		b.Run(fmt.Sprint("r =", i), func(b *testing.B) {
			benchmarkGenericGetPermutationsOf([]int{1, 2, 3, 4, 5, 6, 7, 8, 9}, i, b)
		})
	}
}

func benchmarkGetPermutationsOf(input []int, r int, b *testing.B) {
	for i := 0; i < b.N; i++ {
		genericPermutations := GetPermutationsOf(input, r)
		for genericPermutations.Next() {
			result = genericPermutations.Value().([]int)
		}
	}
}

func BenchmarkReflectionGetPermutationsOf(b *testing.B) {
	for i := 1; i <= 5; i++ {
		b.Run(fmt.Sprint("r =", i), func(b *testing.B) {
			benchmarkGetPermutationsOf([]int{1, 2, 3, 4, 5, 6, 7, 8, 9}, i, b)
		})
	}
}

func TestGenericGetPermutationsOf(t *testing.T) {
	input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := 4
	for i := 0; i < 1000; i++ {
		genericPermutations := GenericGetPermutationsOf(input, r)
		for genericPermutations.Next() {
			result = genericPermutations.Value()
		}
	}
}

func TestReflectionGetPermutationsOf(t *testing.T) {
	input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := 4
	for i := 0; i < 1000; i++ {
		genericPermutations := GetPermutationsOf(input, r)
		for genericPermutations.Next() {
			result = genericPermutations.Value().([]int)
		}
	}
}
