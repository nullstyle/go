package envcheck

import (
	"errors"
	"testing"

	"github.com/nullstyle/go/envcheck/mocks"
	"github.com/nullstyle/go/test"
	"github.com/stretchr/testify/assert"
)

func TestExecutable(t *testing.T) {
	fs, done := test.FS(t, "envcheck")
	defer done()

	be := &mocks.Backend{}

	DefaultPathLooker = be
	DefaultFS = fs

	test.WriteFile(t, fs, "/bin/found", "", 0755)
	test.WriteFile(t, fs, "/bin/non-executable", "", 0644)

	// happy path
	be.On("LookupPath", "foo").Return("/bin/found", nil)
	resolved, err := Executable("foo")
	if assert.NoError(t, err) {
		assert.Equal(t, "/bin/found", resolved)
	}

	// sad path: missing
	be.On("LookupPath", "missing").Return("", errors.New("not found"))
	_, err = Executable("missing")
	assert.Error(t, err)

	// sad path: not executable
	be.On("LookupPath", "non-executable").Return("/bin/non-executable", nil)
	_, err = Executable("non-executable")
	assert.EqualError(t, err, "error not executable")

	be.AssertExpectations(t)
}
