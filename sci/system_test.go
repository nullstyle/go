package sci

import "testing"
import "github.com/stretchr/testify/assert"

func TestParse_Magnitudes(t *testing.T) {
	sys := testSystem()

	cases := []struct {
		Name     string
		Input    string
		Expected string
	}{
		{"integral", "10 meter", "10"},
		{"integral-negative", "-10 meter", "-10"},
		{"integral-nospace", "10meter", "10"},
		{"integral-negative-nospace", "-10meter", "-10"},
		{"floating", "10.001 meter", "10.001"},
		{"floating-negative", "-10.001 meter", "-10.001"},
		{"floating-nospace", "10.001meter", "10.001"},
		{"floating-negative-nospace", "-10.001meter", "-10.001"},
	}

	for _, kase := range cases {
		t.Run(kase.Name, func(t *testing.T) {
			v, err := sys.Parse(kase.Input)
			assert.NoError(t, err)
			assert.Equal(t, kase.Expected, v.M)
		})
	}

	// TODO:
	// assert.Equal(t, sys.MustLookup("meter"), v.U)

}
