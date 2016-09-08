package sci

// PopulateNormalizedUnit implements Unit
func (u *DivUnit) PopulateNormalizedUnit(nu *NormalizedUnit, inv bool) {
	ninv := inv
	dinv := !inv

	u.N.PopulateNormalizedUnit(nu, ninv)
	u.D.PopulateNormalizedUnit(nu, dinv)
}

// String implements fmt.Stringer
func (u *DivUnit) String() string {
	return u.System().getUnitName(u)
}

// System implements Unit
func (u *DivUnit) System() *System {
	return u.N.System()
}
