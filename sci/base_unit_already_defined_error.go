package sci

import "fmt"

// Error implements the error interface
func (berr *BaseUnitAlreadyDefinedError) Error() string {
	return fmt.Sprintf(
		"cannot redefine the base unit for %s: %s is already defined",
		berr.Existing.Measure,
		berr.Existing.Name,
	)
}
