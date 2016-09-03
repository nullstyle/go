package sci

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizedUnit_Add(t *testing.T) {
	sys := testSystem()
	var nu NormalizedUnit

	l := sys.BaseUnits[Length]
	nu.Add(l, 2)

	assert.Equal(t, 2, nu.Components[l])
}
