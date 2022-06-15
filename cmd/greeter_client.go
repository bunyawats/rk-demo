package main

import (
	"context"
	"fmt"
	greeterV1 "github.com/rookie-ninja/rk-demo/api/gen/v1"
	greeterV2 "github.com/rookie-ninja/rk-demo/api/gen/v2"
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

	c1 := greeterV1.NewGreeterClient(cc)
	c2 := greeterV2.NewGreeterClient(cc)
	//fmt.Printf("Create client %f\n\n", c)

	doGreeterUnaryV1(c1)
	doGreeterUnaryV2(c2)
}

func doGreeterUnaryV1(c greeterV1.GreeterClient) {

	fmt.Println("Start to do a Unary RPC...")
	req := &greeterV1.HelloRequest{
		Name: "Bunyawat",
	}
	res, err := c.Hello(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Greet RPC: %v\n", err)
	}
	fmt.Printf("Response from Greet: %v\n", res.MyMessage)

}

func doGreeterUnaryV2(c greeterV2.GreeterClient) {

	fmt.Println("Start to do a Unary RPC...")
	req := &greeterV2.HelloRequest{
		Name: "Bunyawat",
	}
	res, err := c.Hello(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Greet RPC: %v\n", err)
	}
	fmt.Printf("Response from Greet: %v\n", res.MyMessage)

}
