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
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
)

const (
	configFileName = ".env"

	dbUsername = "DB_USERNAME"
	dbPassword = "DB_PASSWORD"
	dbHost     = "DB_HOSTNAME"
	dbName     = "DB_NAME"
)

var (
	dbService *service.DbService
	err       error
)

func init() {

	viper.SetConfigFile(configFileName)
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}

	dbService, err = service.NewDbService(&service.DbConnCfg{
		DbUsername: viper.GetString(dbUsername),
		DbPassword: viper.GetString(dbPassword),
		DbHost:     viper.GetString(dbHost),
		DbName:     viper.GetString(dbName),
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

	dbService.SelectAll()
}
