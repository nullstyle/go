package unit

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	cases := []struct {
		Name     string
		Input    string
		Expected U
	}{
		{
			"empty",
			"",
			&Nil{},
		},
		{
			"simple",
			"meter",
			&Ref{Name: "meter"},
		},
		{
			"div",
			"meter / sec",
			&Div{N: &Ref{Name: "meter"}, D: &Ref{Name: "sec"}},
		},
		{
			"mul",
			"meter * sec",
			&Mul{&Ref{Name: "meter"}, &Ref{Name: "sec"}},
		},
		{
			"mul-div",
			"meter * sec / sec",
			&Div{
				N: &Mul{&Ref{Name: "meter"}, &Ref{Name: "sec"}},
				D: &Ref{Name: "sec"},
			},
		},
		{
			"mul-div",
			"meter / sec * sec",
			&Mul{
				&Div{N: &Ref{Name: "meter"}, D: &Ref{Name: "sec"}},
				&Ref{Name: "sec"},
			},
		},
		{
			"parens",
			"meter * (sec / sec)",
			&Mul{
				&Ref{Name: "meter"},
				&Div{N: &Ref{Name: "sec"}, D: &Ref{Name: "sec"}},
			},
		},
		{
			"parens-2",
			"meter / (sec * sec)",
			&Div{
				N: &Ref{Name: "meter"},
				D: &Mul{&Ref{Name: "sec"}, &Ref{Name: "sec"}},
			},
		},
		{
			"exp",
			"meter / sec^2",
			&Div{
				N: &Ref{Name: "meter"},
				D: &Exp{U: &Ref{Name: "sec"}, Exp: 2},
			},
		},
		{
			"exp-parens",
			"(meter / sec)^2",
			&Exp{
				U:   &Div{N: &Ref{Name: "meter"}, D: &Ref{Name: "sec"}},
				Exp: 2,
			},
		},
	}

	for _, kase := range cases {
		t.Run(kase.Name, func(t *testing.T) {
			u, err := Parse(kase.Input)
			if assert.NoError(t, err) {
				assert.EqualValues(t, kase.Expected, u)
			}
		})
	}
}
