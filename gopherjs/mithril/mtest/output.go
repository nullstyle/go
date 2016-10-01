package mtest

import "github.com/gopherjs/gopherjs/js"

func (o *Output) Contains(str string) bool {
	return o.Call("contains", str).Bool()
}

// Has returns true if any element in tree matches the selector, otherwise false
func (o *Output) Has(sel string) bool {
	return o.Call("has", sel).Bool()
}

func (o *Output) Find(sel string) *js.Object {
	return o.Call("find", sel)
}

func (o *Output) First(sel string) *js.Object {
	return o.Call("first", sel)
}
