// Package build implements functions for building go-electron applications
package build

import (
	"path/filepath"

	"os/exec"

	"github.com/nullstyle/go/env"
	"github.com/nullstyle/go/gopherjs"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

//go:generate go-bindata -pkg=build skel/...

// Run builds `pkg` as an electron app for the os/arch pair provided, writing
// the result to the build directory underneath pkg. returns the absolute path
// of the built electron directory.
func Run(pkg string, os string, arch string) (string, error) {
	path, err := env.PkgExists(pkg)
	if err != nil {
		return "", errors.Wrap(err, "env/PkgExists failed")
	}

	outDir := filepath.Join(path, ".go-electron")
	err = env.FS.MkdirAll(outDir, 0755)
	if err != nil {
		return "", errors.Wrap(err, "failed to make build dir")
	}

	// build gopherjs
	mainPath := filepath.Join(outDir, "main.js")
	err = gopherjs.Build(pkg, mainPath, true)
	if err != nil {
		return "", errors.Wrap(err, "compile main.js failed")
	}

	browserPath := filepath.Join(outDir, "browser.js")
	err = gopherjs.Build(pkg+"/browser", browserPath, true)
	if err != nil {
		return "", errors.Wrap(err, "compile browser.js failed")
	}

	err = writePackageJSON(path, outDir)
	if err != nil {
		return "", errors.Wrap(err, "add package.json failed")
	}

	err = writeIndex(path, outDir)
	if err != nil {
		return "", errors.Wrap(err, "create index.html failed")
	}

	return outDir, nil
}

func writePackageJSON(pkgPath string, outDir string) error {

	// TODO: don't assume main.go is the entry point
	main := filepath.Join(pkgPath, "main.go")
	raw, err := exec.Command("go", "run", main, "-writePackageJSON").Output()
	if err != nil {
		return errors.Wrap(err, "get package.json content failed")
	}

	outPath := filepath.Join(outDir, "package.json")

	err = afero.WriteFile(env.FS, outPath, raw, 0644)
	if err != nil {
		return errors.Wrap(err, "failed writing file")
	}

	return nil
}
