package sci

// PopulateNormalizedUnit implements Unit
func (u *NilUnit) PopulateNormalizedUnit(nu *NormalizedUnit, inv bool) {
	return
}

// String implements fmt.Stringer
func (u *NilUnit) String() string {
	return "nil unit"
}

// System implements Unit
func (u *NilUnit) System() *System {
	return u.system
}
