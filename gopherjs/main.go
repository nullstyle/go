// Package gopherjs implements functions for working with gopherjs packages
package gopherjs

import (
	"os/exec"

	"github.com/nullstyle/go/env"
	"github.com/pkg/errors"
)

type BuildError struct {
	Stdout  []byte
	ExitErr *exec.ExitError
}

// Build builds the pkg into javascript
func Build(pkg string, outPath string) error {
	_, err := env.PkgExists(pkg)
	if err != nil {
		return errors.Wrap(err, "env/PkgExists failed")
	}

	realPath, err := env.RealPath(outPath)
	if err != nil {
		return errors.Wrap(err, "env/RealPath failed")
	}

	cmd := exec.Command("gopherjs", "build", pkg, "-o", realPath)

	raw, err := cmd.Output()
	if err != nil {
		eerr, ok := err.(*exec.ExitError)
		if ok {
			return &BuildError{
				Stdout:  raw,
				ExitErr: eerr,
			}
		}

		return errors.Wrap(err, "gopherjs failed")
	}

	return nil
}
