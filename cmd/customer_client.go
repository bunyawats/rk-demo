package main

import (
	"context"
	"fmt"
	greeterV1 "github.com/rookie-ninja/rk-demo/api/gen/v1"
	greeterV2 "github.com/rookie-ninja/rk-demo/api/gen/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
	"math/rand"
	"time"
)

func main() {
	println("Call Customer Service")

	opts := grpc.WithTransportCredentials(insecure.NewCredentials())

	cc, err := grpc.Dial("localhost:8080", opts)
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()

	rand.Seed(time.Now().UnixNano())

	c1 := greeterV1.NewCustomerClient(cc)
	cusId := createNewCustomer(c1)
	cusId = createNewCustomer(c1)
	getAllCustomer(c1)
	updateExistCustomer(c1, cusId)
	getAllCustomer(c1)
	deleteExistCustomer(c1, cusId-1)
	getAllCustomer(c1)

	c2 := greeterV2.NewCustomerClient(cc)
	cusId = createNewCustomerV2(c2)
	cusId = createNewCustomerV2(c2)
	getAllCustomerV2(c2)
	updateExistCustomerV2(c2, cusId)
	getAllCustomerV2(c2)
	deleteExistCustomerV2(c2, cusId-1)
	getAllCustomerV2(c2)
	getAllCustomerStreamV2(c2)
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
