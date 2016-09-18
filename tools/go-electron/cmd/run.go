package cmd

import (
	"log"
	"os"
	"os/exec"
	"syscall"

	"github.com/nullstyle/go/electron/build"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run [PATH]",
	Short: "run the go-electron app at PATH",
	Long:  `run compiles app at PATH and then starts it up`,
	Run: func(cmd *cobra.Command, args []string) {

		dir, err := build.Run(args[0], "all", "all")
		if err != nil {
			log.Fatal(err)
		}

		electron, err := exec.LookPath("electron")
		if err != nil {
			log.Fatal(err)
		}

		err = os.Chdir(dir)
		if err != nil {
			log.Fatalln("couldn't change to build dir:", err)
		}

		eargs := []string{"electron", "."}
		env := os.Environ()
		err = syscall.Exec(electron, eargs, env)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
