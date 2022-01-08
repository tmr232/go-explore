package itertools

import (
	"fmt"
	"reflect"
	"testing"
)

func TestIntRange(t *testing.T) {
	type args struct {
		stop int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{"0", args{}, []int{}},
		{"1", args{1}, []int{0}},
		{"10", args{10}, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToSlice(IntRange(tt.args.stop)); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IntRange() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilterIn(t *testing.T) {
	got := ToSlice(FilterIn(IntRange(10), func(n int) bool {
		return n%2 == 0
	}))

	want := []int{0, 2, 4, 6, 8}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got = %v, want %v", got, want)
	}
}

func TestFilterOut(t *testing.T) {
	got := ToSlice(FilterOut(IntRange(10), func(n int) bool {
		return n%2 == 0
	}))

	want := []int{1, 3, 5, 7, 9}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got = %v, want %v", got, want)
	}
}

func TestFilter(t *testing.T) {
	want := ToSlice(IntRange(10))
	filter := func(n int) bool {
		return n%2 == 0
	}
	got := ToSlice(InterleaveFlat(FilterIn(IntRange(10), filter), FilterOut(IntRange(10), filter)))
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got = %v, want %v", got, want)
	}
}

func TestTake(t *testing.T) {
	want := []int{0, 1, 2, 3}

	got := ToSlice(Take(4, IntRange(10)))

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got = %v, want %v", got, want)
	}
}

func TestCycle(t *testing.T) {
	want := []int{0, 1, 2, 3, 0, 1, 2, 3, 0, 1}

	got := ToSlice(Take(10, Cycle(Take(4, IntRange(10)))))

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got = %v, want %v", got, want)
	}
}

func TestZip(t *testing.T) {
	want := []Pair[int, string]{{1, "a"}, {2, "b"}, {3, "c"}}

	got := ToSlice(Zip[int, string](FromSlice([]int{1, 2, 3}), FromSlice([]string{"a", "b", "c"})))

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got = %v, want %v", got, want)
	}
}

func TestZipLongest(t *testing.T) {
	want := []Pair[int, string]{{1, "a"}, {2, "b"}, {3, "b"}}

	got := ToSlice(ZipLongest[int, string](FromSlice([]int{1, 2, 3}), FromSlice([]string{"a"}), Pair[int, string]{0, "b"}))

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got = %v, want %v", got, want)
	}
}

func TestReflectionThingy(t *testing.T) {
	v := 1
	vv := reflect.ValueOf(v)
	vvv := vv.Interface().(int)
	fmt.Println(vvv)
}

func ConsumeA[T any](n int, iter Iterator[T]) {
	for i := 0; i < n; i++ {
		if !iter.Next() {
			break
		}
	}
}

func ConsumeB[T any, I Iterator[T]](n int, iter I) {
	for i := 0; i < n; i++ {
		if !iter.Next() {
			break
		}
	}
}

func TestConsumeA(t *testing.T) {
	source := Repeat(0)
	ConsumeA[int](100, source)
}

func BenchmarkConsumeA(b *testing.B) {
	source := Repeat(0)
	for i := 0; i < b.N; i++ {
		ConsumeA[int](1, source)
	}
}

func BenchmarkConsumeB(b *testing.B) {
	source := Repeat(0)
	for i := 0; i < b.N; i++ {
		ConsumeB[int](1, source)
	}
}

func BenchmarkConsumeC(b *testing.B) {
	source := func() Iterator[int] { return Repeat(0) }()
	for i := 0; i < b.N; i++ {
		ConsumeB[int](1, source)
	}
}

func TestScan(t *testing.T) {
	want := []int{1, 3, 6, 10, 15}
	got := ToSlice(Scan(FromSlice([]int{1, 2, 3, 4, 5}), func(a, b int) int { return a + b }))
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got = %v, want %v", got, want)
	}
}

func TestChain(t *testing.T) {
	want := []int{1, 2, 5, 6, 9}
	got := ToSlice(Chain(Literal(1, 2), Literal(5, 6), Literal(9)))
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got = %v, want %v", got, want)
	}
}

func TestCompress(t *testing.T) {
	want := []int{1, 4, 5}
	got := ToSlice(Compress(Literal(1, 2, 3, 4, 5, 6), Literal(true, false, false, true, true)))
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got = %v, want %v", got, want)
	}
}

func TestFlatten(t *testing.T) {
	want := []int{1, 2, 3, 4}
	got := ToSlice(Flatten(Literal(Literal(1, 2), Literal(3, 4))))
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got = %v, want %v", got, want)
	}
}

