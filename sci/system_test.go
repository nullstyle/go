package sci

import "testing"
import "github.com/stretchr/testify/assert"

func TestParse_Magnitudes(t *testing.T) {
	sys := testSystem()

	cases := []struct {
		Name      string
		Input     string
		ExpectedM string
		ExpectedU Unit
	}{
		{"integral", "10 meter", "10", sys.BaseUnits[Length]},
		{"integral-negative", "-10 meter", "-10", sys.BaseUnits[Length]},
		{"integral-nospace", "10meter", "10", sys.BaseUnits[Length]},
		{"integral-negative-nospace", "-10meter", "-10", sys.BaseUnits[Length]},
		{"floating", "10.001 sec", "10.001", sys.BaseUnits[Time]},
		{"floating-negative", "-10.001 meter", "-10.001", sys.BaseUnits[Length]},
		{"floating-nospace", "10.001meter", "10.001", sys.BaseUnits[Length]},
		{"floating-negative-nospace", "-10.001meter", "-10.001", sys.BaseUnits[Length]},
	}

	for _, kase := range cases {
		t.Run(kase.Name, func(t *testing.T) {
			v, err := sys.Parse(kase.Input)
			assert.NoError(t, err)
			assert.Equal(t, kase.ExpectedM, v.M)
			assert.EqualValues(t, kase.ExpectedU, v.U)
		})
	}

	// TODO:
	// assert.Equal(t, sys.MustLookup("meter"), v.U)

}

func TestParseUnit(t *testing.T) {
	sys := testSystem()

	u, err := sys.ParseUnit("")
	if assert.NoError(t, err) {
		assert.Equal(t, Nil, u)
	}

	u, err = sys.ParseUnit("meter")
	if assert.NoError(t, err) {
		assert.Equal(t, sys.BaseUnits[Length], u)
	}

	u, err = sys.ParseUnit("/sec")
	if assert.NoError(t, err) {
		assert.EqualValues(t, &DivUnit{N: Nil, D: sys.BaseUnits[Time]}, u)
	}

	u, err = sys.ParseUnit("meter/sec")
	if assert.NoError(t, err) {
		assert.EqualValues(t, &DivUnit{
			N: sys.BaseUnits[Length],
			D: sys.BaseUnits[Time],
		}, u)
	}

	u, err = sys.ParseUnit("meter*sec")
	if assert.NoError(t, err) {
		assert.EqualValues(t, &DivUnit{
			N: sys.BaseUnits[Length],
			D: sys.BaseUnits[Time],
		}, u)
	}

}
