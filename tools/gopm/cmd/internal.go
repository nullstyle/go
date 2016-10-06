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
	"github.com/spf13/afero"
)

var output string
var verbose bool

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

// getPackageJSON returns either the manually written or automatically generated
// package.json contents for a gopherjs package.
func getPackageJSON(pkg string) (*packageJson, error) {
	gjs, err := isGopherJS(pkg)
	if err != nil {
		return nil, errors.Wrap(err, "isGopherJS failed")
	}

	if !gjs {
		return autoPackage(pkg)
	}

	jsonPath, err := jsonPath(pkg)
	if err != nil {
		return nil, errors.Wrap(err, "json-path failed")
	}

	var result packageJson
	raw, err := afero.ReadFile(env.FS, jsonPath)
	if err != nil {
		return nil, errors.Wrap(err, "read-file failed")
	}

	// load the dependnecies into a temporary location
	err = json.Unmarshal(raw, &result)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal failed")
	}

	return &result, nil
}

func autoPackage(pkg string) (*packageJson, error) {
	imports, testImports, err := deps(pkg)
	if err != nil {
		return nil, errors.Wrap(err, "load deps failed")
	}

	jsonImports, err := jsonPkgs(imports)
	if err != nil {
		return nil, errors.Wrap(err, "filter imports failed")
	}

	jsonTestImports, err := jsonPkgs(testImports)
	if err != nil {
		return nil, errors.Wrap(err, "filter test imports failed")
	}

	merged, err := mergeJSONDeps(jsonImports, jsonTestImports)
	if err != nil {
		return nil, errors.Wrap(err, "building package.json failed")
	}
	return merged, nil
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

func gotoPkgDir(pkg string) error {
	pkgPath, err := env.PkgPath(pkg)
	if err != nil {
		return errors.Wrap(err, "pkg-path failed")
	}

	err = os.Chdir(pkgPath)
	if err != nil {
		return errors.Wrap(err, "chdir failed")
	}

	return nil
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

func installModules(pkg string) error {
	gjs, err := isGopherJS(pkg)
	if err != nil {
		return errors.Wrap(err, "isGopherJS failed")
	}

	jsonPath, err := jsonPath(pkg)
	if err != nil {
		return errors.Wrap(err, "get package.json path failed")
	}

	// write a temporary package.json by calculating the
	// package.json for the package under test
	if !gjs {

		packageJSON, err := autoPackage(pkg)
		if err != nil {
			return errors.Wrap(err, "create package.json contents failed")
		}

		raw, err := json.MarshalIndent(packageJSON, "", "  ")
		if err != nil {
			return errors.Wrap(err, "marshal package.json failed")
		}

		err = afero.WriteFile(env.FS, jsonPath, raw, 0644)
		if err != nil {
			return errors.Wrap(err, "write package.json failed")
		}

		defer env.FS.Remove(jsonPath)
	}

	pkgPath := filepath.Dir(jsonPath)
	realPath, err := env.RealPath(pkgPath)
	if err != nil {
		return errors.Wrap(err, "resolve real package.json path failed")
	}

	icmd := exec.Command("npm", "i")
	icmd.Dir = realPath
	err = icmd.Run()
	if err != nil {
		// TODO: output the installs stderr
		return errors.Wrap(err, "`npm i` failed")
	}

	return nil
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

// reads the package.json files at imports and timports, merges them into a
// single json result.
func mergeJSONDeps(imports, timports []string) (*packageJson, error) {
	results := newPackage()

	load := func(pkgs []string, dest map[string]string) error {
		loaded := newPackage()
		for _, pkg := range pkgs {
			path, err := jsonPath(pkg)
			if err != nil {
				return errors.Wrap(err, "json-path failed")
			}

			raw, err := afero.ReadFile(env.FS, path)
			if err != nil {
				return errors.Wrap(err, "read-file failed")
			}

			// load the dependnecies into a temporary location
			err = json.Unmarshal(raw, &loaded)
			if err != nil {
				return errors.Wrap(err, "unmarshal failed")
			}

			// copy the loaded dependencies into the results
			for mod, version := range loaded.Dependencies {
				dest[mod] = version
			}
		}

		return nil
	}

	err := load(imports, results.Dependencies)
	if err != nil {
		return nil, errors.Wrap(err, "load imports failed")
	}

	err = load(timports, results.DevDependencies)
	if err != nil {
		return nil, errors.Wrap(err, "load test imports failed")
	}

	return &results, nil
}
