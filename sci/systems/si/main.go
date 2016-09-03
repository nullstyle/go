package si

import "github.com/nullstyle/go/sci"

// System is the set of units that make up the SI unit system.
// https://en.wikipedia.org/wiki/International_System_of_Units
var System = sci.NewSystem("SI")

var (
	// Meter represents the SI base unit of length, the meter.
	// https://en.wikipedia.org/wiki/SI_base_unit
	Meter = System.MustDefineBaseUnit(
		"meter",
		sci.Length,
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

var _ sci.Unit = Meter
