package v3

import (
	"context"
	greeterV3 "github.com/rookie-ninja/rk-demo/api/gen/v3"
)

//GreeterServer GreeterServer struct
type GreeterServer struct{}

// Hello response with hello message
func (server *GreeterServer) Hello(
	_ context.Context,
	request *greeterV3.HelloRequest) (*greeterV3.HelloResponse, error) {

	return &greeterV3.HelloResponse{
		MyMessage: "hello V3! " + request.GetName(),
	}, nil
}
