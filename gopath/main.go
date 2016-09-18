// Package gopath implements functions for working with a GOPATH value
package gopath

import (
	"os"
	"strings"
)

// Current returns the current gopath
func Current() string {
	return os.Getenv("GOPATH")
}

// IsCurrentSingular returns true if current gopath is a single path
func IsCurrentSingular() bool {
	return IsSingular(Current())
}

// IsSingular returns true if the provided path represents a single path
func IsSingular(path string) bool {
	return !strings.Contains(path, ":")
}
