package v1

import (
	"context"
	greeterV1 "github.com/rookie-ninja/rk-demo/api/gen/v1"
)

//GreeterServer GreeterServer struct
type GreeterServer struct{}

// Hello response with hello message
func (server *GreeterServer) Hello(
	_ context.Context,
	request *greeterV1.HelloRequest) (*greeterV1.HelloResponse, error) {

	return &greeterV1.HelloResponse{
		MyMessage: "hello! v1 " + request.GetName(),
	}, nil
}
