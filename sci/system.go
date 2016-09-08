package sci

import (
	"github.com/gedex/inflector"
	"github.com/nullstyle/go/sci/parse"
	"github.com/pkg/errors"
)

// DefineBaseUnit defines a base unit of measure in the unit system. Only one
// base unit per measure is allowed, and every other unit in a given system must
// be defined in terms of the base units of the system.
func (sys *System) DefineBaseUnit(name string, m Measure) (*BaseUnit, error) {

	// NOTE: if you run into this panic and are using a singular name, please open
	// a github issue and include the name you used as it signifies that our
	// inflector is broken.
	if inflector.Singularize(name) != name {
		panic(name + " is not singular; base units must be specified using a singular name")
	}

	existing, ok := sys.BaseUnits[m]
	if ok {
		return nil, &BaseUnitAlreadyDefinedError{
			Existing: existing,
		}
	}
	unit := &BaseUnit{
		Name:    name,
		Measure: m,
		system:  sys,
	}

	err := sys.addUnit(name, unit)
	if err != nil {
		return nil, errors.Wrap(err, "failed to add unit to index")
	}

	sys.BaseUnits[m] = unit

	return unit, nil
}

// DefineUnit adds a new DerivedUnit to the system using val as the definition.
// val must be expressed in terms of units previously defined within the system.
func (sys *System) DefineUnit(name string, valstr string) (Unit, error) {
	val, err := sys.Parse(valstr)
	if err != nil {
		return nil, errors.Wrap(err, "parse value failed")
	}

	unit := &DerivedUnit{
		Value: val,
	}

	err = sys.addUnit(name, unit)
	if err != nil {
		return nil, errors.Wrap(err, "failed to add unit to index")
	}

	return unit, nil
}

// LookupUnit finds a unit by name or alias in the system of units. An error of
// type *UnitNotFoundError will be returned if a unit with the given name has
// not been previously defined.
func (sys *System) LookupUnit(name string) (Unit, error) {
	one := inflector.Singularize(name)

	found, ok := sys.units[one]
	if !ok {
		return nil, &UnitNotDefinedError{Name: name}
	}

	return found, nil
}

// MustDefineBaseUnit is the panicking version of define base unit
func (sys *System) MustDefineBaseUnit(name string, m Measure) *BaseUnit {
	u, err := sys.DefineBaseUnit(name, m)
	if err != nil {
		panic(err)
	}
	return u
}

// MustDefineUnit is the panicking version of DefineUnit
func (sys *System) MustDefineUnit(name string, val string) Unit {
	u, err := sys.DefineUnit(name, val)
	if err != nil {
		panic(err)
	}
	return u
}

// MustParse is the panicking version of Parse
func (sys *System) MustParse(val string) *Value {
	v, err := sys.Parse(val)
	if err != nil {
		panic(err)
	}
	return v
}

func (sys *System) Nil() *NilUnit {
	return sys.nilu
}

// Parse parses a single value using the units defined (i.e. previously added to
// the system using Add()) in sys.
func (sys *System) Parse(val string) (*Value, error) {
	v, err := parse.Value(val)
	if err != nil {
		return nil, errors.Wrap(err, "parse value")
	}

	u, err := sys.ParseUnit(v.U)
	if err != nil {
		return nil, errors.Wrap(err, "parse unit")
	}

	return &Value{M: v.M, U: u}, nil
}

// ParseUnit converts the provided string into a Unit value, looking up defined
// units as necessary and forming new algebraic units as specified.
func (sys *System) ParseUnit(unitstr string) (Unit, error) {
	ast, err := parse.Unit(unitstr)
	if err != nil {
		return nil, errors.Wrap(err, "parse failed")
	}

	eu := &evalUnit{System: sys}
	u, err := eu.eval(ast)
	if err != nil {
		return nil, errors.Wrap(err, "eval failed")
	}

	return u, nil
}

// addUnit records the provided unit at the given name, providing it wouldn't
// redefine an existing name.
func (sys *System) addUnit(name string, unit Unit) error {
	one := inflector.Singularize(name)

	found, ok := sys.units[one]

	if ok {
		return &UnitAlreadyDefinedError{Existing: found, Name: name}
	}

	sys.units[one] = unit

	return nil
}
