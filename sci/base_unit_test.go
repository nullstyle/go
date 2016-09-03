package sci

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBaseUnit_PopulateNormalizedUnit(t *testing.T) {
	sys := testSystem()
	var nu NormalizedUnit

	l := sys.BaseUnits[Length]

	l.PopulateNormalizedUnit(&nu, false)
	assert.Equal(t, 1, nu.Components[l])
	l.PopulateNormalizedUnit(&nu, true)
	assert.Equal(t, 0, nu.Components[l])
	l.PopulateNormalizedUnit(&nu, true)
	assert.Equal(t, -1, nu.Components[l])
}
