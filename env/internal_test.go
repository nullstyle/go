package env

import "github.com/nullstyle/go/env/mocks"

func mockBackend(fn func(*mocks.Backend)) {
	be := &mocks.Backend{}
	old := DefaultBackend

	DefaultBackend = be
	fn(be)
	DefaultBackend = old
}
