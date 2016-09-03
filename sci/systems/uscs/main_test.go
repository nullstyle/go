package uscs

import (
	"testing"

	"github.com/nullstyle/go/sci"
	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	cases := []struct {
		Name              string
		Value             string
		ExpectedMagnitude string
		ExpectedUnit      sci.Unit
	}{
		{"simple", "10 inch", "10", Inch},
	}

	for _, kase := range cases {
		t.Run(kase.Name, func(t *testing.T) {
			v, err := Parse(kase.Value)
			if assert.NoError(t, err) {
				assert.Equal(t, kase.ExpectedMagnitude, v.M)
				assert.Equal(t, kase.ExpectedUnit, v.U)
			}
		})
	}
}
