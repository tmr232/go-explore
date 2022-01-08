package itertools

import (
	"reflect"
	"testing"
)

func TestRange(t *testing.T) {
	type args struct {
		options []RangeOption
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			"range(10)",
			args{[]RangeOption{Stop(10)}},
			[]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
		},
		{
			"range(5,10)",
			args{[]RangeOption{Start(5), Stop(10)}},
			[]int{5, 6, 7, 8, 9},
		},
		{
			"range(5, 10, 2)",
			args{[]RangeOption{Start(5), Stop(10), Step(2)}},
			[]int{5, 7, 9},
		},

		{
			"range(stop=-10, step=-1)",
			args{[]RangeOption{Stop(-10), Step(-1)}},
			[]int{0, -1, -2, -3, -4, -5, -6, -7, -8, -9},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToSlice(Range(tt.args.options...)); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Range() = %v, want %v", got, tt.want)
			}
		})
	}
}
