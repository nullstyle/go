package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/nullstyle/go/env"
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

		npmDeps, err := autoPackage(pkg)
		if err != nil {
			log.Fatal(err)
		}

		raw, err := json.MarshalIndent(npmDeps, "", "  ")
		if err != nil {
			log.Fatal(err)
		}

		if output == "-" {
			fmt.Println(string(raw))
			os.Exit(0)
		}

		err = afero.WriteFile(env.FS, output, raw, 0644)
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
