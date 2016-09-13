package sci

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValue_Convert(t *testing.T) {
	sys := testSystem()
	v1, err := sys.Parse("10 foot")
	require.NoError(t, err)
	target, err := sys.ParseUnit("meter")
	require.NoError(t, err)

	var result Value
	err = result.Convert(v1, target)

	if assert.NoError(t, err) {
		assert.Equal(t, "3.048", result.M)
		assert.Equal(t, target, result.U)
	}
}

func TestValue_String(t *testing.T) {
	sys := testSystem()
	cases := []struct {
		Name     string
		In       string
		Expected string
	}{
		{"singular", "1 foot", "1 foot"},
		{"singular-wrong", "1 feet", "1 foot"},
		{"plural", "10 feet", "10 feet"},
		{"plural-wrong", "10 foot", "10 feet"},
	}

	for _, kase := range cases {
		t.Run(kase.Name, func(t *testing.T) {
			v := sys.MustParse(kase.In)
			assert.Equal(t, kase.Expected, v.String())
		})
	}
}
