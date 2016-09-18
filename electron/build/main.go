// Package build implements functions for building go-electron applications
package build

import (
	"path/filepath"

	"github.com/nullstyle/go/env"
	"github.com/nullstyle/go/gopherjs"
	"github.com/pkg/errors"
)

// Run builds `pkg` as an electron app for the os/arch pair provided, writing
// the result to the build directory underneath pkg.
func Run(pkg string, os string, arch string) error {
	path, err := env.PkgExists(pkg)
	if err != nil {
		return errors.Wrap(err, "env/PkgExists failed")
	}

	outDir := filepath.Join(path, ".go-electron")
	err = env.FS.MkdirAll(outDir, 0755)
	if err != nil {
		return errors.Wrap(err, "failed to make build dir")
	}

	outPath := filepath.Join(outDir, "main.js")

	err = gopherjs.Build(pkg, outPath)
	if err != nil {
		return errors.Wrap(err, "compile js failed")
	}

	return nil
}
