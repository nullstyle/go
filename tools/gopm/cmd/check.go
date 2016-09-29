package cmd

import (
	"log"

	"github.com/nullstyle/go/env"
	"github.com/spf13/cobra"
)

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "check the local execution environment",
	Long: `check ensures that the dependencies for gopm
  are available in the local environment.`,
	Run: func(cmd *cobra.Command, args []string) {
		ensureExecutable("npm")
		ensureExecutable("browserify")
	},
}

func init() {
	RootCmd.AddCommand(checkCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// checkCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// checkCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

// ensureExecutable exits the application with a failure if an application with
// the name provided is not in the local environment's PATH.
func ensureExecutable(name string) {
	_, err := env.Executable(name)
	if err != nil {
		log.Printf("program `%s` MISSING", name)
		log.Fatal(err)
	}
	log.Printf("[good] program `%s` found", name)
}

func ensureNpm(name string) {
	_, err := env.NpmPkgExists(name)
	if err != nil {
		log.Printf("npm `%s` MISSING", name)
		log.Fatal(err)
	}
	log.Printf("[good] npm `%s` found", name)
}