func TestTee(t *testing.T) {
	want := []int{1, 2, 3, 4}
	tees := Tee(FromSlice(want), 2)
	got0 := ToSlice(tees[0])
	got1 := ToSlice(tees[1])
	if !reflect.DeepEqual(got0, want) {
		t.Errorf("got0 = %v, want %v", got0, want)
	}
	if !reflect.DeepEqual(got1, want) {
		t.Errorf("got1 = %v, want %v", got1, want)
	}
}

func TestMap(t *testing.T) {
	want := []int{2, 4, 6, 8}
	got := ToSlice(Map(func(x int) int { return x * 2 }, Literal(1, 2, 3, 4)))
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got = %v, want %v", got, want)
	}
}

func TestCount(t *testing.T) {
	want := []int{4, 5, 6}
	got := ToSlice(Take(3, Count(4)))
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got = %v, want %v", got, want)
	}
}

func TestTabulate(t *testing.T) {
	want := []int{2, 3, 4, 5}
	got := ToSlice(Take(4, Tabulate(func(i int) int {
		return i + 1
	}, 1)))
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got = %v, want %v", got, want)
	}
}

func TestTail(t *testing.T) {
	want := []int{4, 5}
	got := ToSlice(Tail(2, Literal(0, 1, 2, 3, 4, 5)))
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got = %v, want %v", got, want)
	}
}

func TestGroupByKey(t *testing.T) {
	want := []int{1, 2, 3}
	got := ToSlice(
		Map(
			func(g Group[int, int]) int { return len(g.slice) },
			GroupByKey(
				Literal(1, 2, 2, 3, 3, 3),
				MakeKey(
					func(i int) int { return i },
					func(a, b int) bool { return a == b },
				),
			),
		),
	)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got = %v, want %v", got, want)
	}
}

func TestGroupByValue(t *testing.T) {
	want := []int{1, 2, 3}
	got := ToSlice(
		Map(
			func(g Group[int, int]) int { return len(g.slice) },
			GroupByValue(
				Literal(1, 2, 2, 3, 3, 3),
			),
		),
	)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got = %v, want %v", got, want)
	}
}

func TestAllEqualValue(t *testing.T) {
	want := true
	got := AllEqualValue(Literal(1, 1, 1, 1))
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got = %v, want %v", got, want)
	}
}

func TestAllEqualByKey(t *testing.T) {
	want := true
	got := AllEqualByKey(
		Literal("a", "b", "c", "d"),
		MakeKey(func(s string) int { return len(s) }, IsSameValue[int]),
	)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got = %v, want %v", got, want)
	}
}

func TestPairwise(t *testing.T) {
	want := [][]int{{1, 2}, {2, 3}, {3, 4}}
	got := ToSlice(Map(PairToSlice[int], Pairwise(Literal(1, 2, 3, 4))))
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got = %v, want %v", got, want)
	}
}

func TestProduct(t *testing.T) {
	want := [][]int{{0, 0}, {0, 1}, {1, 0}, {1, 1}}
	got := ToSlice(Product(Literal(0, 1), Literal(0, 1)))
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got = %v, want %v", got, want)
	}
}

func TestTakeWhile(t *testing.T) {
	want := []int{0, 1, 2, 3, 4}
	got := ToSlice(TakeWhile(LessThan(5), IntRange(10)))
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got = %v, want %v", got, want)
	}
}

func TestDropWhile(t *testing.T) {
	want := []int{5, 6, 7, 8, 9}
	got := ToSlice(DropWhile(Not(GreaterThan(4)), IntRange(10)))
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got = %v, want %v", got, want)
	}
}

func TestChunked(t *testing.T) {
	want := [][]int{{1, 2, 3}, {4, 5, 6}, {7}}
	got := ToSlice(Chunked(Literal(1, 2, 3, 4, 5, 6, 7), 3))
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got = %v, want %v", got, want)
	}
}

func TestChunkBy(t *testing.T) {
	want := [][]int{{1, 1, 1}, {2, 2}, {3}}
	got := ToSlice(ChunkBy(Literal(1, 1, 1, 2, 2, 3), Identity[int]))
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got = %v, want %v", got, want)
	}
}

func TestFromAdvanceSafe(t *testing.T) {
	hasNext := true
	want := []int{1}
	safe := FromAdvanceSafe(func() (next bool, value int) {
		next = hasNext
		hasNext = !hasNext
		value = 1
		return
	})
	got := ToSlice(Take(10, Chain(safe, safe, safe)))
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got = %v, want %v", got, want)
	}
}

func TestAllButLast(t *testing.T) {
	want := [][]int{{1, 2, 3}, {4, 5, 6}}
	first, last := AllButLast(Literal(1, 2, 3, 4, 5, 6), 3)
	firstSlice := ToSlice(first)
	got := [][]int{firstSlice, ToSlice(last)}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got = %v, want %v", got, want)
	}
}

