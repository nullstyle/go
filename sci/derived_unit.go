package sci

// PopulateNormalizedUnit implements Unit
func (u *DerivedUnit) PopulateNormalizedUnit(nu *NormalizedUnit, inv bool) {
	u.Value.U.PopulateNormalizedUnit(nu, inv)
}
