package impl

import (
	"context"
	greeter "github.com/rookie-ninja/rk-demo/api/gen/v1"
)

//GreeterServer GreeterServer struct
type GreeterServer struct{}

// Hello response with hello message
func (server *GreeterServer) Hello(
	_ context.Context,
	request *greeter.HelloRequest) (*greeter.HelloResponse, error) {

	return &greeter.HelloResponse{
		MyMessage: "hello! " + request.GetName(),
	}, nil
}
