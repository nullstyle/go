package cmd

import (
	"os/exec"
	"strings"

	"github.com/pkg/errors"
)

func which(program string) (string, error) {
	raw, err := exec.Command("which", program).Output()
	if err != nil {
		return "", errors.Wrap(err, "`which` command failed")
	}

	return strings.TrimSpace(string(raw)), nil
}
