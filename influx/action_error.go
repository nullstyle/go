package influx

import "fmt"

// Error implements error
func (err *ActionError) Error() string {
	return fmt.Sprintf("action error: %s", err.Cause)
}
