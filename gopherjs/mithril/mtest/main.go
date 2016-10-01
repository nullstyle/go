// Package mtest implements functions for testing mithril components.
package mtest

import (
	"testing"

	"github.com/gopherjs/gopherjs/js"
	"github.com/nullstyle/go/gopherjs/module"
)

var mq *js.Object

// MithrilTest provides helper functions for testing gopherjs-based mithril
// applications
type MithrilTest struct {
	T         *testing.T
	Component interface{}
	Output    *Output
}

type Output struct {
	*js.Object
}

func init() {
	mq = module.Require("mithril-query")
}

// New starts a new mithril test using the provided mithril component.
func New(t *testing.T, comp interface{}) *MithrilTest {
	result := &MithrilTest{
		T:         t,
		Component: comp,
		Output:    Query(comp),
	}

	return result
}

// Query invokes mithril-query on the provided component or view
func Query(compOrView interface{}) *Output {
	return &Output{
		Object: mq.Invoke(compOrView),
	}
}
