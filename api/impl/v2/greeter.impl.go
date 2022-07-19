package v2

import (
	"context"
	greeterV2 "github.com/rookie-ninja/rk-demo/api/gen/v2"
)

//GreeterServer GreeterServer struct
type GreeterServer struct{}

// Hello response with hello message
func (server *GreeterServer) Hello(
	_ context.Context,
	request *greeterV2.HelloRequest) (*greeterV2.HelloResponse, error) {

	return &greeterV2.HelloResponse{
		MyMessage: "hello V2! " + request.GetName(),
	}, nil
}
