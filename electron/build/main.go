// Package build implements functions for building go-electron applications
package build

import (
	"path/filepath"

	"strings"

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
	err = gopherjs.Build(pkg, mainPath)
	if err != nil {
		return "", errors.Wrap(err, "compile main.js failed")
	}

	browserPath := filepath.Join(outDir, "browser.js")
	err = gopherjs.Build(pkg+"/browser", browserPath)
	if err != nil {
		return "", errors.Wrap(err, "compile browser.js failed")
	}

	// copy application skeleton
	err = copyAssetDir("skel", outDir)
	if err != nil {
		return "", errors.Wrap(err, "copy skell failed")
	}

	err = writePackageJSON(path, outDir)
	if err != nil {
		return "", errors.Wrap(err, "add package.json failed")
	}

	return outDir, nil
}

func copyAsset(asset string, outDir string) error {
	fi, err := AssetInfo(asset)
	if err != nil {
		return errors.Wrap(err, "AssetInfo failed")
	}

	// copy children if current asset is a directory
	if fi.IsDir() {
		return copyAssetDir(asset, outDir)
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

func copyAssetDir(dir string, outDir string) error {
	children, err := AssetDir(dir)
	if err != nil {
		return errors.Wrap(err, "AssetDir failed")
	}

	for _, child := range children {
		childPath := filepath.Join(dir, child)
		err = copyAsset(childPath, outDir)
		if err != nil {
			return errors.Wrap(err, "failed copying child")
		}
	}

	return nil
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
