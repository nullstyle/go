package cmd

import (
	"log"

	"github.com/nullstyle/go/electron/build"
	"github.com/nullstyle/go/gopherjs"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build [PATH]",
	Short: "Build the go-electron app at PATH",
	Long:  `TODO`,
	Run: func(cmd *cobra.Command, args []string) {
		_, err := build.Run(args[0], "all", "all")
		switch err := errors.Cause(err).(type) {
		case nil:
			// noop
		case *gopherjs.BuildError:
			log.Fatalf("gopherjs build error: %s", err.ExitErr.Stderr)
		default:
			log.Fatal(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(buildCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// buildCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// buildCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
