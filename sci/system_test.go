package sci

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSystem(t *testing.T) {
	sys := NewSystem("test")
	assert.Equal(t, "test", sys.Name)
	assert.NotNil(t, sys.BaseUnits)
}

func TestSystem_DefineBaseUnit(t *testing.T) {
	sys := NewSystem("test")
	var (
		meter *BaseUnit
		err   error
	)

	t.Run("happy path", func(t *testing.T) {
		meter, err = sys.DefineBaseUnit("meter", Length)
		if assert.NoError(t, err) {
			assert.Equal(t, meter, sys.BaseUnits[Length])
		}
	})

	t.Run("redefinition failse", func(t *testing.T) {
		_, err = sys.DefineBaseUnit("foot", Length)
		assert.Error(t, err)
	})

	t.Run("plural definition panics", func(t *testing.T) {
		assert.Panics(t, func() {
			sys.DefineBaseUnit("seconds", Length)
		})
	})
}

func TestSystem_DefineUnit(t *testing.T) {
	sys := NewSystem("test")

	var (
		inch Unit
		foot Unit
		err  error
	)

	inch, err = sys.DefineBaseUnit("inch", Length)
	require.NoError(t, err)

	t.Run("happy path", func(t *testing.T) {
		foot, err = sys.DefineUnit("foot", "12 inches")
		if assert.NoError(t, err) {
			assert.Equal(t, foot, sys.units["foot"])

			du := foot.(*DerivedUnit)
			assert.Equal(t, inch, du.Value.U)
			assert.Equal(t, "12", du.Value.M)
		}
	})
}

func TestSystem_Parse_Magnitudes(t *testing.T) {
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
		assert.Equal(t, sys.Nil(), u)
	}

	u, err = sys.ParseUnit("meter")
	if assert.NoError(t, err) {
		assert.Equal(t, sys.BaseUnits[Length], u)
	}

	u, err = sys.ParseUnit("/sec")
	if assert.NoError(t, err) {
		assert.EqualValues(t, &DivUnit{N: sys.Nil(), D: sys.BaseUnits[Time]}, u)
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
		assert.EqualValues(t, &MulUnit{
			sys.BaseUnits[Length],
			sys.BaseUnits[Time],
		}, u)
	}

}
