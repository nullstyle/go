// Package produce implements the developer site for the produce project.
package main

import (
	"github.com/nullstyle/go/produce"
	// "github.com/nullstyle/go/produce/plugins/oss"
)

var product = produce.New("oss-produce")

func main() {
	product.OSSProject(func(prj *produce.OSSProject) error {
		prj.Github("nullstyle/go/produce")
	})
	product.Run()
}
