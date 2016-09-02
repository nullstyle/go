package sci

func testSystem() *System {
	var ret System

	ret.MustAdd(&BaseUnit{
		Name:    "meter",
		Measure: Length,
	})

	ret.MustAdd(&BaseUnit{
		Name:    "sec",
		Measure: Time,
	})

	return &ret
}
