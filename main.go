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
	rkentry "github.com/rookie-ninja/rk-entry/v2/entry"
	rkgrpc "github.com/rookie-ninja/rk-grpc/v2/boot"
	"google.golang.org/grpc"
)

const (
	configName = "ssc-config"

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
	boot        *rkboot.Boot
)

func init() {

	boot = rkboot.NewBoot()

	db = repository.NewDbConnectionEnv(&repository.DbConnCfg{
		DbUsername: getConfigString(dbUsername),
		DbPassword: getConfigString(dbPassword),
		DbHost:     getConfigString(dbHost),
		DbName:     getConfigString(dbName),
	})

}

func getConfigString(name string) string {
	return rkentry.GlobalAppCtx.GetConfigEntry(configName).GetString(name)
}

func getDbConn() *sql.DB {
	return db
}

func main() {

	defer db.Close()

	// register grpc
	entry := rkgrpc.GetGrpcEntry("ssc-poc")
	entry.AddRegFuncGrpc(registerGreeter)
	entry.AddRegFuncGw(greeterV1.RegisterGreeterHandlerFromEndpoint)
	entry.AddRegFuncGw(greeterV2.RegisterGreeterHandlerFromEndpoint)

	// Bootstrap
	boot.Bootstrap(context.TODO())
	//db = repository.NewDbConnectionRKDB()

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
	testService()
}

func testService() {
	dbService.SelectAll()
	sqlcService.SelectAll()
}
