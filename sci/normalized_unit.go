package sci

// Add contributes the provided base unit at the provided exponent (negative
// numbers representing to the negative power) to the normalized unit.
func (nu *NormalizedUnit) Add(bu *BaseUnit, exp int) {
	nu.mutex.Lock()
	if nu.Components == nil {
		nu.Components = map[*BaseUnit]int{}
	}
	nu.Components[bu] = nu.Components[bu] + exp
	nu.mutex.Unlock()
}

// Invert returns a copy of the normalized unit, inverted.
func (nu *NormalizedUnit) Invert() *NormalizedUnit {
	var result NormalizedUnit
	for k, v := range nu.Components {
		result.Add(k, -v)
	}
	return &result
}
