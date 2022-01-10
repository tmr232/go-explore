package main

import (
	"fmt"
	"github.com/tmr232/go-explore/itertools"
)

func Fibonacci() itertools.Iterator[int] {
	a, b := 1, 1
	advance := func() (bool, int) {
		retval := a
		a, b = b, a+b
		return true, retval
	}
	return itertools.FromAdvance(advance)
}

func main() {
	itertools.ForEach(
		itertools.Take( // Take the first
			10,          // 10 elements
			Fibonacci(), // from our Fibonacci sequence
		),
		func(v int) { fmt.Println(v) }, // and print them
	)
}
