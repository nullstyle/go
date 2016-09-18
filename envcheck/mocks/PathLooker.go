package mocks

import "github.com/stretchr/testify/mock"

// PathLooker is an autogenerated mock type for the PathLooker type
type PathLooker struct {
	mock.Mock
}

// LookupPath provides a mock function with given fields: program
func (_m *PathLooker) LookupPath(program string) (string, error) {
	ret := _m.Called(program)

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(program)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(program)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
