// Copyright (c) 2021 rookie-ninja
//
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"github.com/rookie-ninja/rk-boot/v2"
	"github.com/rookie-ninja/rk-demo/api/gen/v1"
	"github.com/rookie-ninja/rk-demo/api/impl/v1"
	"github.com/rookie-ninja/rk-demo/service"
	"github.com/rookie-ninja/rk-grpc/v2/boot"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	dbService *service.DbService
	err       error

	db_username = "test"
	db_password = "test"
	db_host     = "127.0.0.1:3306"
	db_name     = "test"
)

func init() {

	dbService, err = service.NewDbService(&service.DbConnCfg{
		DbUsername: db_username,
		DbPassword: db_password,
		DbHost:     db_host,
		DbName:     db_name,
	})
	if err != nil {
		panic(error.Error)
	}
}

func main() {

	// defer the close till after the main function has finished
	// executing
	defer dbService.Close()

	boot := rkboot.NewBoot()

	// register grpc
	entry := rkgrpc.GetGrpcEntry("rk-demo")
	entry.AddRegFuncGrpc(registerGreeter)
	entry.AddRegFuncGw(greeter.RegisterGreeterHandlerFromEndpoint)

	// Bootstrap
	boot.Bootstrap(context.TODO())

	// Wait for shutdown sig
	boot.WaitForShutdownSig(context.TODO())
}

func registerGreeter(server *grpc.Server) {
	greeter.RegisterGreeterServer(server, &v1.GreeterServer{})
	greeter.RegisterCustomerServer(server, v1.NewCustomerServer(context.TODO(), dbService))
	reflection.Register(server)
}
