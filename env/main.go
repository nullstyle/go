// package env implements functions that make assertions on the state of
// the local execution environment.
package env

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/nullstyle/go/gopath"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

//go:generate mockery -name Backend

var DefaultBackend Backend = OS

//FS uses the local machine's filesystem
var FS = afero.NewOsFs()

// OS represents a backend that uses the real os
var OS Backend = &osBackend{}

// Backend is a backend
type Backend interface {
	Getwd() (string, error)
	Getenv(string) string

	// LookupPath searches the local environment for program, returning the
	// absolute path for the program.
	LookupPath(program string) (string, error)
}

// TODO: add BuildTime func

// Executable asserts that program is present and executable on the local
// system. Returns the resolved path to the program.
func Executable(program string) (string, error) {
	path, err := DefaultBackend.LookupPath(program)
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

// ExpandPkg takes a package spec that may be relative to the local directory
// and returns the full package name. Errors if the expanded path is not
// underneath an element of the current GOPATH.
func ExpandPkg(pkg string) (string, error) {
	if len(pkg) == 0 {
		return "", nil
	}

	// if the path isn't relative, treat it as a full path spec
	if pkg[0] != '.' {
		// TODO:  confirm pkg is in GOPATH
		return pkg, nil
	}

	wd, err := DefaultBackend.Getwd()
	if err != nil {
		return "", errors.Wrap(err, "get wd failed")
	}

	abs := wd + pkg[1:]
	gpath := DefaultBackend.Getenv("GOPATH")
	paths := gopath.Split(gpath)

	for _, path := range paths {
		// if abs pkg path is underneath the path
		if strings.HasPrefix(abs, path) {
			src := filepath.Join(path, "src")
			return strings.TrimPrefix(abs, src+"/"), nil
		}
	}

	return "", &NotOnGoPathError{Path: abs, GoPath: paths}
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
	gpath := DefaultBackend.Getenv("GOPATH")
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
		gpath := DefaultBackend.Getenv("GOPATH")
		return filepath.Join(gopath.First(gpath), "src", pkg), nil
	}

	if err != nil {
		return "", errors.Wrap(err, "env/PkgExists failed")
	}

	return path, nil
}

// RealPath resolves path against the current FS implementation, returning the
// actual local filesystem path if possible.  Useful for calling out to external
// commands that requires paths while still being able to override FS in tests.
func RealPath(path string) (string, error) {
	bp, ok := FS.(*afero.BasePathFs)
	if !ok {
		return path, nil
	}

	real, err := bp.RealPath(path)
	if err != nil {
		return "", errors.Wrap(err, "real path failed")
	}

	return real, nil
}

// Version returns the version string that was included when building the
// current program, or "devel" should one not be set.
func Version() string {
	return version
}
