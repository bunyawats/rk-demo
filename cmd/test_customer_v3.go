package main

import (
	"context"
	"fmt"
	greeterV3 "github.com/rookie-ninja/rk-demo/api/gen/v3"
	"google.golang.org/grpc"
	"io"
	"log"
	"math/rand"
)

func TestCustomerV3(cc *grpc.ClientConn) {

	c := greeterV3.NewCustomerClient(cc)
	cusId := createNewCustomerV3(c)
	getAllCustomerV3(c)
	updateExistCustomerV3(c, cusId)
	getAllCustomerV3(c)
	deleteExistCustomerV3(c, cusId)
	getAllCustomerStreamV3(c)
}

func deleteExistCustomerV3(c greeterV3.CustomerClient, id string) {
	res, err := c.Delete(context.Background(), &greeterV3.DeleteRequest{
		CusId: id,
	})
	if err != nil {
		log.Fatalf("error while calling Customer.Delete RPC: %v\n", err)
	}
	fmt.Printf("Response from Customer.Delete: %v\n", res.DeletedCount)
}

func updateExistCustomerV3(c greeterV3.CustomerClient, id string) {
	ranNumber := int32(rand.Intn(100))

	res, err := c.Update(context.Background(), &greeterV3.UpdateRequest{
		Customer: &greeterV3.CustomerModel{
			CusId:     id,
			FirstName: fmt.Sprintf("%v_%v", "Bunyawat", ranNumber),
			LastName:  fmt.Sprintf("%v_%v", "Singchai", ranNumber),
			Age:       ranNumber,
		},
	})
	if err != nil {
		log.Fatalf("error while calling Customer.Update RPC: %v\n", err)
	}
	fmt.Printf("Response from Customer.Update V2: %v\n", res.UpdatedCount)
}

func getAllCustomerV3(c greeterV3.CustomerClient) {
	res, err := c.ReadAll(context.Background(), &greeterV3.ReadAllRequest{})
	if err != nil {
		log.Fatalf("error while calling Customer.ReadAll RPC: %v\n", err)
	}
	fmt.Printf("Response from Customer.ReadAll V2: %v\n", res.CustomerList)
}

func getAllCustomerStreamV3(c greeterV3.CustomerClient) {
	resStream, err := c.ReadAllStream(context.Background(), &greeterV3.ReadAllRequest{})
	if err != nil {
		log.Fatalf("error while calling Customer.ReadAll RPC: %v\n", err)
	}

	for {
		cus, err := resStream.Recv()
		if err == io.EOF {
			// we have reached the end of the straem
			break
		}
		if err != nil {
			log.Fatalf("error while reading stream %v\n", err)
		}
		fmt.Printf("Response from Customer.getAllCustomerStream V2: %v\n", cus)
	}
}

func createNewCustomerV3(c greeterV3.CustomerClient) string {
	ranNumber := int32(rand.Intn(100))

	res, err := c.Create(context.Background(), &greeterV3.CreateRequest{
		FirstName: fmt.Sprintf("%v_%v", "Bunyawat", ranNumber),
		LastName:  fmt.Sprintf("%v_%v", "Singchai", ranNumber),
		Age:       ranNumber,
	})
	if err != nil {
		log.Fatalf("error while calling Customer.Create RPC: %v\n", err)
	}
	fmt.Printf("Response from Customer.Create V2: %v\n", res.Customer)

	return res.Customer.CusId
}
