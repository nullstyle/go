package trace

import (
	"testing"
	"time"

	"github.com/nullstyle/go/influx"
	"github.com/nullstyle/go/influx/influxtest"
	"github.com/stretchr/testify/assert"
)

func TestHook(t *testing.T) {
	influxtest.Hook(t, []influxtest.HookCase{
		{
			Name: "records actions",
			// Setup the tracing hook
			Hook: &Hook{
				TargetAge:  1 * time.Minute,
				TargetSize: DefaultMaxSize,
				MaxSize:    DefaultMaxSize,
			},

			// run the actions
			Actions: []influx.Action{
				influxtest.ActionInc,
				influxtest.ActionInc,
				influxtest.ActionDec,
			},

			// Assert hook state
			Test: func(hook *Hook) {
				t.Log(hook.Current.Dispatches)
				if assert.Len(t, hook.Current.Dispatches, 3) {
					disp := hook.Current.Dispatches
					assert.Equal(t, influxtest.ActionInc, disp[0].Action)
					assert.Equal(t, influxtest.ActionDec, disp[2].Action)
				}
			},
		},
	})
}

// TODO: add snapshotting tests, ensuring that
