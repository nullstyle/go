package test

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"path/filepath"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

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
}
