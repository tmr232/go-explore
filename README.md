# go-explore

In this repo I experiment with Golang.

While the packages within may be in a usable state,
they are still highly experimental and subject to change.

## Itertools

Porting Python's `itertools` module (and a bit more) into Go.

It is built around the `Iterator[T]` interface: 

```go
type Iterator[T any] interface {
	// Next tries to advance to the next value.
	// Returns true if a value exists, false if not.
	// Once an iterator returns false to indicate exhaustion,
	// it should continue returning false.
	Next() bool
	// Value returns the current value of the iterator.
	// Next() must be called and return true before every call to Value().
	Value() T
}
```

Iteration is fairly straightforward:

```go
for iter := Literal(1,2,3); iter.Next(); {
	fmt.Println(iter.Value())
}
```

Would result in

```text
1
2
3
```

#### Fibonacci

A simple Fibonacci iterator can be written as follows (`itertools.` remove for brevity):

```go
func Fibonacci() Iterator[int] {
	a, b := 1, 1
	advance := func() (bool, int) {
        retval := a
		a, b = b, a + b
		return true, retval
	}
	return FromAdvance(advance)
}
```

Note that this iterator is "infinite".
It does not stop after a set number, but keeps going.
If we only want to print the first 10 elements, 
we can `Take` the first 10 elements and print them:

```go
ForEach(
    Take(                           // Take the first
        10,                         // 10 elements
        Fibonacci(),                // from our Fibonacci sequence
    ),
    func(v int) { fmt.Println(v) }, // and print them
)
```

Resulting in 
```text
1
1 
2 
3 
5 
8 
13
21
34
55
```