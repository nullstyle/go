package sci

import (
	"fmt"
	"math/big"
)

// Interface conformity confirmations
var _ Unit = &BaseUnit{}
var _ Unit = &DerivedUnit{}
var _ Unit = &NilUnit{}
var _ Unit = &DivUnit{}
var _ Unit = &MulUnit{}

var _ error = &MagnitudeError{}
var _ error = &ExpToBigError{}
var _ error = &BaseUnitAlreadyDefinedError{}

var _ fmt.Stringer = &Value{}

// converter converts from one magnitude to another using a conversion factor.
type converter struct {
	factor big.Float
}
