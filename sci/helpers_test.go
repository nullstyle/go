package sci

func testSystem() *System {
	ret := NewSystem("test")

	ret.MustDefineBaseUnit("meter", Length)
	ret.MustDefineBaseUnit("sec", Time)

	return ret
}
