package sci

// PopulateNormalizedUnit implements Unit
func (u *NilUnit) PopulateNormalizedUnit(nu *NormalizedUnit, inv bool) {
	return
}

// System implements Unit
func (u *NilUnit) System() *System {
	return u.system
}
