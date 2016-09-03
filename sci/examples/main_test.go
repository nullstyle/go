package examples

import (
	"fmt"

	"github.com/nullstyle/go/sci"
	"github.com/nullstyle/go/sci/systems/si"
)

func ExampleAdd_LowLevel() {
	// same unit
	x := &sci.Value{M: "10", U: si.Meter}
	y := &sci.Value{M: "100", U: si.Meter}
	z := &sci.Value{}

	z.Add(x, y)
	fmt.Print(z.M)
	// output: 110
}
