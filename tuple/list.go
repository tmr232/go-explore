package tuple

import "fmt"

type End struct{}

type List[First any, Second any] struct {
	first  First
	second Second
}

func (l List[F, S]) First() F {
	return l.first
}
func (l List[F, S]) Second() S {
	return l.second
}

type pretty[F, S any] List[F, S]

func Pretty[F, S any](l List[F, S]) fmt.Stringer {
	return pretty[F, S](l)
}

func (l List[F, S]) String() string {
	_, ok := any(l.Second()).(End)
	if ok {
		return fmt.Sprint(l.First())
	} else {
		return fmt.Sprintf("%v, %v", l.First(), l.Second())
	}
}

func (p pretty[F, S]) String() string {
	return fmt.Sprintf("<< %v >>", List[F, S](p))
}

func MakeList[First any, Second any](first First, second Second) List[First, Second] {
	return List[First, Second]{first, second}
}

func Get0[A, B any](l List[A, B]) A {
	return l.First()
}

func Get1[A, B, C any](l List[A, List[B, C]]) B {
	return Get0(l.Second())
}

func Get2[A, B, C, D any](l List[A, List[B, List[C, D]]]) C {
	return Get1(l.Second())
}

func Get3[A, B, C, D, E any](l List[A, List[B, List[C, List[D, E]]]]) D {
	return Get2(l.Second())
}

func Get4[A, B, C, D, E, F any](l List[A, List[B, List[C, List[D, List[E, F]]]]]) E {
	return Get3(l.Second())
}

func Get5[A, B, C, D, E, F, G any](l List[A, List[B, List[C, List[D, List[E, List[F, G]]]]]]) F {
	return Get4(l.Second())
}

func MakeList1[A any](a A) List[A, End] {
	return MakeList(a, End{})
}

func MakeList2[A, B any](a A, b B) List[A, List[B, End]] {
	return MakeList(a, MakeList1(b))
}

func MakeList3[A, B, C any](a A, b B, c C) List[A, List[B, List[C, End]]] {
	return MakeList(a, MakeList2(b, c))
}

func MakeList4[A, B, C, D any](a A, b B, c C, d D) List[A, List[B, List[C, List[D, End]]]] {
	return MakeList(a, MakeList3(b, c, d))
}

func MakeList5[A, B, C, D, E any](a A, b B, c C, d D, e E) List[A, List[B, List[C, List[D, List[E, End]]]]] {
	return MakeList(a, MakeList4(b, c, d, e))
}

func MakeList6[A, B, C, D, E, F any](a A, b B, c C, d D, e E, f F) List[A, List[B, List[C, List[D, List[E, List[F, End]]]]]] {
	return MakeList(a, MakeList5(b, c, d, e, f))
}
