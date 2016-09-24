package env

import (
	"testing"

	"github.com/nullstyle/go/test"
)

func TestNotOnGoPathError(t *testing.T) {
	test.Errors(t, []test.ErrorCase{
		{
			Name: "simple",
			Err: &NotOnGoPathError{
				Path:   "bar",
				GoPath: []string{"/go"},
			},
			Msg: "path bar not on GOPATH: /go",
		}, {
			Name: "multiple go paths",
			Err: &NotOnGoPathError{
				Path:   "bar",
				GoPath: []string{"/go:/go2"},
			},
			Msg: "path bar not on GOPATH: /go:/go2",
		},
	})
}

func TestPkgNotFoundError(t *testing.T) {
	test.Errors(t, []test.ErrorCase{
		{
			"simple",
			&PkgNotFoundError{
				Pkg: "blah",
			},
			"cannot find pkg blah",
		},
	})
}
