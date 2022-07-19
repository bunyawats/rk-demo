package v3

import (
	"context"
	greeterV3 "github.com/rookie-ninja/rk-demo/api/gen/v3"
	"github.com/rookie-ninja/rk-demo/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"time"
)

type CustomerServer struct {
	context        context.Context
	mongoDbService *service.MongoService
}

func NewCustomerServer(cx context.Context, dbs *service.MongoService) *CustomerServer {
	return &CustomerServer{
		context:        cx,
		mongoDbService: dbs,
	}
}

func (server *CustomerServer) ReadAll(
	_ context.Context, _ *greeterV3.ReadAllRequest) (*greeterV3.ReadAllResponse, error) {

	customerList, err := server.mongoDbService.ListCustomer()
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	cusListFromDB := make([]*greeterV3.CustomerModel, 0)
	for _, cus := range *customerList {
		cusListFromDB = append(cusListFromDB, &greeterV3.CustomerModel{
			FirstName: cus.Fname,
			LastName:  cus.Lname,
			Age:       int32(cus.Age),
		})
	}

	return &greeterV3.ReadAllResponse{
		CustomerList: cusListFromDB,
	}, nil
}

func (server *CustomerServer) Create(
	_ context.Context, request *greeterV3.CreateRequest) (*greeterV3.CreateResponse, error) {

	cusId, err := server.mongoDbService.CreateCustomer(&service.CustomerDoc{
		Fname: request.FirstName,
		Lname: request.LastName,
		Age:   int(request.Age),
	})
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return &greeterV3.CreateResponse{
		Customer: &greeterV3.CustomerModel{
			CusId:     cusId,
			FirstName: request.FirstName,
			LastName:  request.LastName,
			Age:       request.Age,
		},
	}, nil
}

func (server *CustomerServer) Update(
	_ context.Context, request *greeterV3.UpdateRequest) (*greeterV3.UpdateResponse, error) {

	reqCus := request.Customer
	objectId, err := primitive.ObjectIDFromHex(reqCus.CusId)

	rowsAffected, err := server.mongoDbService.UpdateCustomer(&service.CustomerDoc{
		CusId: objectId,
		Fname: reqCus.FirstName,
		Lname: reqCus.LastName,
		Age:   int(reqCus.Age),
	})
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return &greeterV3.UpdateResponse{
		UpdatedCount: int32(rowsAffected),
	}, nil
}

func (server *CustomerServer) Delete(
	_ context.Context, request *greeterV3.DeleteRequest) (*greeterV3.DeleteResponse, error) {

	rowsAffected, err := server.mongoDbService.DeleteCustomer(request.CusId)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return &greeterV3.DeleteResponse{
		DeletedCount: int32(rowsAffected),
	}, nil
}

func (server *CustomerServer) ReadAllStream(
	_ *greeterV3.ReadAllRequest, resStream greeterV3.Customer_ReadAllStreamServer) error {

	customerList, err := server.mongoDbService.ListCustomer()
	if err != nil {
		log.Println(err.Error())
		return err
	}

	for _, cus := range *customerList {
		res := &greeterV3.CustomerModel{
			FirstName: cus.Fname,
			LastName:  cus.Lname,
			Age:       int32(cus.Age),
		}
		err := resStream.Send(res)
		if err != nil {
			return err
		}
		time.Sleep(time.Second * 2)
	}

	return nil
}
