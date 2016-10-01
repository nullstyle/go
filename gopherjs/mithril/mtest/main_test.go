package mtest

import (
	"testing"

	"github.com/gopherjs/gopherjs/js"
	"github.com/nullstyle/go/gopherjs/mithril"
	"github.com/stretchr/testify/assert"
)

var app = js.M{
	"controller": js.MakeFunc(ctrl),
	"view":       js.MakeFunc(view),
}

func ctrl(this *js.Object, args []*js.Object) interface{} {
	return js.M{}
}

func view(this *js.Object, args []*js.Object) interface{} {
	return mithril.M("h1", js.M{}, "hello world")
}

func TestQuery(t *testing.T) {
	out := Query(app)

	if assert.True(t, out.Has("h1"), "couldn't find header") {
		assert.Contains(t, out.First("h1").Get("children").String(), "hello world")
		assert.True(t, out.Contains("hello world"))
	}

	if assert.False(t, out.Has("div.missing")) {
		assert.Panics(t, func() {
			out.First("div.missing")
		})
	}
}
