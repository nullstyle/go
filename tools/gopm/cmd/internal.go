package cmd

import (
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"

	"github.com/nullstyle/go/env"
	"github.com/pkg/errors"
)

var output string

type packageJson struct {
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
}

func newPackage() packageJson {
	return packageJson{
		Dependencies:    map[string]string{},
		DevDependencies: map[string]string{},
	}
}

func deps(pkg string) ([]string, []string, error) {

	// onedep is the self-recursive function that
	// populates the found map by using `go list``
	var onedep func(pkg string, ret map[string]struct{}) error
	onedep = func(pkg string, ret map[string]struct{}) error {
		// get test
		i, _, err := imports(pkg)
		if err != nil {
			return errors.Wrap(err, "load imports failed")
		}

		for _, dep := range i {
			_, ok := ret[dep]
			if ok {
				return nil
			}

			ret[dep] = struct{}{}
			err := onedep(dep, ret)
			if err != nil {
				return err
			}
		}
		return nil
	}

	collect := func(in map[string]struct{}) (out []string) {
		for path := range in {
			out = append(out, path)
		}
		sort.Strings(out)
		return
	}

	importMap := map[string]struct{}{}
	testImportMap := map[string]struct{}{}

	err := onedep(pkg, importMap)
	if err != nil {
		return nil, nil, errors.Wrap(err, "load imports failed")
	}

	_, ti, err := imports(pkg)
	for _, tdep := range ti {
		testImportMap[tdep] = struct{}{}
		err := onedep(tdep, testImportMap)
		if err != nil {
			return nil, nil, errors.Wrap(err, "load test imports failed")
		}
	}

	imports := collect(importMap)
	testImports := collect(testImportMap)

	return imports, testImports, nil
}

func expandPkg(arg string) string {
	pkg, err := env.ExpandPkg(arg)
	switch err := errors.Cause(err).(type) {
	case nil:
		return pkg
	case *env.NotOnGoPathError:
		log.Fatalf("bad path: %s", err)
	default:
		log.Fatal(err)
	}
	return ""
}

func imports(pkg string) ([]string, []string, error) {
	raw, err := exec.Command("go", "list", "-tags", "'js'", "-json", pkg).Output()
	if err != nil {
		return nil, nil, errors.Wrap(err, "run go list failed")
	}

	var parsed struct {
		Imports     []string
		TestImports []string
	}
	err = json.Unmarshal(raw, &parsed)
	if err != nil {
		return nil, nil, errors.Wrap(err, "parse list output failed")
	}
	return parsed.Imports, parsed.TestImports, nil
}

func isGopherJS(pkg string) (bool, error) {
	jsonPath, err := jsonPath(pkg)
	if err != nil {
		return false, errors.Wrap(err, "json-path failed")
	}

	_, err = env.FS.Stat(jsonPath)
	switch {
	case os.IsNotExist(err):
		return false, nil
	case err != nil:
		return false, errors.Wrap(err, "stat failed")
	default:
		return true, nil
	}
}

func jsonPath(pkg string) (string, error) {
	pkgPath, err := env.PkgPath(pkg)
	if err != nil {
		return "", errors.Wrap(err, "pkg-path failed")
	}

	return filepath.Join(pkgPath, "package.json"), nil
}

// jsonPkgs returns a copy of pkgs filtered by which have a package.json file
// within their directory
func jsonPkgs(pkgs []string) ([]string, error) {
	var ret []string
	for _, pkg := range pkgs {
		gjs, err := isGopherJS(pkg)
		if err != nil {
			return nil, err
		}

		if !gjs {
			continue
		}

		ret = append(ret, pkg)
	}

	return ret, nil
}
