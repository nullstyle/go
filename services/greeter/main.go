package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/nullstyle/go/protocols/example"
	"golang.org/x/net/context"
)

var (
	tls      = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	certFile = flag.String("cert_file", "", "The TLS cert file")
	keyFile  = flag.String("key_file", "", "The TLS key file")
	port     = flag.Int("port", 10000, "The server port")
)

type greeterServer struct {
}

func (s *greeterServer) Greet(
	ctx context.Context,
	req *example.GreetRequest,
) (*example.GreetResponse, error) {

	return &example.GreetResponse{
		Msg: fmt.Sprintf("Hello, %s!", req.Name),
	}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	if *tls {
		creds, err := credentials.NewServerTLSFromFile(*certFile, *keyFile)
		if err != nil {
			log.Fatalf("Failed to generate credentials %v", err)
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}
	grpcServer := grpc.NewServer(opts...)
	example.RegisterGreeterServer(grpcServer, &greeterServer{})
	grpcServer.Serve(lis)
}

var _ example.GreeterServer = &greeterServer{}
