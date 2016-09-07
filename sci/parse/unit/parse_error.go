package unit

import "fmt"

// Error implements the error interface
func (perr *ParseError) Error() string {
	return fmt.Sprintf("parse unit: %s failed", perr.FailurePhase)
}
