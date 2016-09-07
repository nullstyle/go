package unit

// Abs returns the absolute value of the exponent
func (e *Exp) Abs() int64 {
	if e.Exp < 0 {
		return -e.Exp
	}

	return e.Exp
}
