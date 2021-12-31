package tuple

import (
	"fmt"
	"testing"
)

func TestList(t *testing.T) {
	p := MakeList6(1, "b", struct{ int }{3}, "four", "4+1", 6)
	fmt.Println(Pretty(p))
	fmt.Println(Get0(p))
	fmt.Println(Get1(p))
	fmt.Println(Get2(p))
	fmt.Println(Get3(p))
	fmt.Println(Get4(p))
	fmt.Println(Get5(p))
}
