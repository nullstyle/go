package influx

import "reflect"

// handlert is the cached reflect.Type of the influx handler type. Used during
// dispatch to see if a given value or field implements Handler.
var handlert = reflect.TypeOf((*Handler)(nil)).Elem()
