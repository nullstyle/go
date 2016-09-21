package influx

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {
	cases := []struct {
		Name     string
		Error    HookError
		Expected string
	}{
		{
			"index",
			HookError{
				Index: 1,
				Err:   errors.New("boom"),
				Hook:  struct{}{},
			},
			"hook [1] failed: boom",
		},
		{
			"named",
			HookError{
				Index: 1,
				Err:   errors.New("boom"),
				Hook:  &TestHook{},
			},
			"hook `test-hook` failed: boom",
		},
		{
			"nil err",
			HookError{
				Index: 1,
				Hook:  &TestHook{},
			},
			"hook `test-hook` failed: %!s(<nil>)",
		},
	}

	for _, kase := range cases {
		t.Run(kase.Name, func(t *testing.T) {
			assert.Equal(t, kase.Expected, kase.Error.Error())
		})
	}
}
