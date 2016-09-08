package sci

func testSystem() *System {
	ret := NewSystem("test")

	ret.MustDefineBaseUnit("meter", Length)
	ret.MustDefineBaseUnit("sec", Time)
	ret.MustDefineUnit("foot", "0.3048 meter")

	return ret
}
