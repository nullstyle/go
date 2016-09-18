// Package build implements functions for building go-electron applications
package build

import (
	"path/filepath"

	"strings"

	"github.com/nullstyle/go/env"
	"github.com/nullstyle/go/gopherjs"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

//go:generate go-bindata -pkg=build skel/...

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

	err = copyAsset("skel/index.html", outDir)
	if err != nil {
		return errors.Wrap(err, "copy skell failed")
	}

	return nil
}

func copyAsset(asset string, outDir string) error {
	fi, err := AssetInfo(asset)

	// copy children if current asset is a directory
	if fi.IsDir() {
		children, err := AssetDir(asset)
		if err != nil {
			return errors.Wrap(err, "AssetDir failed")
		}

		for _, child := range children {
			childPath := filepath.Join(asset, child)
			err = copyAsset(childPath, outDir)
			if err != nil {
				return errors.Wrap(err, "failed copying child")
			}
		}

		return nil
	}

	// otherwise copy the file
	raw, err := Asset(asset)
	if err != nil {
		return errors.Wrap(err, "Asset failed")
	}

	// strip off first directory on the asset path, since it represents the
	// bindata folder we don't care about
	assetPieces := strings.Split(asset, string(filepath.Separator))
	targetPath := filepath.Join(assetPieces[1:]...)
	outPath := filepath.Join(outDir, targetPath)

	err = env.FS.MkdirAll(outDir, 0755)
	if err != nil {
		return errors.Wrap(err, "failed creating outDir")
	}

	err = afero.WriteFile(env.FS, outPath, raw, fi.Mode())
	if err != nil {
		return errors.Wrap(err, "failed writing asset")
	}

	return nil
}
