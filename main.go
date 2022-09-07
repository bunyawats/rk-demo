// Copyright (c) 2021 rookie-ninja
//
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/rookie-ninja/rk-boot/v2"
	greeterV1 "github.com/rookie-ninja/rk-demo/api/gen/v1"
	greeterV2 "github.com/rookie-ninja/rk-demo/api/gen/v2"
	greeterV3 "github.com/rookie-ninja/rk-demo/api/gen/v3"
	"github.com/rookie-ninja/rk-demo/api/impl/v1"
	"github.com/rookie-ninja/rk-demo/api/impl/v2"
	"github.com/rookie-ninja/rk-demo/api/impl/v3"
	"github.com/rookie-ninja/rk-demo/repository"
	"github.com/rookie-ninja/rk-demo/service"
	"github.com/rookie-ninja/rk-grpc/v2/boot"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
	"log"
	"os"
)

var (
	mysqlDB *sql.DB
	mongoDB *mongo.Database

	dbService   *service.DbService
	sqlcService *service.SQLcService
	mongoSerice *service.MongoService

	boot *rkboot.Boot
)

func getMySqlDbConn() *sql.DB {
	return mysqlDB
}

func getMongDbConn() *mongo.Database {
	return mongoDB
}

func init() {
	err := os.Setenv("RK_MYSQL_0_USER", "bunyawat444")
	if err != nil {
		return
	}
}

func main() {

	boot = rkboot.NewBoot()

	// register grpc
	entry := rkgrpc.GetGrpcEntry("ssc-grpc")
	entry.AddRegFuncGrpc(registerGreeter)
	entry.AddRegFuncGw(greeterV1.RegisterGreeterHandlerFromEndpoint)
	entry.AddRegFuncGw(greeterV2.RegisterGreeterHandlerFromEndpoint)
	entry.AddRegFuncGw(greeterV3.RegisterGreeterHandlerFromEndpoint)

	entry2 := rkgrpc.GetGrpcEntry("ssc-app")
	entry2.AddRegFuncGrpc(registerGreeter)
	entry2.AddRegFuncGw(greeterV1.RegisterGreeterHandlerFromEndpoint)
	entry2.AddRegFuncGw(greeterV2.RegisterGreeterHandlerFromEndpoint)
	entry2.AddRegFuncGw(greeterV3.RegisterGreeterHandlerFromEndpoint)

	// Bootstrap
	boot.Bootstrap(context.TODO())
	mysqlDB = repository.NewMySqlDbConnectionRKDB()
	defer func(mysqlDB *sql.DB) {
		err := mysqlDB.Close()
		if err != nil {
			log.Fatal(err.Error())
		}
	}(mysqlDB)

	mongoDB = repository.NewMongoDbConnectionRKDB()

	testService()

	// Wait for shutdown sig
	boot.WaitForShutdownSig(context.TODO())
}

func registerGreeter(server *grpc.Server) {

	dbService = service.NewDbService(getMySqlDbConn)
	sqlcService = service.NewSQLcService(getMySqlDbConn, context.TODO())
	mongoSerice = service.NewMongoService(getMongDbConn, context.TODO())

	greeterV1.RegisterGreeterServer(server, &v1.GreeterServer{})
	greeterV1.RegisterCustomerServer(server, v1.NewCustomerServer(context.TODO(), dbService))

	greeterV2.RegisterGreeterServer(server, &v2.GreeterServer{})
	greeterV2.RegisterCustomerServer(server, v2.NewCustomerServer(context.TODO(), sqlcService))

	greeterV3.RegisterGreeterServer(server, &v3.GreeterServer{})
	greeterV3.RegisterCustomerServer(server, v3.NewCustomerServer(context.TODO(), mongoSerice))

	//	reflection.Register(server)

}

func testService() {
	_, err := dbService.SelectAll()
	if err != nil {
		return
	}
	_, err = sqlcService.SelectAll()
	if err != nil {
		return
	}

	c := &service.CustomerDoc{
		Fname: "Bunyawat",
		Lname: "Singchai",
		Age:   51,
	}
	cusId, _ := mongoSerice.CreateCustomer(c)
	c = mongoSerice.GetCustomer(cusId)
	fmt.Printf("Customer id: %v is %v", cusId, c)

	customerList, _ := mongoSerice.ListCustomer()
	for ix, u := range *customerList {
		fmt.Printf("Customer id: %v is %v \n", ix, u)
	}
	c.Fname = "Waraporn"
	modifiedCount, _ := mongoSerice.UpdateCustomer(c)
	fmt.Printf("Modify count: %v ", modifiedCount)
	deleteCount, _ := mongoSerice.DeleteCustomer(cusId)
	fmt.Printf("Delete count: %v ", deleteCount)

}
