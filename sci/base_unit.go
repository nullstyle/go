package sci

func (u *BaseUnit) Compat(other Unit) bool {
	switch other := other.(type) {
	case *BaseUnit:
		return u == other
	case *DefinedUnit:
		return u == other.Base
	case *NilUnit:
    return true
	default:
		return false
	}
}
