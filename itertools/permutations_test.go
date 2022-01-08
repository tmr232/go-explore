package itertools

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

var result []int

func TestPermutationsOf(t *testing.T) {
	t.Run("With Strings", func(t *testing.T) {
		join := func(elements []string) string { return strings.Join(elements, "") }
		type args struct {
			str []string
			r   int
		}
		tests := []struct {
			name string
			args args
			want []string
		}{
			{
				"3 / 3",
				args{[]string{"a", "b", "c"}, 3},
				[]string{"abc", "acb", "bac", "bca", "cab", "cba"},
			},
			{
				"2 / 3",
				args{[]string{"a", "b", "c"}, 2},
				[]string{"ab", "ac", "ba", "bc", "ca", "cb"},
			},
			{
				"1 / 3",
				args{[]string{"a", "b", "c"}, 1},
				[]string{"a", "b", "c"},
			},

			{
				"4 / 3",
				args{[]string{"a", "b", "c"}, 4},
				[]string{},
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if got := ToSlice(Map[[]string](join, PermutationsOf(tt.args.str, tt.args.r))); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("GenericGetPermutationsOf() = %v, want %v", got, tt.want)
				}
			})
		}
	})

	t.Run("With Ints", func(t *testing.T) {
		type args struct {
			input []int
			r     int
		}
		tests := []struct {
			name string
			args args
			want [][]int
		}{
			{
				"3 / 3",
				args{[]int{1, 2, 3}, 3},
				[][]int{{1, 2, 3}, {1, 3, 2}, {2, 1, 3}, {2, 3, 1}, {3, 1, 2}, {3, 2, 1}},
			},
			{
				"2 / 3",
				args{[]int{1, 2, 3}, 2},
				[][]int{{1, 2}, {1, 3}, {2, 1}, {2, 3}, {3, 1}, {3, 2}},
			},
			{
				"1 / 3",
				args{[]int{1, 2, 3}, 1},
				[][]int{{1}, {2}, {3}},
			},
			{
				"4 / 3",
				args{[]int{1, 2, 3}, 4},
				[][]int{},
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if got := ToSlice(CopySlices[int](PermutationsOf(tt.args.input, tt.args.r))); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("GenericGetPermutationsOf() = %v, want %v", got, tt.want)
				}
			})
		}
	})
}

func TestSafePermutationsOf(t *testing.T) {
	t.Run("With Strings", func(t *testing.T) {
		join := func(elements []string) string { return strings.Join(elements, "") }
		type args struct {
			str []string
			r   int
		}
		tests := []struct {
			name string
			args args
			want []string
		}{
			{
				"3 / 3",
				args{[]string{"a", "b", "c"}, 3},
				[]string{"abc", "acb", "bac", "bca", "cab", "cba"},
			},
			{
				"2 / 3",
				args{[]string{"a", "b", "c"}, 2},
				[]string{"ab", "ac", "ba", "bc", "ca", "cb"},
			},
			{
				"1 / 3",
				args{[]string{"a", "b", "c"}, 1},
				[]string{"a", "b", "c"},
			},

			{
				"4 / 3",
				args{[]string{"a", "b", "c"}, 4},
				[]string{},
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if got := ToSlice(Map[[]string](join, SafePermutationsOf(tt.args.str, tt.args.r))); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("GenericGetPermutationsOf() = %v, want %v", got, tt.want)
				}
			})
		}
	})

	t.Run("With Ints", func(t *testing.T) {
		type args struct {
			input []int
			r     int
		}
		tests := []struct {
			name string
			args args
			want [][]int
		}{
			{
				"3 / 3",
				args{[]int{1, 2, 3}, 3},
				[][]int{{1, 2, 3}, {1, 3, 2}, {2, 1, 3}, {2, 3, 1}, {3, 1, 2}, {3, 2, 1}},
			},
			{
				"2 / 3",
				args{[]int{1, 2, 3}, 2},
				[][]int{{1, 2}, {1, 3}, {2, 1}, {2, 3}, {3, 1}, {3, 2}},
			},
			{
				"1 / 3",
				args{[]int{1, 2, 3}, 1},
				[][]int{{1}, {2}, {3}},
			},
			{
				"4 / 3",
				args{[]int{1, 2, 3}, 4},
				[][]int{},
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if got := ToSlice(SafePermutationsOf(tt.args.input, tt.args.r)); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("GenericGetPermutationsOf() = %v, want %v", got, tt.want)
				}
			})
		}
	})
}

func benchmarkPermutationsOf(input []int, r int, b *testing.B) {
	for i := 0; i < b.N; i++ {
		permutations := PermutationsOf(input, r)
		for permutations.Next() {
			result = permutations.Value()
		}
	}
}

func BenchmarkPermutationsOf(b *testing.B) {
	for i := 1; i <= 5; i++ {
		b.Run(fmt.Sprint("r =", i), func(b *testing.B) {
			benchmarkPermutationsOf([]int{1, 2, 3, 4, 5, 6, 7, 8, 9}, i, b)
		})
	}
}
