package main

import (
	"context"
	"fmt"
	greeter "github.com/rookie-ninja/rk-demo/api/gen/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

	c := greeter.NewCustomerClient(cc)
	//fmt.Printf("Create client %f\n\n", c)

	rand.Seed(time.Now().UnixNano())
	cusId := createNewCustomer(c)
	cusId = createNewCustomer(c)
	getAllCustomer(c)
	updateExistCustomer(c, cusId)
	getAllCustomer(c)
	deleteExistCustomer(c, cusId-1)
	getAllCustomer(c)
}

func deleteExistCustomer(c greeter.CustomerClient, id int32) {
	res, err := c.Delete(context.Background(), &greeter.DeleteRequest{
		CusId: id,
	})
	if err != nil {
		log.Fatalf("error while calling Customer.Delete RPC: %v\n", err)
	}
	fmt.Printf("Response from Customer.Delete: %v\n", res.DeletedCount)
}

func updateExistCustomer(c greeter.CustomerClient, id int32) {

	ranNumber := int32(rand.Intn(100))

	res, err := c.Update(context.Background(), &greeter.UpdateRequest{
		Customer: &greeter.CustomerModel{
			CusId:     id,
			FirstName: fmt.Sprintf("%v_%v", "Bunyawat_", ranNumber),
			LastName:  fmt.Sprintf("%v_%v", "Singchai", ranNumber),
			Age:       ranNumber,
		},
	})
	if err != nil {
		log.Fatalf("error while calling Customer.Update RPC: %v\n", err)
	}
	fmt.Printf("Response from Customer.Update: %v\n", res.UpdatedCount)
}

func createNewCustomer(c greeter.CustomerClient) int32 {

	res, err := c.Create(context.Background(), &greeter.CreateRequest{
		FirstName: "Bunyawat 12",
		LastName:  "Singchai12",
		Age:       int32(rand.Intn(100)),
	})
	if err != nil {
		log.Fatalf("error while calling Customer.Create RPC: %v\n", err)
	}
	fmt.Printf("Response from Customer.Create: %v\n", res.Customer)

	return res.Customer.CusId
}

func getAllCustomer(c greeter.CustomerClient) {

	res, err := c.ReadAll(context.Background(), &greeter.ReadAllRequest{})
	if err != nil {
		log.Fatalf("error while calling Customer.ReadAll RPC: %v\n", err)
	}
	fmt.Printf("Response from Customer.ReadAll: %v\n", res.CustomerList)

}
