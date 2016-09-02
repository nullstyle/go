package sci

import "strings"

func (sys *System) Add(u Unit) error {
	// TODO: make sure the unit can be added... e.g. no two base units may share
	// the same measure in a system, no defined unit can be against a base unit
	// not defined in the system.

	sys.Units = append(sys.Units, u)
	return nil
}

func (sys *System) MustAdd(u Unit) {
	err := sys.Add(u)
	if err != nil {
		panic(err)
	}
}

func (sys *System) MustParse(val string) *Value {
	v, err := sys.Parse(val)
	if err != nil {
		panic(err)
	}
	return v
}

// Parse parses a single value using the units defined (i.e. previously added to
// the system using Add()) in sys.
func (sys *System) Parse(val string) (*Value, error) {
	stripped := strings.TrimSpace(val)

	if stripped == "" {
		return nil, ErrBlankValue
	}

	matches := MagnitudeRegexp.FindStringSubmatchIndex(val)
	if matches == nil {
		return nil, &ParseError{Input: val, FailurePhase: "extract magnitude"}
	}

	// the first submatch of the match is the whole magnitude (in indexes 2 and 3
	// according to package regexp rules)
	magstart, magend := matches[2], matches[3]

	magnitude := val[magstart:magend]
	unitstr := ""

	// if the magnitude match does not consume the whole input, everything past
	// the match is to be considered the unit of the value.
	if magend < len(val) {
		unitstr = val[magend+1:]
	}

	_ = unitstr
	// unit, err := sys.ParseUnit(unitstr)
	// if err != nil {
	// 	return nil, errors.Wrap(err, "parse unit")
	// }

	return &Value{M: magnitude, U: nil}, nil
}

// ParseUnit converts the provided string into a Unit value, looking up defined
// units as necessary and forming new algebraic units as specified.
func (sys *System) ParseUnit(unitstr string) (Unit, error) {
	//TODO
	return nil, nil
}
