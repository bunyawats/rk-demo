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
	println("Call to GRPC server")

	opts := grpc.WithTransportCredentials(insecure.NewCredentials())

	cc, err := grpc.Dial("localhost:8080", opts)
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()

	c := greeter.NewGreeterClient(cc)
	//fmt.Printf("Create client %f\n\n", c)

	doGreeterUnary(c)
}

func doGreeterUnary(c greeter.GreeterClient) {

	fmt.Println("Start to do a Unary RPC...")
	req := &greeter.HelloRequest{
		Name: "Bunyawat",
	}
	res, err := c.Hello(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Greet RPC: %v\n", err)
	}
	fmt.Printf("Response from Greet: %v\n", res.MyMessage)

}
