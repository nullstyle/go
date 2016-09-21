package influx

import (
	"errors"
	"testing"

	"github.com/nullstyle/go/test"
)

func TestActionError(t *testing.T) {
	test.Errors(t, []test.ErrorCase{
		{
			"simple",
			&ActionError{
				Err: errors.New("boom"),
			},
			"action error: boom",
		},
	})
}

func TestHookError(t *testing.T) {
	test.Errors(t, []test.ErrorCase{
		{
			"hook: index-based",
			&HookError{
				Index: 1,
				Err:   errors.New("boom"),
				Hook:  struct{}{},
			},
			"hook [1] failed: boom",
		},
		{
			"named",
			&HookError{
				Index: 1,
				Err:   errors.New("boom"),
				Hook:  &TestHook{},
			},
			"hook `test-hook` failed: boom",
		},
		{
			"nil err",
			&HookError{
				Index: 1,
				Hook:  &TestHook{},
			},
			"hook `test-hook` failed: %!s(<nil>)",
		},
	})
}
