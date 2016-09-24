package influx

import "fmt"

// ActionError represents an error triggered during the application of an action
// to a store.
type ActionError struct {
	Action Action
	Store  *Store
	Err    error
}

// Error implements error
func (err *ActionError) Error() string {
	return fmt.Sprintf("action error: %s", err.Err)
}

// HookError represents an error that occurred while running a hook function
type HookError struct {
	Index int
	Hook  Hook
	Err   error
}

// Error implements error
func (herr *HookError) Error() string {
	name := fmt.Sprintf("[%d]", herr.Index)

	if named, ok := herr.Hook.(Named); ok {
		name = fmt.Sprintf("`%s`", named.Name())
	}

	return fmt.Sprintf("hook %s failed: %s", name, herr.Err)
}

var _ error = &ActionError{}
var _ error = &HookError{}
