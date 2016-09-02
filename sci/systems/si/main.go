package si

import "github.com/nullstyle/go/sci"

var System = &sci.System{}

var (
	Meter = &sci.BaseUnit{
		Measure: "length",
	}
)

func MustParse(val string) *sci.Value {
	return System.MustParse(val)
}

func Parse(val string) (*sci.Value, error) {
	return System.Parse(val)
}

var _ sci.Unit = Meter
