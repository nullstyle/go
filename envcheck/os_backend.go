package envcheck

import (
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

var _ PathLooker = &osBackend{}
