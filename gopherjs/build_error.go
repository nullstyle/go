package gopherjs

import "fmt"

// Error implements error
func (err *BuildError) Error() string {
	return fmt.Sprintf("gopherjs failed: %s", err.ExitErr.Error())
}

var _ error = &BuildError{}
