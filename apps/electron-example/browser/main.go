package main

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/nullstyle/go/gopherjs/mithril"
)

//go:generate gopm build

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

func main() {
	mithril.Mount(
		js.Global.Get("document").Get("body"),
		app,
	)
}
