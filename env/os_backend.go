package env

import (
	"os"
	"os/exec"

	"github.com/pkg/errors"
)

// LookPath implements PathLooker
func (be *osBackend) LookupPath(pogram string) (string, error) {
	path, err := exec.LookPath(pogram)
	if err != nil {
		return "", errors.Wrap(err, "exec/LookPath failed")
	}

	return path, nil
}

func (be *osBackend) Getenv(key string) string {
	return os.Getenv(key)
}

var _ PathLooker = &osBackend{}
