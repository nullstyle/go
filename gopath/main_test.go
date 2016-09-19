package gopath

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFirst(t *testing.T) {
	cases := []struct {
		Name     string
		GOPATH   string
		Expected string
	}{
		{"empty", "", ""},
		{"single", "/home/nullstyle/go", "/home/nullstyle/go"},
		{"multi", "/go:/go2", "/go"},
	}

	for _, kase := range cases {
		t.Run(kase.Name, func(t *testing.T) {
			actual := First(kase.GOPATH)
			assert.Equal(t, kase.Expected, actual)
		})
	}
}

func TestIsSingular(t *testing.T) {
	cases := []struct {
		Name     string
		GOPATH   string
		Expected bool
	}{
		{"empty", "", true},
		{"single", "/home/nullstyle/go", true},
		{"multi", "/go:/go2", false},
	}

	for _, kase := range cases {
		t.Run(kase.Name, func(t *testing.T) {
			actual := IsSingular(kase.GOPATH)
			assert.Equal(t, kase.Expected, actual)
		})
	}
}

func TestSplit(t *testing.T) {
	cases := []struct {
		Name     string
		GOPATH   string
		Expected []string
	}{
		{"empty", "", []string{""}},
		{"single", "/home/nullstyle/go", []string{"/home/nullstyle/go"}},
		{"multi", "/go:/go2", []string{"/go", "/go2"}},
	}

	for _, kase := range cases {
		t.Run(kase.Name, func(t *testing.T) {
			actual := Split(kase.GOPATH)
			assert.Equal(t, kase.Expected, actual)
		})
	}
}
