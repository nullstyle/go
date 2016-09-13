package sci

// Add contributes the provided base unit at the provided exponent (negative
// numbers representing to the negative power) to the normalized unit.
func (nu *NormalizedUnit) Add(m Measure, exp int) {
	nu.mutex.Lock()
	nu.Components[m] = nu.Components[m] + exp
	nu.mutex.Unlock()
}
