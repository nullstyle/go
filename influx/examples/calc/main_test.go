package calc

import (
	"testing"

	"github.com/nullstyle/go/influx"
	"github.com/nullstyle/go/influx/influxtest"
	"github.com/stretchr/testify/assert"
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
		}, {
			Name:     "1.01+2",
			Actions:  []influx.Action{One, Dot, Zero, One, Add{}, Two, Equals{}},
			Expected: "3.01",
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
		{
			Name: "plus-minus",
			Actions: []influx.Action{
				One, PlusMinus{},
			},
			Expected: "-1",
		},
		{
			Name: "plus-minus: after calc",
			Actions: []influx.Action{
				One, Add{}, One, Equals{}, PlusMinus{},
			},
			Expected: "-2",
		},
		{
			Name: "clear",
			Actions: []influx.Action{
				One, One, Clear{},
			},
			Expected: "",
		},
	}

	for _, kase := range cases {
		t.Run(kase.Name, func(t *testing.T) {

			// setup
			var c Calculator
			influxtest.NewFromState(t, &c, kase.Actions...)

			// check result
			actual := c.Display()
			assert.Equal(t,
				kase.Expected, actual,
				"bad: %s != %s", kase.Name, actual,
			)
		})
	}

	// failure cases: every action but the last should succeed
	fails := []struct {
		Name    string
		Actions []influx.Action
	}{
		{
			Name: "invalid decimal: left side",
			Actions: []influx.Action{
				One, Dot, Dot, One, Add{}, One, Equals{},
			},
		},
		{
			Name: "invalid decimal: right side",
			Actions: []influx.Action{
				One, Add{}, One, Dot, Dot, One, Equals{},
			},
		},
		{
			Name: "invalid decimal: plus-minus",
			Actions: []influx.Action{
				One, Dot, Dot, One, PlusMinus{},
			},
		},
		{
			Name: "invalid decimal: multiple operators",
			Actions: []influx.Action{
				One, Add{}, One, Dot, Dot, One, Add{},
			},
		},
		{
			Name: "invalid decimal: multiple operators",
			Actions: []influx.Action{
				One, Add{}, One, Add{}, One, Dot, Dot, One, Equals{},
			},
		},
	}

	for _, kase := range fails {
		t.Run(kase.Name, func(t *testing.T) {
			actions := kase.Actions
			setup := actions[1 : len(actions)-1]
			trigger := actions[len(actions)-1]

			// setup
			var c Calculator
			store := influxtest.NewFromState(t, &c, setup...)

			err := store.Dispatch(trigger)
			assert.Error(t, err)

		})
	}
}
