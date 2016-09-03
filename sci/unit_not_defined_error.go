package sci

import "fmt"

// Error implements the error interface
func (uerr *UnitNotDefinedError) Error() string {
	return fmt.Sprintf("%s is not defined", uerr.Name)
}
