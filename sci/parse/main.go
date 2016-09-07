package parse

import (
	"github.com/nullstyle/go/sci/parse/unit"
	"github.com/nullstyle/go/sci/parse/value"
)

// Unit parses the provided input into the abstract syntax tree for a unit.
func Unit(input string) (unit.U, error) {
	return unit.Parse(input)
}

// Value parses the provided input into the (magnitude, unit) tuple.
func Value(input string) (value.V, error) {
	return value.Parse(input)
}
