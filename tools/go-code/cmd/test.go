package cmd

import (
	"log"
	"os"
	"os/exec"
	"syscall"

	"github.com/nullstyle/go/env"
	"github.com/spf13/cobra"
	do "gopkg.in/godo.v2"
)

var testCmd = &cobra.Command{
	Use:   "test [PATH]",
	Short: "testing+vscode debugger",
	Long:  `Launch delve to allow visual studio code to debug the tests for package at PATH`,
	Run: func(cmd *cobra.Command, args []string) {

		pkg := expandPkg(args[0])
		electron, err := exec.LookPath("dlv")
		if err != nil {
			log.Fatal(err)
		}

		pkgpath, err := env.PkgPath(pkg)
		if err != nil {
			log.Fatal(err)
		}

		err = do.Inside(pkgpath, func() {
			eargs := []string{"dlv", "test", "--headless", "--listen=:2345", "--log"}
			env := os.Environ()
			err = syscall.Exec(electron, eargs, env)
			if err != nil {
				log.Fatal(err)
			}
		})
		if err != nil {
			log.Fatal(err)
		}

	},
}

func init() {
	RootCmd.AddCommand(testCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// buildCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// buildCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
