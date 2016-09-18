// Package gopherjs implements functions for working with gopherjs packages
package gopherjs

import (
	"path/filepath"

	"os/exec"

	"github.com/nullstyle/go/env"
	"github.com/pkg/errors"
)

type BuildError struct {
	Stdout  []byte
	ExitErr *exec.ExitError
}

// Build builds the pkg into javascript
func Build(pkg string) error {
	path, err := env.PkgExists(pkg)
	if err != nil {
		return errors.Wrap(err, "env/PkgExists failed")
	}

	outPath := filepath.Join(path, "main.js")
	cmd := exec.Command("gopherjs", "build", pkg, "-o", outPath)

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
