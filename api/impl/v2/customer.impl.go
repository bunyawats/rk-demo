package v2

import (
	"context"
	greeterV2 "github.com/rookie-ninja/rk-demo/api/gen/v2"
	"github.com/rookie-ninja/rk-demo/service"
	"log"
	"time"
)

type CustomerServer struct {
	context     context.Context
	sqlcService *service.SQLcService
}

func NewCustomerServer(cx context.Context, dbs *service.SQLcService) *CustomerServer {
	return &CustomerServer{
		context:     cx,
		sqlcService: dbs,
	}
}

func (server *CustomerServer) ReadAll(
	_ context.Context, _ *greeterV2.ReadAllRequest) (*greeterV2.ReadAllResponse, error) {

	customerList, err := server.sqlcService.SelectAll()
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	cusListFromDB := make([]*greeterV2.CustomerModel, 0)
	for _, cus := range customerList {
		cusListFromDB = append(cusListFromDB, &greeterV2.CustomerModel{
			FirstName: cus.Fname.String,
			LastName:  cus.Lname.String,
			Age:       cus.Age.Int32,
		})
	}

	return &greeterV2.ReadAllResponse{
		CustomerList: cusListFromDB,
	}, nil
}

func (server *CustomerServer) Create(
	_ context.Context, request *greeterV2.CreateRequest) (*greeterV2.CreateResponse, error) {

	cusId, err := server.sqlcService.InsertNewCustomer(&service.CustomerRecord{
		Fname: request.FirstName,
		Lname: request.LastName,
		Age:   int(request.Age),
	})
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return &greeterV2.CreateResponse{
		Customer: &greeterV2.CustomerModel{
			CusId:     cusId,
			FirstName: request.FirstName,
			LastName:  request.LastName,
			Age:       request.Age,
		},
	}, nil
}

func (server *CustomerServer) Update(
	_ context.Context, request *greeterV2.UpdateRequest) (*greeterV2.UpdateResponse, error) {

	reqCus := request.Customer
	rowsAffected, err := server.sqlcService.UpdateCustomer(&service.CustomerRecord{
		CusId: reqCus.CusId,
		Fname: reqCus.FirstName,
		Lname: reqCus.LastName,
		Age:   int(reqCus.Age),
	})
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return &greeterV2.UpdateResponse{
		UpdatedCount: rowsAffected,
	}, nil
}

func (server *CustomerServer) Delete(
	_ context.Context, request *greeterV2.DeleteRequest) (*greeterV2.DeleteResponse, error) {

	rowsAffected, err := server.sqlcService.DeleteCustomer(request.CusId)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return &greeterV2.DeleteResponse{
		DeletedCount: rowsAffected,
	}, nil
}

func (server *CustomerServer) ReadAllStream(
	_ *greeterV2.ReadAllRequest, resStream greeterV2.Customer_ReadAllStreamServer) error {

	customerList, err := server.sqlcService.SelectAll()
	if err != nil {
		log.Println(err.Error())
		return err
	}

	for _, cus := range customerList {
		res := &greeterV2.CustomerModel{
			FirstName: cus.Fname.String,
			LastName:  cus.Lname.String,
			Age:       cus.Age.Int32,
		}
		err := resStream.Send(res)
		if err != nil {
			return err
		}
		time.Sleep(time.Second * 2)
	}

	return nil
}
