package cmd

import (
	"log"

	"github.com/nullstyle/go/env"
	"github.com/pkg/errors"
)

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
