package itertools

import (
	"reflect"
	"testing"
)

func TestIndexPermutations(t *testing.T) {
	type args struct {
		n int
		r int
	}
	tests := []struct {
		name string
		args args
		want [][]int
	}{
		{
			"3 / 3",
			args{3, 3},
			[][]int{{0, 1, 2}, {0, 2, 1}, {1, 0, 2}, {1, 2, 0}, {2, 0, 1}, {2, 1, 0}},
		},
		{
			"1 / 3",
			args{3, 2},
			[][]int{{0, 1}, {0, 2}, {1, 0}, {1, 2}, {2, 0}, {2, 1}},
		},
		{
			"1 / 3",
			args{3, 1},
			[][]int{{0}, {1}, {2}},
		},
		{
			"4 / 3",
			args{3, 4},
			[][]int{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToSlice(CopySlices[int](IndexPermutations(tt.args.n, tt.args.r))); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IndexPermutations() = %v, want %v", got, tt.want)
			}
		})
	}
}
