// Package greeter implements a client for accessing a greeter grpc service.
package greeter

import (
	"github.com/nullstyle/go/protocols/example"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

// Client represents a grpc client connected to a greeter server.
type Client struct {
	example.GreeterClient
	conn *grpc.ClientConn
}

// New creates a new greeter client, connected to the server at addr
func New(addr string, opts ...grpc.DialOption) (*Client, error) {
	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		return nil, errors.Wrap(err, "grpc dial failed")
	}

	cl := example.NewGreeterClient(conn)

	return &Client{
		GreeterClient: cl,
		conn:          conn,
	}, nil
}
