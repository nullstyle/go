package cmd

import (
	"bytes"
	"fmt"
	"log"
	"path/filepath"

	"os"

	"os/exec"

	"github.com/nullstyle/go/env"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var output string

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build [PATH]",
	Short: "Build creates the gopm bundle for the app at PATH",
	Long:  `TODO`,
	Run: func(cmd *cobra.Command, args []string) {
		var pkgDir string
		var err error

		switch len(args) {
		case 0:
			pkgDir, err = os.Getwd()
			if err != nil {
				log.Fatal(err)
			}
		case 1:
			pkg := expandPkg(args[0])
			pkgDir, err = env.PkgPath(pkg)
			if err != nil {
				log.Fatal(err)
			}
		default:
			log.Fatal("too many args")
		}

		var buf bytes.Buffer

		nodeDir := filepath.Join(pkgDir, "node_modules")
		nodepkgs, err := afero.ReadDir(env.FS, nodeDir)
		if err != nil {
			log.Fatal(err)
		}

		// preamble
		fmt.Fprintln(&buf, "var gopm_modules = {};")

		for _, fi := range nodepkgs {
			if !fi.IsDir() {
				continue
			}
			fmt.Fprintf(
				&buf,
				"gopm_modules[\"%s\"] = require(\"%s\");\n",
				fi.Name(),
				fi.Name())
		}

		// postamble
		fmt.Fprintln(&buf, "global[\"gopm_modules\"] = gopm_modules;")

		afero.WriteFile(env.FS, "tmp-gopm.js", buf.Bytes(), 0644)
		// defer env.FS.Remove("tmp-gopm.js")

		err = exec.Command("browserify", "tmp-gopm.js", "-o", output).Run()
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
		"name of generated file")
}
