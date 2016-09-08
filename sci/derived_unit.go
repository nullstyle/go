package sci

// PopulateNormalizedUnit implements Unit
func (u *DerivedUnit) PopulateNormalizedUnit(nu *NormalizedUnit, inv bool) {
	u.Value.U.PopulateNormalizedUnit(nu, inv)
}

// String implements fmt.Stringer
func (u *DerivedUnit) String() string {
	return u.System().getUnitName(u)
}

// System implements Unit
func (u *DerivedUnit) System() *System {
	return u.Value.U.System()
}
