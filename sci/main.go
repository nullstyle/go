// Package sci implements a system for performing calculations on physical
// quantities.  The math API is inspired by the big package of the go standard
// library: The receiver of a .
//
// The core
package sci

import (
	"regexp"
	"sync"

	"github.com/pkg/errors"
)

//go:generate peg -switch -inline unit_parser.peg

const (
	// Length represents the length measure
	Length Measure = "Length"

	// Time represents the time measure
	Time Measure = "Time"
)

const (
	// MaxExp represents the largest absolute value that an exponent is allowed to
	// have within this library.  It's a conservative cap for now, as I (scott) am
	// not very sure what sort of issues will arise if left higher or unbounded.
	MaxExp = 4
)

// Nil is the singleton instance of *NilUnit
var Nil = &NilUnit{}

var (
	MagnitudeRegexp = regexp.MustCompile(
		"^((-)?([1-9][0-9]*)(\\.[0-9]+)?)([^0-9]|$)",
	)
)

var (
	// ErrIncompatibleTypes is returned when attempted to perform an operation
	// (such as addition) on two incompatible types.
	ErrIncompatibleTypes = errors.New("incompatible types")
)

// Unit represents any unit of measure
type Unit interface {
	PopulateNormalizedUnit(nu *NormalizedUnit, inverted bool)
}

// BaseUnit represents the a base unit of a given measure against which other
// units can be defined.  For example, a meter is a base unit of length.
type BaseUnit struct {
	Name    string
	Measure Measure

	system *System
}

// DerivedUnit represents a unit expressed in relation to some base unit.
type DerivedUnit struct {
	Value *Value
}

// DivUnit represents a compound unit such as "feet / hour"
type DivUnit struct {
	N Unit
	D Unit
}

// MulUnit represents a compound unit such as "foot*pound"
type MulUnit []Unit

// NilUnit represents "no unit".  Values of it represents "a mignitude without a
// unit".  NilUnit is also used to represent inverse units, such as hz (1 / s)
// when combind using DivUnit.
type NilUnit struct {
}

// NormalizedUnit represents the non-aliased form of a unit, expressed
// completely in terms of base units.  derived units are expanded into base
// units and then contribute themselves to either the numerator by increasing
// the count by one, or to the denominator by decreasing the count by one.  This
// will be a recursive process.
type NormalizedUnit struct {
	Components map[*BaseUnit]int
	mutex      sync.Mutex
}

// BaseUnitAlreadyDefinedError is an error that occurs when attempting to
// redefine the base unit used for some measure in a system of units.  A system
// of units may only have one base unit per measure to ensure that we can define
// any value belonging to a given measure in relation to a single base unit.
type BaseUnitAlreadyDefinedError struct {
	Existing *BaseUnit
}

// Converter represents a type that can convert a value of one unit into
// another.
type Converter interface {
	Convert(in *Value) (*Value, error)
}

// ExpToBigError is the error that is returned when a string that is being
// parsed into a unit has an exponent that is too large.
type ExpToBigError struct {
	Exp int64
}

// MagnitudeError represents the error produces when trying to operate on a
// value whose magnitude (the M field) is invalid.
type MagnitudeError struct {
	M string
}

// Measure represents a domain of measurement, such as length, time, or mass.
type Measure string

// Prefix represents a unit prefix.
type Prefix struct {
	Name    string
	Scalar  string
	Aliases []string
}

// UnitAlreadyDefinedError is an error that occurs when attempting to redefine
// the a named unit.  A given system may not have multiple units that have the
// same name dfined within it.
type UnitAlreadyDefinedError struct {
	Existing Unit
	Name     string
}

// UnitNotDefinedError is an error that occurs when attempting to lookup a unit
// in the system.
type UnitNotDefinedError struct {
	// Name is the name for the unit attempted to be found
	Name string
}

// System represents a system of measurement.
type System struct {
	// Name is an optional name for a system
	Name string

	// BaseUnits represents the collection of defined base units
	BaseUnits map[Measure]*BaseUnit

	// units represents all of the units defined in this system, base units
	// included.
	units map[string]Unit
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

// NewSystem creates a new unit system with the given name
func NewSystem(name string) *System {
	var ret System
	ret.Name = name
	ret.BaseUnits = make(map[Measure]*BaseUnit)
	ret.units = make(map[string]Unit)

	return &ret
}

// Interface conformity confirmations
var _ Unit = &BaseUnit{}
var _ Unit = &DerivedUnit{}
var _ Unit = &NilUnit{}
var _ Unit = &DivUnit{}
var _ Unit = &MulUnit{}

var _ error = &MagnitudeError{}
var _ error = &ExpToBigError{}
var _ error = &BaseUnitAlreadyDefinedError{}
