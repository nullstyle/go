package trace

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"time"

	"github.com/miolini/datacounter"
	"github.com/nullstyle/go/influx"
	"github.com/pkg/errors"
)

// Age returns the current age of the snapshot
func (s *Snapshot) Age() time.Duration {
	return time.Since(s.InitialState.CreatedAt)
}

// Checkpoint records the
func (s *Snapshot) Checkpoint(store *influx.Store) error {
	s.InitialState.CreatedAt = time.Now()

	var err error
	s.InitialState, err = store.TakeSnapshot()
	if err != nil {
		return errors.Wrap(err, "store snapshot failed")
	}

	// reset the dispatches record
	s.Dispatches = s.Dispatches[0:0]

	return nil
}

// Empty returns true if the snapshot is considered empty.  A snapshot is
// considered empty if it has no initial state defined.
func (s *Snapshot) Empty() bool {
	return len(s.InitialState.State) == 0
}

// Size returns the size of an encoded snapshot in bytes.  NOTE: this
func (s *Snapshot) Size() (uint64, error) {
	if s.Empty() {
		return 0, nil
	}

	counter := datacounter.NewWriterCounter(ioutil.Discard)

	err := s.Save(counter)
	if err != nil {
		return 0, errors.Wrap(err, "save failed")
	}

	return counter.Count(), nil
}

// Save serializes the snapshot and writes it to the provided writer
func (s *Snapshot) Save(w io.Writer) error {
	if len(s.InitialState.State) == 0 {
		return errors.New("empty initial state")
	}

	enc := json.NewEncoder(w)
	err := enc.Encode(s)

	if err != nil {
		return errors.Wrap(err, "marshal failed")
	}

	return nil
}
