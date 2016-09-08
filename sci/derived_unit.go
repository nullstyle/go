package sci

// PopulateNormalizedUnit implements Unit
func (u *DerivedUnit) PopulateNormalizedUnit(nu *NormalizedUnit, inv bool) {
	u.Value.U.PopulateNormalizedUnit(nu, inv)
}

// System implements Unit
func (u *DerivedUnit) System() *System {
	return u.Value.U.System()
}
