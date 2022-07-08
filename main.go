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
	rkgrpc "github.com/rookie-ninja/rk-grpc/v2/boot"
	"google.golang.org/grpc"
	"os"
)

var (
	db          *sql.DB
	dbService   *service.DbService
	sqlcService *service.SQLcService
	boot        *rkboot.Boot
)

func getDbConn() *sql.DB {
	return db
}

func main() {

	os.Setenv("RK_MYSQL_0_USER", "bunyawat444")

	boot = rkboot.NewBoot()

	// register grpc
	entry := rkgrpc.GetGrpcEntry("ssc-grpc")
	entry.AddRegFuncGrpc(registerGreeter)
	entry.AddRegFuncGw(greeterV1.RegisterGreeterHandlerFromEndpoint)
	entry.AddRegFuncGw(greeterV2.RegisterGreeterHandlerFromEndpoint)

	entry2 := rkgrpc.GetGrpcEntry("ssc-app")
	entry2.AddRegFuncGrpc(registerGreeter)
	entry2.AddRegFuncGw(greeterV1.RegisterGreeterHandlerFromEndpoint)
	entry2.AddRegFuncGw(greeterV2.RegisterGreeterHandlerFromEndpoint)

	// Bootstrap
	boot.Bootstrap(context.TODO())
	db = repository.NewDbConnectionRKDB()
	defer db.Close()

	testService()

	// Wait for shutdown sig
	boot.WaitForShutdownSig(context.TODO())
}

func registerGreeter(server *grpc.Server) {

	dbService = &service.DbService{DbConn: getDbConn}
	sqlcService = service.NewSQLcService(getDbConn, context.TODO())

	greeterV1.RegisterGreeterServer(server, &v1.GreeterServer{})
	greeterV1.RegisterCustomerServer(server, v1.NewCustomerServer(context.TODO(), dbService))

	greeterV2.RegisterGreeterServer(server, &v2.GreeterServer{})
	greeterV2.RegisterCustomerServer(server, v2.NewCustomerServer(context.TODO(), sqlcService))

	//	reflection.Register(server)

}

func testService() {
	dbService.SelectAll()
	sqlcService.SelectAll()
}
