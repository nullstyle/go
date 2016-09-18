package env

import (
	"testing"

	"strings"

	"github.com/nullstyle/go/env/mocks"
	"github.com/nullstyle/go/test"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestExecutable(t *testing.T) {
	fs, done := test.FS(t, "envcheck")
	defer done()

	be := &mocks.Backend{}

	DefaultPathLooker = be
	FS = fs

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

func TestIsPkgNotFound(t *testing.T) {
	cases := []struct {
		Name     string
		Err      error
		Expected bool
	}{
		{"happy: bare", &PkgNotFoundError{"foo"}, true},
		{"happy: wrapped", errors.Wrap(&PkgNotFoundError{"foo"}, "once"), true},
		{
			"happy: wrapped-twice",
			errors.Wrap(errors.Wrap(&PkgNotFoundError{"foo"}, "once"), "twiice"),
			true,
		},
		{"sad: nil", nil, false},
		{"sad: other", errors.New("foo"), false},
	}

	for _, kase := range cases {
		t.Run(kase.Name, func(t *testing.T) {
			ret := IsPkgNotFound(kase.Err)
			assert.Equal(t, kase.Expected, ret)
		})
	}

}

func TestPkgExists(t *testing.T) {
	fs, done := test.FS(t, "envcheck")
	defer done()
	be := &mocks.Backend{}
	DefaultPathLooker = be
	DefaultEnvGetter = be
	FS = fs

	test.WriteFile(t, fs, "/go1/src/pkg1/main.go", "", 0644)
	test.WriteFile(t, fs, "/go2/src/pkg2/main.go", "", 0644)

	be.On("Getenv", "GOPATH").Return("/go1:/go2")

	// happy path: in first entry
	path, err := PkgExists("pkg1")
	if assert.NoError(t, err) {
		assert.Equal(t, "/go1/src/pkg1", path)
	}
	// happy path: in second entry
	path, err = PkgExists("pkg2")
	if assert.NoError(t, err) {
		assert.Equal(t, "/go2/src/pkg2", path)
	}

	// sad path: doesn't exist
	_, err = PkgExists("missing")
	assert.Error(t, err)

	be.AssertExpectations(t)
}

func TestPkgPath(t *testing.T) {
	fs, done := test.FS(t, "envcheck")
	defer done()
	be := &mocks.Backend{}
	DefaultPathLooker = be
	DefaultEnvGetter = be
	FS = fs

	test.WriteFile(t, fs, "/go1/src/pkg1/main.go", "", 0644)
	test.WriteFile(t, fs, "/go2/src/pkg2/main.go", "", 0644)

	be.On("Getenv", "GOPATH").Return("/go1:/go2")

	// happy path: in first entry
	path, err := PkgPath("pkg1")
	if assert.NoError(t, err) {
		assert.Equal(t, "/go1/src/pkg1", path)
	}
	// happy path: in second entry
	path, err = PkgPath("pkg2")
	if assert.NoError(t, err) {
		assert.Equal(t, "/go2/src/pkg2", path)
	}

	// happy path: new pkg
	path, err = PkgPath("pkg3")
	if assert.NoError(t, err) {
		assert.Equal(t, "/go1/src/pkg3", path)
	}
}

func TestRealPath(t *testing.T) {
	fs, done := test.FS(t, "envcheck")
	defer done()

	// happy path: local fs
	FS = afero.NewOsFs()

	real, err := RealPath("/bin")
	if assert.NoError(t, err) {
		assert.Equal(t, "/bin", real)
	}

	// happy path: test fs
	FS = fs
	real, err = RealPath("/bin")
	if assert.NoError(t, err) {
		assert.True(t, strings.HasSuffix(real, "/bin"), "real path doesn't end with /bin")
	}
}
