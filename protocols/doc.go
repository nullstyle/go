// Package protocols contains subpackages that implement grpc protocols.
package protocols

//go:generate bash -c "protoc --go_out=plugins=grpc:. ./example/*.proto"
//go:generate pwd
