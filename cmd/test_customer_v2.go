package main

import (
	"context"
	"fmt"
	greeterV2 "github.com/rookie-ninja/rk-demo/api/gen/v2"
	"google.golang.org/grpc"
	"io"
	"log"
	"math/rand"
)

func TestCustomerV2(cc *grpc.ClientConn) {

	c := greeterV2.NewCustomerClient(cc)
	cusId := createNewCustomerV2(c)
	getAllCustomerV2(c)
	updateExistCustomerV2(c, cusId)
	getAllCustomerV2(c)
	deleteExistCustomerV2(c, cusId)
	getAllCustomerStreamV2(c)
}

func deleteExistCustomerV2(c greeterV2.CustomerClient, id int32) {
	res, err := c.Delete(context.Background(), &greeterV2.DeleteRequest{
		CusId: id,
	})
	if err != nil {
		log.Fatalf("error while calling Customer.Delete RPC: %v\n", err)
	}
	fmt.Printf("Response from Customer.Delete: %v\n", res.DeletedCount)
}

func updateExistCustomerV2(c greeterV2.CustomerClient, id int32) {
	ranNumber := int32(rand.Intn(100))

	res, err := c.Update(context.Background(), &greeterV2.UpdateRequest{
		Customer: &greeterV2.CustomerModel{
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

func getAllCustomerV2(c greeterV2.CustomerClient) {
	res, err := c.ReadAll(context.Background(), &greeterV2.ReadAllRequest{})
	if err != nil {
		log.Fatalf("error while calling Customer.ReadAll RPC: %v\n", err)
	}
	fmt.Printf("Response from Customer.ReadAll V2: %v\n", res.CustomerList)
}

func getAllCustomerStreamV2(c greeterV2.CustomerClient) {
	resStream, err := c.ReadAllStream(context.Background(), &greeterV2.ReadAllRequest{})
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

func createNewCustomerV2(c greeterV2.CustomerClient) int32 {
	ranNumber := int32(rand.Intn(100))

	res, err := c.Create(context.Background(), &greeterV2.CreateRequest{
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
