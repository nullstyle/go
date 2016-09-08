package sci

// PopulateNormalizedUnit implements Unit
func (bu *BaseUnit) PopulateNormalizedUnit(nu *NormalizedUnit, inverted bool) {
	if inverted {
		nu.Add(bu, -1)
	} else {
		nu.Add(bu, 1)
	}
}

// String implements fmt.Stringer
func (bu *BaseUnit) String() string {
	return bu.Name
}

// System implements Unit
func (bu *BaseUnit) System() *System {
	return bu.system
}
