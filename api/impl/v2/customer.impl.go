package v2

import (
	"context"
	greeterV2 "github.com/rookie-ninja/rk-demo/api/gen/v2"
	"github.com/rookie-ninja/rk-demo/service"
	"log"
	"time"
)

type CustomerServer struct {
	context   context.Context
	dbService *service.SQLcService
}

var (
	staticCustomerList = []*greeterV2.CustomerModel{
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

func NewCustomerServer(cx context.Context, dbs *service.SQLcService) *CustomerServer {
	return &CustomerServer{
		context:   cx,
		dbService: dbs,
	}
}

func (server *CustomerServer) ReadAll(
	_ context.Context, _ *greeterV2.ReadAllRequest) (*greeterV2.ReadAllResponse, error) {

	customerList, err := server.dbService.SelectAll()
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

	//cusId, err := server.dbService.InsertNewCustomer(&service.CustomerRecord{
	//	Fname: request.FirstName,
	//	Lname: request.LastName,
	//	Age:   int(request.Age),
	//})
	//if err != nil {
	//	log.Println(err.Error())
	//	return nil, err
	//}
	cusId := server.dbService.InsertNewCustomer()

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
	//rowsAffected, err := server.dbService.UpdateCustomer(&service.CustomerRecord{
	//	CusId: reqCus.CusId,
	//	Fname: reqCus.FirstName,
	//	Lname: reqCus.LastName,
	//	Age:   int(reqCus.Age),
	//})
	//if err != nil {
	//	log.Println(err.Error())
	//	return nil, err
	//}
	server.dbService.UpdateCustomer(reqCus.CusId)

	return &greeterV2.UpdateResponse{
		UpdatedCount: reqCus.CusId,
	}, nil
}

func (server *CustomerServer) Delete(
	_ context.Context, request *greeterV2.DeleteRequest) (*greeterV2.DeleteResponse, error) {

	//rowsAffected, err := server.dbService.DeleteCustomer(request.CusId)
	//if err != nil {
	//	log.Println(err.Error())
	//	return nil, err
	//}
	server.dbService.DeleteCustomer(request.CusId)

	return &greeterV2.DeleteResponse{
		DeletedCount: request.CusId,
	}, nil
}

func (server *CustomerServer) ReadAllStream(
	_ *greeterV2.ReadAllRequest, resStream greeterV2.Customer_ReadAllStreamServer) error {

	customerList, err := server.dbService.SelectAll()
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
		resStream.Send(res)
		time.Sleep(time.Second * 2)
	}

	return nil
}
