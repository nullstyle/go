package sci

import "fmt"

// Error implements the error interface
func (experr *ExpToBigError) Error() string {
	return fmt.Sprintf(
		"exponent (%d) is to large, max allowed is (%d)",
		experr.Exp,
		MaxExp,
	)
}
