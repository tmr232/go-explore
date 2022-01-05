package main

import (
	"fmt"
	"testing"
)

//go:generate go run ./generate

func TestMyGen(t *testing.T) {
	for gen := MyGen(); gen.Next(); {
		fmt.Println(gen.Value())
	}
}
