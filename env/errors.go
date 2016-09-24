package env

import (
	"fmt"
	"strings"

	"github.com/nullstyle/go/gopath"
)

// NotOnGoPathError is returned when a process failed because the path provided
// was not on the GOPATH of the current process.
type NotOnGoPathError struct {
	Path   string
	GoPath []string
}

func (err *NotOnGoPathError) Error() string {
	return fmt.Sprintf("path %s not on GOPATH: %s", err.Path, strings.Join(err.GoPath, gopath.Sep))
}

// PkgNotFoundError is the error returns when a package cannot be found in the
// GOPATH.
type PkgNotFoundError struct {
	Pkg string
}

func (err *PkgNotFoundError) Error() string {
	return fmt.Sprintf("cannot find pkg %s", err.Pkg)
}

var _ error = &NotOnGoPathError{}
var _ error = &PkgNotFoundError{}
