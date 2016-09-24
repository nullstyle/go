package trace

import (
	"bytes"
	"io/ioutil"
	"testing"
	"time"

	"encoding/json"

	"github.com/nullstyle/go/influx/influxtest"
	"github.com/stretchr/testify/assert"
)

func TestSnapshot_Save_Smoke(t *testing.T) {
	hook := &Hook{
		MaxSize:    DefaultMaxSize,
		TargetSize: DefaultMaxSize,
		TargetAge:  1 * time.Minute,
	}
	store := influxtest.New(t)
	store.UseHooks(hook)

	influxtest.Do(t, store,
		influxtest.ActionInc,
		influxtest.ActionInc,
		influxtest.ActionInc,
		influxtest.ActionInc,
	)

	var saved bytes.Buffer
	err := hook.Current.Save(&saved)
	if assert.NoError(t, err) {
		var loaded Snapshot
		err := json.Unmarshal(saved.Bytes(), &loaded)
		assert.NoError(t, err)

		assert.Len(t, loaded.Dispatches, 4)
	}

	// saving a blank snapshot fails
	snap := &Snapshot{}
	err = snap.Save(ioutil.Discard)
	assert.EqualError(t, err, "empty initial state")
}

func TestSnapshot_Age(t *testing.T) {

	s := Snapshot{}
	now := time.Now()
	expectedAge := 5 * time.Minute
	s.InitialState.CreatedAt = now.Add(-expectedAge)

	actualAge := s.Age()
	assert.InDelta(t, expectedAge.Seconds(), actualAge.Seconds(), 1)
}
