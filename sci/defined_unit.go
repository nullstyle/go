package sci

// Compat implements the `Unit` interface
func (u *DefinedUnit) Compat(other Unit) bool {
	switch other := other.(type) {
	case *BaseUnit:
		return u.Base == other
	case *DefinedUnit:
		return u.Base == other.Base
	default:
		return false
	}
}
