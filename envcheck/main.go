// Package envcheck implements functions that make assertions on the state of
// the local execution environment.
package envcheck

import (
	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

//go:generate mockery -all

var DefaultPathLooker = OS
var DefaultFS = afero.NewOsFs()
var OS Backend = &osBackend{}

// Backend is a backend
type Backend interface {
	PathLooker
}

// PathLooker represents a type that can look paths up from the local execution
// environment.
type PathLooker interface {

	// LookupPath searches the local environment for program, returning the
	// absolute path for the program.
	LookupPath(program string) (string, error)
}

// Executable asserts that program is present and executable on the local
// system. Returns the resolved path to the program.
func Executable(program string) (string, error) {
	path, err := DefaultPathLooker.LookupPath(program)
	if err != nil {
		return "", errors.Wrap(err, "lookup path failed")
	}

	file, err := DefaultFS.Stat(path)
	if err != nil {
		return "", errors.Wrap(err, "file stat failed")
	}

	if (file.Mode() & 0111) == 0 {
		return "", errors.New("error not executable")
	}

	return path, nil
}
