// Package uscs provides the United states customary unit system.
// https://en.wikipedia.org/wiki/United_States_customary_units
package uscs

import "github.com/nullstyle/go/sci"

// System is the container for the units of this package
var System = sci.NewSystem("USCS")

var (
	Inch = System.MustDefineBaseUnit(
		"inch",
		sci.Length,
	)

	Foot = System.MustDefineUnit(
		"foot",
		"12 inch",
	)

	Yard = System.MustDefineUnit(
		"yard",
		"3 foot",
	)

	Mile = System.MustDefineUnit(
		"mile",
		"5280 foot",
	)
)

// MustParse is the panicking version of Parse
func MustParse(val string) *sci.Value {
	return System.MustParse(val)
}

// Parse parses a value expressed with SI units
func Parse(val string) (*sci.Value, error) {
	return System.Parse(val)
}
