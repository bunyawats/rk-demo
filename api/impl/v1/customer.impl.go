package v1

import (
	"context"
	greeter "github.com/rookie-ninja/rk-demo/api/gen/v1"
)

type CustomerServer struct{}

func (server *CustomerServer) Create(
	_ context.Context, request *greeter.CreateRequest) (*greeter.CreateResponse, error) {

	return nil, nil
}

func (server *CustomerServer) ReadAll(
	_ context.Context, request *greeter.ReadAllRequest) (*greeter.ReadAllResponse, error) {

	res := &greeter.ReadAllResponse{
		CustomerList: []*greeter.CustomerModel{
			{
				CusId:     1,
				FirstName: "Bunyawat",
				LastName:  "Singchai",
				Age:       51,
			},
			{
				CusId:     2,
				FirstName: "Waraporn",
				LastName:  "Singchai",
				Age:       44,
			},
		},
	}

	return res, nil
}

func (server *CustomerServer) Update(
	_ context.Context, request *greeter.UpdateRequest) (*greeter.UpdateResponse, error) {

	return nil, nil
}

func (server *CustomerServer) Delete(
	_ context.Context, request *greeter.DeleteRequest) (*greeter.DeleteResponse, error) {

	return nil, nil
}
