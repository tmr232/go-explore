package main

import (
	"github.com/tmr232/go-explore/itertools"
	"reflect"
	"testing"
)

//go:generate go run ./generate

func generate_MyGen() int {
	return 1
	return 2
	return 3
}

func TestMyGen(t *testing.T) {
	want := []int{1, 2, 3}
	got := itertools.ToSlice(MyGen())
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got = %v, want %v", got, want)
	}
}
