package sci

import "fmt"

// Error implements the error interface
func (uerr *UnitAlreadyDefinedError) Error() string {
	return fmt.Sprintf("cannot redefine %s", uerr.Name)
}
