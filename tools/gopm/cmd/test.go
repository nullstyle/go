package cmd

import (
	"log"

	"os"
	"os/exec"

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

		err = installModules(pkg)
		if err != nil {
			log.Fatal(err)
		}

		err = gotoPkgDir(pkg)
		if err != nil {
			log.Fatal(err)
		}

		eargs := []string{"test"}

		if verbose {
			eargs = append(eargs, "-v")
		}

		gcmd := exec.Command("gopherjs", eargs...)
		gcmd.Stderr = os.Stderr
		gcmd.Stdout = os.Stdout

		err = gcmd.Run()
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
