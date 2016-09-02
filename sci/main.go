// Package sci implements a system for performing calculations on physical
// quantities.  The math API is inspired by the big package of the go standard
// library: The receiver of a .
//
// The core
package sci

import "github.com/pkg/errors"

const (
	Length Measure = "Length"
	Time   Measure = "Time"
)

var (
	ErrIncompatibleTypes = errors.New("incompatible types")
)

// Unit represents any unit of measure
type Unit interface {
	Compat(Unit) bool
}

// BaseUnit represents the a base unit of a given measure against which other
// units can be defined.  For example, a meter is a base unit of length.
type BaseUnit struct {
	Name    string
	Measure Measure
}

// DefinedUnit represents a unit expressed in relation to some base unit.
type DefinedUnit struct {
	Scalar string
	Base   *BaseUnit
}

type DivUnit struct {
	N Unit
	M Unit
}

type MulUnit []Unit

// NilUnit represents "no unit".  Values of it represents "a mignitude without a
// unit".  NilUnit is also used to represent inverse units, such as hz (1 / s)
// when combind using DivUnit.
type NilUnit struct {
}

// Measure represents a domain of measurement, such as length, time, or mass.
type Measure string

// Prefix represents a unit prefix.
type Prefix struct {
	Name    string
	Scalar  string
	Aliases []string
}

// System represents a system of measurement.
type System struct {
	// Name is an optional name for a system
	Name string

	// Units represents the collection of defined units
	Units []Unit

	// ByName is an index of the units in the system by name.
	ByName map[string]int
}

// Value represents a value. Examples include "3 mm" or "10 m/s"
type Value struct {
	// M is the magnitude or the multitude of the value (depending upon whether or
	// not the unit is collective or non-collective), expressed as a string
	// (parsing rules TBD)
	M string

	// U is the unit of the value.
	U Unit
}

var _ Unit = &BaseUnit{}
var _ Unit = &DefinedUnit{}
var _ Unit = &NilUnit{}
var _ Unit = &DivUnit{}
var _ Unit = &MulUnit{}
