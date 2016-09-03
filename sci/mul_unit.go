package sci

// PopulateNormalizedUnit implements Unit
func (u *MulUnit) PopulateNormalizedUnit(nu *NormalizedUnit, inv bool) {
	for _, mu := range *u {
		mu.PopulateNormalizedUnit(nu, inv)
	}
}
