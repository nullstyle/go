package calc

import (
	"testing"

	"github.com/nullstyle/go/influx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCalculator(t *testing.T) {
	cases := []struct {
		Name     string
		Actions  []influx.Action
		Expected string
	}{

		// basic operations and digit exercise
		{
			Name:     "1+2",
			Actions:  []influx.Action{One, Add{}, Two, Equals{}},
			Expected: "3",
		},
		{
			Name:     "45-6",
			Actions:  []influx.Action{Four, Five, Sub{}, Six, Equals{}},
			Expected: "39",
		},
		{
			Name:     "78*901",
			Actions:  []influx.Action{Seven, Eight, Mul{}, Nine, Zero, One, Equals{}},
			Expected: "70278",
		},
		{
			Name: "2345/67",
			Actions: []influx.Action{
				Two, Three, Four, Five,
				Div{},
				Six, Seven, Equals{},
			},
			Expected: "35",
		},

		{
			Name: "multiple operators: 2+2+2",
			Actions: []influx.Action{
				Two, Add{}, Two, Add{}, Two, Equals{},
			},
			Expected: "6",
		},
		{
			Name: "=",
			Actions: []influx.Action{
				Equals{},
			},
			Expected: "",
		},
		{
			Name: "1=",
			Actions: []influx.Action{
				One, Equals{},
			},
			Expected: "1",
		},
	}

	for _, kase := range cases {
		t.Run(kase.Name, func(t *testing.T) {

			// setup
			var c Calculator
			store, err := influx.New(&c)
			require.NoError(t, err)

			// run actions
			for i, a := range kase.Actions {
				err := store.Dispatch(a)
				assert.NoError(t, err, "key entry failed: %d", i)
			}

			// check result
			actual := c.Display()
			assert.Equal(t,
				kase.Expected, actual,
				"bad: %s != %s", kase.Name, actual,
			)
		})
	}
}
