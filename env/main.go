// package env implements functions that make assertions on the state of
// the local execution environment.
package env

import (
	"os"
	"path/filepath"
	"strings"

	"os/exec"

	"github.com/nullstyle/go/gopath"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

//go:generate mockery -name Backend

// DefaultEnvGetter looks up environment variables from the local system
var DefaultEnvGetter = OS

// DefaultPathLooker looks up paths using the golang stdlib
var DefaultPathLooker = OS

//FS uses the local machine's filesystem
var FS = afero.NewOsFs()

// OS represents a backend that uses the real os
var OS Backend = &osBackend{}

// Backend is a backend
type Backend interface {
	PathLooker
	EnvGetter
}

// EnvGetter represents a type that can lookup an environment variable
type EnvGetter interface {
	Getenv(string) string
}

// PathLooker represents a type that can look paths up from the local execution
// environment.
type PathLooker interface {

	// LookupPath searches the local environment for program, returning the
	// absolute path for the program.
	LookupPath(program string) (string, error)
}

// PkgNotFoundError is the error returns when a package cannot be found in the
// GOPATH.
type PkgNotFoundError struct {
	Pkg string
}

// Executable asserts that program is present and executable on the local
// system. Returns the resolved path to the program.
func Executable(program string) (string, error) {
	path, err := DefaultPathLooker.LookupPath(program)
	if err != nil {
		return "", errors.Wrap(err, "lookup path failed")
	}

	file, err := FS.Stat(path)
	if err != nil {
		return "", errors.Wrap(err, "file stat failed")
	}

	if (file.Mode() & 0111) == 0 {
		return "", errors.New("error not executable")
	}

	return path, nil
}

// IsPkgNotFound returns true if err's cause is of type PkgNotFoundError
func IsPkgNotFound(err error) bool {
	_, ok := errors.Cause(err).(*PkgNotFoundError)
	return ok
}

// NpmPkgExists ensures an npm package is installed globally.  TODO: make
// testable.
func NpmPkgExists(pkg string) (string, error) {
	out, err := exec.Command("npm", "list", "--parseable", "-g", pkg).Output()
	_, ok := err.(*exec.ExitError)
	if ok {
		return "", &PkgNotFoundError{pkg}
	}

	if err != nil {
		return "", errors.Wrap(err, "npm exec failed")
	}

	return strings.TrimSpace(string(out)), nil
}

// PkgExists returns the absolute path of pkg and ensures that the package
// is present on the local file system underneath GOPATH.
func PkgExists(pkg string) (string, error) {
	gpath := DefaultEnvGetter.Getenv("GOPATH")
	paths := gopath.Split(gpath)

	// search each path in order for the package
	for _, path := range paths {
		full := filepath.Join(path, "src", pkg)

		_, err := FS.Stat(full)
		if os.IsNotExist(err) {
			continue
		}

		if err != nil {
			return "", errors.Wrap(err, "stat failed")
		}

		return full, nil
	}

	return "", &PkgNotFoundError{pkg}
}

// PkgPath returns the absolute path of pkg and ensures that the package
// is present on the local file system underneath GOPATH.
func PkgPath(pkg string) (string, error) {
	path, err := PkgExists(pkg)
	if IsPkgNotFound(err) {
		gpath := DefaultEnvGetter.Getenv("GOPATH")
		return filepath.Join(gopath.First(gpath), "src", pkg), nil
	}

	if err != nil {
		return "", errors.Wrap(err, "envcheck/PkgExists failed")
	}

	return path, nil
}
