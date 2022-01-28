package exitstack

import (
	"fmt"
	"io"
	"log"
)

type ExitStack struct {
	stack       []io.Closer
	isCancelled bool
}

func (es *ExitStack) Push(closer io.Closer) {
	es.stack = append(es.stack, closer)
}

func (es *ExitStack) Close() {
	if es.isCancelled {
		return
	}

	for i := len(es.stack) - 1; i >= 0; i-- {
		_ = es.stack[i].Close()
	}
}

func (es *ExitStack) Cancel() {
	es.isCancelled = true
}

type Pusher[T io.Closer] struct {
	t   T
	err error
}

func (p Pusher[T]) Into(es *ExitStack) (T, error) {
	if p.err == nil {
		es.Push(p.t)
	}
	return p.t, p.err
}

func Push[T io.Closer](t T, err error) Pusher[T] {
	return Pusher[T]{t, err}
}

type File struct {
	name string
}

func (f File) Close() error {
	fmt.Println("Closing ", f.name)
	return nil
}

func Open(name string) (File, error) {
	return File{name}, nil
}

func main() {
	es := new(ExitStack)
	defer es.Close()

	x, err := Push(Open("X")).Into(es)
	if err != nil {
		log.Fatal(err)
	}

	y, err := Push(Open("Y")).Into(es)
	if err != nil {
		log.Fatal(err)
	}

	//es.Cancel()

	fmt.Println(x, y)
}
