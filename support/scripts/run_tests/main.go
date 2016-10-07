package main

// See README.md for a description of this script

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/pkg/errors"
)

func main() {
	// collect packages to test using glide novendor
	// 
}

func build(pkg, dest, version, buildOS, buildArch string) {
	buildTime := time.Now().Format(time.RFC3339)

	timeFlag := fmt.Sprintf("-X github.com/nullstyle/go/env.buildTime=%s", buildTime)
	versionFlag := fmt.Sprintf("-X github.com/nullstyle/go/env.version=%s", version)

	if buildOS == "windows" {
		dest = dest + ".exe"
	}

	cmd := exec.Command("go", "build",
		"-o", dest,
		"-ldflags", fmt.Sprintf("%s %s", timeFlag, versionFlag),
		pkg,
	)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	cmd.Env = append(
		os.Environ(),
		fmt.Sprintf("GOOS=%s", buildOS),
		fmt.Sprintf("GOARCH=%s", buildArch),
	)
	log.Printf("building %s", pkg)

	log.Printf("running: %s", strings.Join(cmd.Args, " "))
	err := cmd.Run()
	if err != nil {
		log.Fatalln(err)
	}
}

// pushdir is a utility function to temporarily change directories.  It returns
// a func that can be called to restore the current working directory to the
// state it was in when first calling pushdir.
func pushdir(dir string) func() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(errors.Wrap(err, "getwd failed"))
	}

	err = os.Chdir(dir)
	if err != nil {
		panic(errors.Wrap(err, "chdir failed"))
	}

	return func() {
		err := os.Chdir(cwd)
		if err != nil {
			panic(errors.Wrap(err, "revert dir failed"))
		}
	}
}

// utility command to run the provided command that echoes any output.  A failed
// command will trigger a panic.
func run(name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	log.Printf("running: %s %s", name, strings.Join(args, " "))
	err := cmd.Run()

	if err != nil {
		panic(err)
	}
}
