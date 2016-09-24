package test

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"path/filepath"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ErrorCase is a test case for an error's Error() method
type ErrorCase struct {
	Name string
	Err  error
	Msg  string
}

// Errors runs the provided error test cases
func Errors(t *testing.T, cases []ErrorCase) {
	for _, kase := range cases {
		t.Run(kase.Name, func(t *testing.T) {
			assert.Equal(t, kase.Msg, kase.Err.Error())
		})
	}
}

// FS creates a new tempdir for pkg and returns an FS rooted at the tempdir and
// a function that deletes the whole filesystem when called.
func FS(t *testing.T, pkg string) (afero.Fs, func()) {
	dir, err := ioutil.TempDir(os.TempDir(), fmt.Sprintf("%s-test", pkg))
	require.NoError(t, err)

	fs := afero.NewBasePathFs(afero.NewOsFs(), dir)

	return fs, func() {
		os.RemoveAll(dir)
	}
}

// WriteFile writes the provided contents at the provided path (while creating
// any intervening directories as needed) using the provided perms.
func WriteFile(
	t *testing.T,
	fs afero.Fs,
	path string,
	contents string,
	perm os.FileMode,
) {
	dir := filepath.Dir(path)
	err := fs.MkdirAll(dir, 0755)
	require.NoError(t, err)

	err = afero.WriteFile(fs, path, []byte(contents), perm)
	require.NoError(t, err)
}
