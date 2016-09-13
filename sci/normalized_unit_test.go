package sci

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizedUnit_Add(t *testing.T) {
	var nu NormalizedUnit

	nu.Add(Length, 2)

	assert.Equal(t, 2, nu.Components[Length])
}
