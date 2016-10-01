//+build js

package mithril

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	assert.Equal(t, "v0.2.5", Version())
}
