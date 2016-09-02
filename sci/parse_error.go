package sci

import "fmt"

// Error implements the error interface
func (perr *ParseError) Error() string {
	return fmt.Sprintf("parse value: %s failed", perr.FailurePhase)
}
