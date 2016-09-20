package influx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	// works with pointers
	var (
		base TestState
		num  int
		pnum *int
		anon struct {
			Counter int
		}
	)

	New(&base)
	New(&num)
	New(&pnum)
	New(&anon)

	// panics when not a pointer
	assert.Panics(t, func() {
		New(TestState{})
	})
}
