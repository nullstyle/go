package env

import "fmt"
import "strings"
import "github.com/nullstyle/go/gopath"

func (err *NotOnGoPathError) Error() string {
	return fmt.Sprintf("path %s not on GOPATH: %s", err.Path, strings.Join(err.GoPath, gopath.Sep))
}

var _ error = &NotOnGoPathError{}
