//+build js

package main

import (
	"testing"

	"github.com/nullstyle/go/gopherjs/mithril/mtest"
	"github.com/stretchr/testify/assert"
)

func TestApp(t *testing.T) {
	mt := mtest.New(t, app)
	assert.True(t, mt.Output.Contains("hello world"))
}
