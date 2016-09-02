package sci

import "fmt"

// Error implements the error interface
func (merr *MagnitudeError) Error() string {
	return fmt.Sprintf("invalid magnitude: %s", merr.M)
}
