package examples

import (
	"fmt"
	"log"

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

func ExampleAdd_HighLevel() {
	x := si.MustParse("10 meters")
	y := si.MustParse("2 meters")
	z := &sci.Value{}

	z.Add(x, y)
	fmt.Print(z.M)
	// output: 110
}

func ExampleAdd_Expression() {
	result, err := si.Evaluate("10 meters + 2 meters")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(result)
	// output: 12 meters
}
