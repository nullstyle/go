package sci

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnitParser(t *testing.T) {
	sys := testSystem()
	cases := []struct {
		Name     string
		Input    string
		Expected Unit
	}{
		{
			"simple",
			"meter",
			sys.BaseUnits[Length],
		},
		{
			"simple-2",
			"sec",
			sys.BaseUnits[Time],
		},
		{
			"div",
			"meter / sec",
			&DivUnit{N: sys.BaseUnits[Length], D: sys.BaseUnits[Time]},
		},
		{
			"mul",
			"meter * sec",
			&MulUnit{sys.BaseUnits[Length], sys.BaseUnits[Time]},
		},
		{
			"mul-div",
			"meter * sec / sec",
			&DivUnit{
				N: &MulUnit{sys.BaseUnits[Length], sys.BaseUnits[Time]},
				D: sys.BaseUnits[Time],
			},
		},
		{
			"mul-div",
			"meter / sec * sec",
			&MulUnit{
				&DivUnit{N: sys.BaseUnits[Length], D: sys.BaseUnits[Time]},
				sys.BaseUnits[Time],
			},
		},
		{
			"parens",
			"meter * (sec / sec)",
			&MulUnit{
				sys.BaseUnits[Length],
				&DivUnit{N: sys.BaseUnits[Time], D: sys.BaseUnits[Time]},
			},
		},
		{
			"parens-2",
			"meter / (sec * sec)",
			&DivUnit{
				N: sys.BaseUnits[Length],
				D: &MulUnit{sys.BaseUnits[Time], sys.BaseUnits[Time]},
			},
		},
	}

	for _, kase := range cases {
		t.Run(kase.Name, func(t *testing.T) {
			u := UnitParser{
				Buffer: kase.Input,
				System: sys,
				Result: Nil,
			}

			u.Init()
			err := u.Parse()
			if !assert.NoError(t, err) {
				return
			}

			u.Execute()
			if !assert.NoError(t, u.Err) {
				return
			}

			assert.EqualValues(t, kase.Expected, u.Result)
		})
	}
}
