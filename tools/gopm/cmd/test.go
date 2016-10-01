package cmd

import (
	"log"
	"syscall"

	"os"
	"os/exec"
	"path/filepath"

	"github.com/nullstyle/go/env"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test [PATH]",
	Short: "Test runs the tests for the package at PATH",
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

		// write a temporary package.json by calculating the
		// package.json for the package under test
		if !gjs {
			jsonPath, err := jsonPath(pkg)
			if err != nil {
				log.Fatal(err)
			}

			packageJson, err := autoPackage(pkg)
			if err != nil {
				log.Fatal(err)
			}

			err = afero.WriteFile(env.FS, jsonPath, packageJson, 0644)
			if err != nil {
				log.Fatal(err)
			}

			defer env.FS.Remove(jsonPath)

			pkgPath := filepath.Dir(jsonPath)
			realPath, err := env.RealPath(pkgPath)
			if err != nil {
				log.Fatal(err)
			}

			err = os.Chdir(realPath)
			if err != nil {
				log.Fatal(err)
			}

			err = exec.Command("npm", "i").Run()
			if err != nil {
				log.Fatal(err)
			}
		}

		pkgPath, err := env.PkgPath(pkg)
		if err != nil {
			log.Fatal(err)
		}

		err = os.Chdir(pkgPath)
		if err != nil {
			log.Fatal(err)
		}

		gopherjs, err := exec.LookPath("gopherjs")
		if err != nil {
			log.Fatal(err)
		}

		eargs := []string{"gopherjs", "test"}

		if verbose {
			eargs = append(eargs, "-v")
		}

		env := os.Environ()
		err = syscall.Exec(gopherjs, eargs, env)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(testCmd)
	testCmd.Flags().BoolVarP(
		&verbose,
		"verbose",
		"v",
		false,
		"run tests in verbose mode",
	)
}
