package cmd

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"

	"path/filepath"

	"github.com/nullstyle/go/env"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build [PATH]",
	Short: "Build creates the gopm bundle for the app at PATH",
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

		var buf bytes.Buffer

		npmDeps, err := getPackageJSON(pkg)
		if err != nil {
			log.Fatal(err)
		}

		// preamble
		fmt.Fprintln(&buf, "var gopm_modules = {};")

		for mod := range npmDeps.Dependencies {

			fmt.Fprintf(
				&buf,
				"gopm_modules[\"%s\"] = require(\"%s\");\n",
				mod,
				mod)
		}

		// postamble
		fmt.Fprintln(&buf, "global[\"gopm_modules\"] = gopm_modules;")

		err = gotoPkgDir(pkg)
		if err != nil {
			log.Fatal(err)
		}

		afero.WriteFile(env.FS, "tmp-gopm.js", buf.Bytes(), 0644)
		defer env.FS.Remove("tmp-gopm.js")

		bargs := []string{"tmp-gopm.js"}
		if output != "-" {
			realo, err := env.RealPath(output)
			if err != nil {
				log.Fatal(err)
			}
			abso, err := filepath.Abs(realo)
			if err != nil {
				log.Fatal(err)
			}
			bargs = append(bargs, "-o", abso)
		}

		bcmd := exec.Command("browserify", bargs...)
		bcmd.Stdout = os.Stdout
		bcmd.Stderr = os.Stderr

		err = bcmd.Run()
		if err != nil {
			log.Fatal(err)
		}

	},
}

func init() {
	RootCmd.AddCommand(buildCmd)

	buildCmd.Flags().StringVarP(
		&output,
		"output",
		"o",
		"gopm_generated.inc.js",
		"output path (use '-' for stdout)")
}
