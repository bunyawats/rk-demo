package main

import (
	"context"
	"fmt"
	greeterV1 "github.com/rookie-ninja/rk-demo/api/gen/v1"
	"google.golang.org/grpc"
	"log"
	"math/rand"
)

func TestCustomerV1(cc *grpc.ClientConn) {
	c := greeterV1.NewCustomerClient(cc)
	cusId := createNewCustomer(c)
	getAllCustomer(c)
	updateExistCustomer(c, cusId)
	getAllCustomer(c)
	//deleteExistCustomer(c, cusId)
	//getAllCustomer(c)
}

func deleteExistCustomer(c greeterV1.CustomerClient, id int32) {
	res, err := c.Delete(context.Background(), &greeterV1.DeleteRequest{
		CusId: id,
	})
	if err != nil {
		log.Fatalf("error while calling Customer.Delete RPC: %v\n", err)
	}
	fmt.Printf("Response from Customer.Delete V1: %v\n", res.DeletedCount)
}

func updateExistCustomer(c greeterV1.CustomerClient, id int32) {

	ranNumber := int32(rand.Intn(100))

	res, err := c.Update(context.Background(), &greeterV1.UpdateRequest{
		Customer: &greeterV1.CustomerModel{
			CusId:     id,
			FirstName: fmt.Sprintf("%v_%v", "Bunyawat", ranNumber),
			LastName:  fmt.Sprintf("%v_%v", "Singchai", ranNumber),
			Age:       ranNumber,
		},
	})
	if err != nil {
		log.Fatalf("error while calling Customer.Update RPC: %v\n", err)
	}
	fmt.Printf("Response from Customer.Update V1: %v\n", res.UpdatedCount)
}

func createNewCustomer(c greeterV1.CustomerClient) int32 {

	ranNumber := int32(rand.Intn(100))

	res, err := c.Create(context.Background(), &greeterV1.CreateRequest{
		FirstName: fmt.Sprintf("%v_%v", "Bunyawat", ranNumber),
		LastName:  fmt.Sprintf("%v_%v", "Singchai", ranNumber),
		Age:       ranNumber,
	})
	if err != nil {
		log.Fatalf("error while calling Customer.Create RPC: %v\n", err)
	}
	fmt.Printf("Response from Customer.Create V1: %v\n", res.Customer)

	return res.Customer.CusId
}

func getAllCustomer(c greeterV1.CustomerClient) {

	res, err := c.ReadAll(context.Background(), &greeterV1.ReadAllRequest{})
	if err != nil {
		log.Fatalf("error while calling Customer.ReadAll RPC: %v\n", err)
	}
	fmt.Printf("Response from Customer.ReadAll V1: %v\n", res.CustomerList)

}
