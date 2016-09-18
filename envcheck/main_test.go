package envcheck

import (
	"errors"
	"testing"

	"io/ioutil"
	"os"

	"github.com/nullstyle/go/envcheck/mocks"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExecutable(t *testing.T) {
	dir, err := ioutil.TempDir(os.TempDir(), "envcheck-test")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	be := &mocks.Backend{}
	fs := afero.NewBasePathFs(afero.NewOsFs(), dir)

	DefaultPathLooker = be
	DefaultFS = fs

	err = fs.MkdirAll("/bin", 0777)
	require.NoError(t, err)
	err = afero.WriteFile(fs, "/bin/found", []byte(""), 0755)
	require.NoError(t, err)
	err = afero.WriteFile(fs, "/bin/non-executable", []byte(""), 0644)
	require.NoError(t, err)

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
