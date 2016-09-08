package sci

// PopulateNormalizedUnit implements Unit
func (u *DivUnit) PopulateNormalizedUnit(nu *NormalizedUnit, inv bool) {
	ninv := inv
	dinv := !inv

	u.N.PopulateNormalizedUnit(nu, ninv)
	u.D.PopulateNormalizedUnit(nu, dinv)
}

// System implements Unit
func (u *DivUnit) System() *System {
	return u.N.System()
}
