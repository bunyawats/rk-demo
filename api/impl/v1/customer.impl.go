package v1

import (
	"context"
	greeter "github.com/rookie-ninja/rk-demo/api/gen/v1"
	"github.com/rookie-ninja/rk-demo/service"
	"log"
)

type CustomerServer struct {
	context   context.Context
	dbService *service.DbService
}

var (
	staticCustomerList = []*greeter.CustomerModel{
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
	}
)

func NewCustomerServer(cx context.Context, dbs *service.DbService) *CustomerServer {
	return &CustomerServer{
		context:   cx,
		dbService: dbs,
	}
}

func (server *CustomerServer) ReadAll(
	_ context.Context, _ *greeter.ReadAllRequest) (*greeter.ReadAllResponse, error) {

	customerList, err := server.dbService.SelectAll()
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	cusListFromDB := make([]*greeter.CustomerModel, 0)
	for _, cus := range customerList {
		cusListFromDB = append(cusListFromDB, &greeter.CustomerModel{
			FirstName: cus.Fname,
			LastName:  cus.Lname,
			Age:       int32(cus.Age),
		})
	}

	return &greeter.ReadAllResponse{
		CustomerList: cusListFromDB,
	}, nil
}

func (server *CustomerServer) Create(
	_ context.Context, request *greeter.CreateRequest) (*greeter.CreateResponse, error) {

	cusId, err := server.dbService.InsertNewCustomer(&service.CustomerRecord{
		Fname: request.FirstName,
		Lname: request.LastName,
		Age:   int(request.Age),
	})
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return &greeter.CreateResponse{
		Customer: &greeter.CustomerModel{
			CusId:     cusId,
			FirstName: request.FirstName,
			LastName:  request.LastName,
			Age:       request.Age,
		},
	}, nil
}

func (server *CustomerServer) Update(
	_ context.Context, request *greeter.UpdateRequest) (*greeter.UpdateResponse, error) {

	reqCus := request.Customer
	rowsAffected, err := server.dbService.UpdateCustomer(&service.CustomerRecord{
		CusId: reqCus.CusId,
		Fname: reqCus.FirstName,
		Lname: reqCus.LastName,
		Age:   int(reqCus.Age),
	})
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return &greeter.UpdateResponse{
		UpdatedCount: rowsAffected,
	}, nil
}

func (server *CustomerServer) Delete(
	_ context.Context, request *greeter.DeleteRequest) (*greeter.DeleteResponse, error) {

	rowsAffected, err := server.dbService.DeleteCustomer(request.CusId)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return &greeter.DeleteResponse{
		DeletedCount: rowsAffected,
	}, nil
}
