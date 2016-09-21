package influx

import "fmt"

func (herr *HookError) Error() string {
	name := fmt.Sprintf("[%d]", herr.Index)

	if named, ok := herr.Hook.(Named); ok {
		name = fmt.Sprintf("`%s`", named.Name())
	}

	return fmt.Sprintf("hook %s failed: %s", name, herr.Err)
}

var _ error = &HookError{}
