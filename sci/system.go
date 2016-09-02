package sci

func (sys *System) Add(u Unit) Unit {
	// TODO: make sure the unit can be added... e.g. no two base units may share
	// the same measure in a system, no defined unit can be against a base unit
	// not defined in the system.

	sys.Units = append(sys.Units, u)
}

func (sys *System) MustParse(val string) *Value {
	v, err := sys.Parse(val)
	if err != nil {
		panic(err)
	}
	return v
}

func (sys *System) Parse(val string) (*Value, error) {
	return nil, nil
}
