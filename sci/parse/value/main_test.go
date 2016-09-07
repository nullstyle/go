package value

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	cases := []struct {
		Name     string
		Input    string
		Expected V
	}{
		{
			"simple",
			"1 meter",
			V{M: "1", U: "meter"},
		},
		{
			"float",
			"1.01 meter",
			V{M: "1.01", U: "meter"},
		},
		{
			"-float",
			"-1.01 meter",
			V{M: "-1.01", U: "meter"},
		},
		{
			"sci",
			"1.234E5 meter",
			V{M: "1.234E5", U: "meter"},
		},
		{
			"-sci",
			"-1.234E5 meter",
			V{M: "-1.234E5", U: "meter"},
		},
		{
			"exp",
			"9.0123 (meter / second^2)",
			V{M: "9.0123", U: "(meter / second^2)"},
		},
	}

	for _, kase := range cases {
		t.Run(kase.Name, func(t *testing.T) {
			v, err := Parse(kase.Input)
			if assert.NoError(t, err) {
				assert.EqualValues(t, kase.Expected, v)
			}
		})
	}
}
