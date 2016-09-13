// Package produce implements a website building toolkit for products.  By
// creating a go application that uses this package you can quickly create an
// static site generator, a live-reloading development server and a acme-enabled
// production quality server and more.
//
package produce

type Product struct {
	Name string
}

func New(name string) *Product {
	return &Product{
		Name: name,
	}
}
