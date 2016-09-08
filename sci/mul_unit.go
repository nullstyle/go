package sci

// PopulateNormalizedUnit implements Unit
func (u *MulUnit) PopulateNormalizedUnit(nu *NormalizedUnit, inv bool) {
	for _, mu := range *u {
		mu.PopulateNormalizedUnit(nu, inv)
	}
}

// String implements fmt.Stringer
func (u *MulUnit) String() string {
	return u.System().getUnitName(u)
}

// System implements Unit
func (u *MulUnit) System() *System {
	return (*u)[0].System()
}
