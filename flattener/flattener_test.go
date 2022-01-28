package main

import (
	"github.com/tmr232/go-explore/itertools"
	"reflect"
	"testing"
)

/*
To Run:

	go build .\flattener\generate\ ; go generate .\flattener\ ; goimports -w .\flattener\generators_gen.go ; go test .\flattener\
*/

func TestMyGen(t *testing.T) {
	want := []int{1, 2, 3}
	got := itertools.ToSlice(MyGen())
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got = %v, want %v", got, want)
	}
}

func TestIfStmt(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		want := []int{1, 3}
		got := itertools.ToSlice(IfStmt(true))
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got = %v, want %v", got, want)
		}
	})
	t.Run("false", func(t *testing.T) {
		want := []int{2, 3}
		got := itertools.ToSlice(IfStmt(false))
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got = %v, want %v", got, want)
		}
	})
}

func TestAnotherIfStmt(t *testing.T) {
	want := []int{0, 1, 2, 3, 4}
	got := itertools.ToSlice(AnotherIfStmt(true))
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got = %v, want %v", got, want)
	}
}

func TestRepeatOne(t *testing.T) {
	want := []int{1, 1, 1, 1, 1}
	got := itertools.ToSlice(itertools.Take(5, RepeatOne()))
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got = %v, want %v", got, want)
	}
}

//func generate_BasicVar() int {
//	var a int
//	b := 0
//	a = 1
//	b = 2
//	c := b
//	return 0
//}
//
//func TestBasicVar(t *testing.T) {
//	want := []int{1}
//	got := itertools.ToSlice(basicVar())
//	if !reflect.DeepEqual(got, want) {
//		t.Errorf("got = %v, want %v", got, want)
//	}
//}

func TestFib(t *testing.T) {
	want := []int{1, 1, 2, 3, 5, 8, 13, 21, 34, 55}
	got := itertools.ToSlice(itertools.Take(10, Fib()))
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got = %v, want %v", got, want)
	}
}

func TestCount(t *testing.T) {
	want := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	got := itertools.ToSlice(itertools.Take(10, Count()))
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got = %v, want %v", got, want)
	}
}
