package v1

import (
	"context"
	greeterV1 "github.com/rookie-ninja/rk-demo/api/gen/v1"
	"github.com/rookie-ninja/rk-demo/service"
	"log"
)

type CustomerServer struct {
	context   context.Context
	dbService *service.DbService
}

var (
	_ = []*greeterV1.CustomerModel{
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
	_ context.Context, _ *greeterV1.ReadAllRequest) (*greeterV1.ReadAllResponse, error) {

	customerList, err := server.dbService.SelectAll()
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	cusListFromDB := make([]*greeterV1.CustomerModel, 0)
	for _, cus := range customerList {
		cusListFromDB = append(cusListFromDB, &greeterV1.CustomerModel{
			FirstName: cus.Fname,
			LastName:  cus.Lname,
			Age:       int32(cus.Age),
		})
	}

	return &greeterV1.ReadAllResponse{
		CustomerList: cusListFromDB,
	}, nil
}

func (server *CustomerServer) Create(
	_ context.Context, request *greeterV1.CreateRequest) (*greeterV1.CreateResponse, error) {

	cusId, err := server.dbService.InsertNewCustomer(&service.CustomerRecord{
		Fname: request.FirstName,
		Lname: request.LastName,
		Age:   int(request.Age),
	})
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return &greeterV1.CreateResponse{
		Customer: &greeterV1.CustomerModel{
			CusId:     cusId,
			FirstName: request.FirstName,
			LastName:  request.LastName,
			Age:       request.Age,
		},
	}, nil
}

func (server *CustomerServer) Update(
	_ context.Context, request *greeterV1.UpdateRequest) (*greeterV1.UpdateResponse, error) {

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
	return &greeterV1.UpdateResponse{
		UpdatedCount: rowsAffected,
	}, nil
}

func (server *CustomerServer) Delete(
	_ context.Context, request *greeterV1.DeleteRequest) (*greeterV1.DeleteResponse, error) {

	rowsAffected, err := server.dbService.DeleteCustomer(request.CusId)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return &greeterV1.DeleteResponse{
		DeletedCount: rowsAffected,
	}, nil
}
