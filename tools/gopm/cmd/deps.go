package cmd

import (
	"fmt"
	"log"

	"encoding/json"

	"os"

	"github.com/nullstyle/go/env"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

// depsCmd represents the deps command
var depsCmd = &cobra.Command{
	Use:   "deps [PATH]",
	Short: "deps outputs a package.json for PATH containing its dependencies",
	Long:  `TODO`,
	Run: func(cmd *cobra.Command, args []string) {
		var pkg string
		var err error

		switch len(args) {
		case 0:
			pkg = expandPkg(".")
		case 1:
			pkg = expandPkg(args[0])
		default:
			log.Fatal("too many args")
		}

		gjs, err := isGopherJS(pkg)
		if err != nil {
			log.Fatal(err)
		}

		if gjs {
			log.Fatalf("pkg `%s` has curated package.json", pkg)
		}

		imports, testImports, err := deps(pkg)
		if err != nil {
			log.Fatal(err)
		}

		jsonImports, err := jsonPkgs(imports)
		if err != nil {
			log.Fatal(err)
		}

		jsonTestImports, err := jsonPkgs(testImports)
		if err != nil {
			log.Fatal(err)
		}

		merged, err := mergeJsonDeps(jsonImports, jsonTestImports)

		if output == "-" {
			fmt.Println(string(merged))
			os.Exit(0)
		}

		err = afero.WriteFile(env.FS, output, merged, 0644)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(depsCmd)

	depsCmd.Flags().StringVarP(
		&output,
		"output",
		"o",
		"-",
		"output path (use '-' for stdout)")
}

// reads the package.json files at imports and timports, merges them into a
// single json response.
func mergeJsonDeps(imports, timports []string) ([]byte, error) {
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

	ret, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return nil, errors.Wrap(err, "marshal failed")
	}

	return ret, nil
}
