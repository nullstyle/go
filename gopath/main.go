// Package gopath implements functions for working with a GOPATH value
package gopath

import "strings"

// Sep is the GOPATH element separator
const Sep = ":"

// First returns the first GOPATH component
func First(gpath string) string {
	parts := Split(gpath)
	if len(parts) == 0 {
		return ""
	}

	return parts[0]
}

// IsSingular returns true if the provided path represents a single path
func IsSingular(path string) bool {
	return !strings.Contains(path, ":")
}

// Split returns the components of a GOPATH value
func Split(gpath string) []string {
	return strings.Split(gpath, Sep)
}
