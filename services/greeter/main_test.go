package main

import (
	"testing"

	"github.com/nullstyle/go/protocols/example"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestServer(t *testing.T) {
	srv := &greeterServer{}

	resp, err := srv.Greet(
		context.Background(),
		&example.GreetRequest{Name: "Scott"},
	)

	if assert.NoError(t, err) {
		assert.Equal(t, "Hello, Scott!", resp.Msg)
	}
}
