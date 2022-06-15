// Copyright (c) 2021 rookie-ninja
//
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"database/sql"
	"github.com/rookie-ninja/rk-boot/v2"
	greeterV1 "github.com/rookie-ninja/rk-demo/api/gen/v1"
	greeterV2 "github.com/rookie-ninja/rk-demo/api/gen/v2"
	"github.com/rookie-ninja/rk-demo/api/impl/v1"
	"github.com/rookie-ninja/rk-demo/api/impl/v2"
	"github.com/rookie-ninja/rk-demo/repository"
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
	db          *sql.DB
	dbService   *service.DbService
	sqlcService *service.SQLcService
	err         error
)

func init() {

	viper.SetConfigFile(configFileName)
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}

	db = repository.NewDbConnection(&repository.DbConnCfg{
		DbUsername: viper.GetString(dbUsername),
		DbPassword: viper.GetString(dbPassword),
		DbHost:     viper.GetString(dbHost),
		DbName:     viper.GetString(dbName),
	})

	dbService = &service.DbService{DB: db}
	sqlcService = service.NewSQLcService(db, context.TODO())
}

func main() {

	// defer the close till after the main function has finished
	// executing
	defer db.Close()

	boot := rkboot.NewBoot()

	// register grpc
	entry := rkgrpc.GetGrpcEntry("ssc-poc")
	entry.AddRegFuncGrpc(registerGreeter)
	entry.AddRegFuncGw(greeterV1.RegisterGreeterHandlerFromEndpoint)
	entry.AddRegFuncGw(greeterV2.RegisterGreeterHandlerFromEndpoint)

	// Bootstrap
	boot.Bootstrap(context.TODO())

	// Wait for shutdown sig
	boot.WaitForShutdownSig(context.TODO())
}

func registerGreeter(server *grpc.Server) {
	greeterV1.RegisterGreeterServer(server, &v1.GreeterServer{})
	greeterV1.RegisterCustomerServer(server, v1.NewCustomerServer(context.TODO(), dbService))

	greeterV2.RegisterGreeterServer(server, &v2.GreeterServer{})
	greeterV2.RegisterCustomerServer(server, v2.NewCustomerServer(context.TODO(), sqlcService))

	reflection.Register(server)

	testService()
}

func testService() {
	dbService.SelectAll()
	sqlcService.SelectAll()
}
