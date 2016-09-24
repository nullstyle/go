// Package influxfmt implements functions for formatting the state of an influx
// application into strings.
//
// Format strings are expressed using go templates, and values extracted from
// the provided context are provided to the template by this package.
package influxfmt

import (
	"fmt"
)

// NOTE TO SELF: aims to cooperate with the go stdlib fmt package, providing a superset when necessary to support the needs of the development experience.

// CRAZYIDEA: have developers use the dispatch ctx as the argument passed into
// influxfmt.* methods.  We extract values from the context and format them
// based upon the format string.
//
// influxfmt.Sprintf(ctx, "%i", )
//

// IDEA: use templates for formatting:
//
//   influxfmt.Sprintf(ctx, "{{.DispatchTime}")

// TODO: implement the standard functions such that we expose an api that allows developers to format actions in ways that make it easier to debug and develop code within the influx framework.

// what interfaces do we want to support: %a

// how do we register formatters such that I can easily add more in the future.

// Sprintf formats according to a format specifier and returns the resulting string.
func Sprintf(format string, a ...interface{}) string {
	// NOTE TO SELF: use the same design for presenting errors as the stdlib.Sprintf

	// interrogate each argument for an interfaces they might implement that influences an influx specific formatter, falling back to fmt.Sprintf if none are found

	ia := make([]interface{}, len(a))
	for i := range a {

		// search for formatting sequences, replace them with %s and wrap the provided argument in a struct that implement fmt.Stringer to provide the desired behavior.
		ia[i] = a[i]

	}

	return fmt.Sprintf(format, ia...)
}
