package main

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"math/rand"
	"time"
)

func main() {
	println("Call Customer Service")

	tls := true
	opts := grpc.WithTransportCredentials(insecure.NewCredentials())
	if tls {
		certFile := "./certs/ca.pem" // Certificate Authority Trust certificate
		creds, sslErr := credentials.NewClientTLSFromFile(certFile, "")
		if sslErr != nil {
			log.Fatalf("Error while loading CA trust certificate: %v", sslErr)
			return
		}
		opts = grpc.WithTransportCredentials(creds)
	}

	cc, err := grpc.Dial("localhost:8080", opts)
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()

	rand.Seed(time.Now().UnixNano())

	TestCustomerV1(cc)
	TestCustomerV2(cc)
	TestCustomerV3(cc)

}
