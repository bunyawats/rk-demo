package main

import (
	"context"
	"fmt"
	greeter "github.com/rookie-ninja/rk-demo/api/gen/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
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

	doCustomerUnary(c)
}

func doCustomerUnary(c greeter.CustomerClient) {

	fmt.Println("Start to do a ReadAll RPC...")

	res, err := c.ReadAll(context.Background(), &greeter.ReadAllRequest{})
	if err != nil {
		log.Fatalf("error while calling Greet RPC: %v\n", err)
	}
	fmt.Printf("Response from Greet: %v\n", res.CustomerList)

}