func TestIntersperse(t *testing.T) {
	want := []int{1, 5, 2, 5, 3}
	got := ToSlice(Intersperse(Literal(1, 2, 3), 5))
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got = %v, want %v", got, want)
	}
}

func TestISliceEx(t *testing.T) {
	type args struct {
		iter    Iterator[int]
		options []RangeOption
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			"islice(range(10), 2, 6)",
			args{Range(Stop(10)), []RangeOption{Start(2), Stop(6)}},
			[]int{2, 3, 4, 5},
		},
		{
			"islice(range(10), step=2, stop=5)",
			args{Range(Stop(10)), []RangeOption{Step(2), Stop(5)}},
			[]int{0, 2, 4},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToSlice(ISlice(tt.args.iter, tt.args.options...)); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ISlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReduce(t *testing.T) {
	type args struct {
		iter  Iterator[int]
		binOp func(a, b int) int
	}
	tests := []struct {
		name       string
		args       args
		wantValue  int
		wantExists bool
	}{
		{
			"Sum(1...10)",
			args{
				IntRange(10 + 1),
				func(a, b int) int { return a + b },
			},
			55,
			true,
		},
		{
			"Mul(1...5)",
			args{
				Range(Start(1), Stop(6)),
				func(a, b int) int { return a * b },
			},
			120,
			true,
		},
		{
			"Sum()",
			args{
				EmptyIterator[int](),
				func(a, b int) int { return a + b },
			},
			120,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotValue, gotExists := Reduce(tt.args.iter, tt.args.binOp)
			if gotExists != tt.wantExists {
				t.Errorf("Reduce() gotExists = %v, want %v", gotExists, tt.wantExists)
			}
			if gotExists && !reflect.DeepEqual(gotValue, tt.wantValue) {
				t.Errorf("Reduce() gotValue = %v, want %v", gotValue, tt.wantValue)
			}
		})
	}

	t.Run("Join strings", func(t *testing.T) {
		want := "abcdefg"
		got, exists := Reduce(Literal("a", "b", "c", "d", "e", "f", "g"), func(a, b string) string {
			return a + b
		})
		if !exists {
			t.Errorf("Failed calculating result.")
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got = %v, want %v", got, want)
		}
	})
}

func TestNth(t *testing.T) {
	type args struct {
		iter Iterator[int]
		n    int
	}
	tests := []struct {
		name       string
		args       args
		wantValue  int
		wantExists bool
	}{
		{
			"Get negative element",
			args{Literal(1, 2, 3), -2},
			0,
			false,
		},
		{
			"Get 1st element",
			args{Literal(1, 2, 3), 0},
			1,
			true,
		},
		{
			"Get element past end",
			args{Literal(1, 2, 3), 5},
			1,
			false,
		},
		{
			"Get last element",
			args{Literal(1, 2, 3), 2},
			3,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotValue, gotExists := Nth(tt.args.iter, tt.args.n)
			if gotExists != tt.wantExists {
				t.Errorf("Nth() gotExists = %v, want %v", gotExists, tt.wantExists)
			}

			if gotExists && !reflect.DeepEqual(gotValue, tt.wantValue) {
				t.Errorf("Nth() gotValue = %v, want %v", gotValue, tt.wantValue)
			}
		})
	}
}

func TestWindowedWithFiller(t *testing.T) {
	type args struct {
		iter   Iterator[int]
		n      int
		filler int
	}
	tests := []struct {
		name string
		args args
		want [][]int
	}{
		{
			"Pairs",
			args{IntRange(5),
				2,
				0},
			[][]int{{0, 1}, {1, 2}, {2, 3}, {3, 4}},
		},
		{
			"Triplets",
			args{IntRange(5),
				3,
				0},
			[][]int{{0, 1, 2}, {1, 2, 3}, {2, 3, 4}},
		},
		{
			"With Filler",
			args{IntRange(2),
				3,
				0},
			[][]int{{0, 1, 0}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToSlice(WindowedWithFiller(tt.args.iter, tt.args.n, tt.args.filler)); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WindowedWithFiller() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWindowed(t *testing.T) {
	type args struct {
		iter Iterator[int]
		n    int
	}
	tests := []struct {
		name string
		args args
		want [][]int
	}{
		{
			"Pairs",
			args{IntRange(5),
				2},
			[][]int{{0, 1}, {1, 2}, {2, 3}, {3, 4}},
		},
		{
			"Triplets",
			args{IntRange(5),
				3},
			[][]int{{0, 1, 2}, {1, 2, 3}, {2, 3, 4}},
		},
		{
			"Too long",
			args{IntRange(2),
				3},
			[][]int{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToSlice(Windowed(tt.args.iter, tt.args.n)); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WindowedWithFiller() = %v, want %v", got, tt.want)
			}
		})
	}
}
