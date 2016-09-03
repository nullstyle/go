package sci

// PopulateNormalizedUnit implements Unit
func (bu *BaseUnit) PopulateNormalizedUnit(nu *NormalizedUnit, inverted bool) {
	if inverted {
		nu.Add(bu, -1)
	} else {
		nu.Add(bu, 1)
	}
}
