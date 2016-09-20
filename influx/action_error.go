package influx

import "fmt"

func (err *ActionError) Error() string {
	return fmt.Sprintf("action error: %s", err.Cause)
}
