package cmd

import (
	"log"

	"github.com/nullstyle/go/electron/build"
	"github.com/nullstyle/go/env"
	"github.com/nullstyle/go/gopherjs"
	"github.com/pkg/errors"
)

func buildPkg(pkg string) string {
	dir, err := build.Run(pkg, "all", "all")
	switch err := errors.Cause(err).(type) {
	case nil:
		return dir
	case *gopherjs.BuildError:
		log.Fatalf("gopherjs build error: %s", err.ExitErr.Stderr)
	default:
		log.Fatal(err)
	}
	return ""
}

func expandPkg(arg string) string {
	pkg, err := env.ExpandPkg(arg)
	switch err := errors.Cause(err).(type) {
	case nil:
		return pkg
	case *env.NotOnGoPathError:
		log.Fatalf("bad path: %s", err)
	default:
		log.Fatal(err)
	}
	return ""
}
