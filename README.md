# Example
Middleware & bootstrapper designed for [gRPC](https://grpc.io/docs/languages/go/) and [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway).

## Documentation
- [Github](https://github.com/rookie-ninja/rk-grpc)
- [Official Docs](https://docs.rkdev.info)

## Installation
- rk-boot: Bootstrapper base
- rk-grpc: Bootstrapper for [gRPC](https://grpc.io/docs/languages/go/) & [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway)

```shell
go get github.com/rookie-ninja/rk-boot/v2
go get github.com/rookie-ninja/rk-grpc/v2
```

## Quick start
### 1.Prepare .proto files
- api/v1/greeter.proto

```protobuf
syntax = "proto3";

package api.v1;

option go_package = "api/v1/hello";

service Greeter {
  rpc Hello (HelloRequest) returns (HelloResponse) {}
}

message HelloRequest {}

message HelloResponse {
  string my_message = 1;
}
```

- api/v1/gw_mapping.yaml

```yaml
type: google.api.Service
config_version: 3

# Please refer google.api.Http in third-party/googleapis/google/api/http.proto file for details.
http:
  rules:
    - selector: api.v1.Greeter.Hello
      get: /v1/hello
```

- buf.yaml

```yaml
version: v1beta1
name: github.com/rk-dev/rk-boot
build:
  roots:
    - api
    - third-party/googleapis
```

- buf.gen.yaml

```yaml
version: v1beta1
plugins:
  - name: go
    out: api/gen
    opt:
     - paths=source_relative
  - name: go-grpc
    out: api/gen
    opt:
      - paths=source_relative
      - require_unimplemented_servers=false
  - name: grpc-gateway
    out: api/gen
    opt:
      - paths=source_relative
      - grpc_api_configuration=api/v1/gw_mapping.yaml
      - allow_repeated_fields_in_body=true
      - generate_unbound_methods=true
  - name: openapiv2
    out: api/gen
    opt:
      - grpc_api_configuration=api/v1/gw_mapping.yaml
      - allow_repeated_fields_in_body=true
```

### 2.Generate .pb.go files with [buf](https://docs.buf.build/introduction)
```
$ buf generate --path api/v1
```

### 4.Create boot.yaml
Important note: rk-boot will bind grpc and grpc-gateway in the same port which we think is a convenient way.

As a result, grpc-gateway will automatically be started.

```yaml
---
grpc:
  - name: rk-demo
    port: 8080
    enabled: true
    commonService:
      enabled: true
    sw:
      enabled: true
    docs:
      enabled: true
    prom:
      enabled: true
    middleware:
      logging:
        enabled: true
      prom:
        enabled: true
```

### 5.Create main.go
```go
// Copyright (c) 2021 rookie-ninja
//
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"github.com/rookie-ninja/rk-boot/v2"
	"github.com/rookie-ninja/rk-demo/api/gen/v1"
	"github.com/rookie-ninja/rk-grpc/v2/boot"
	"google.golang.org/grpc"
)

func main() {
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
	greeter.RegisterGreeterServer(server, &GreeterServer{})
}

//GreeterServer GreeterServer struct
type GreeterServer struct{}

// Hello response with hello message
func (server *GreeterServer) Hello(_ context.Context, _ *greeter.HelloRequest) (*greeter.HelloResponse, error) {
	return &greeter.HelloResponse{
		Message: "hello!",
	}, nil
}
```

### 6.Start server

```go
$ go run main.go
```

### 7.Validation
- Call API:

```shell script
$ curl -X GET localhost:8080/v1/greeter
{"myMessage":"hello!"}

$ curl -X GET localhost:8080/rk/v1/ready
{
  "ready": true
}

$ curl -X GET localhost:8080/rk/v1/alive
{
  "alive": true
}
```

- Swagger UI: [http://localhost:8080/sw](http://localhost:8080/sw)

- Docs UI via: [http://localhost:8080/docs](http://localhost:8080/docs)

- Prometheus client: [http://localhost:8080/metrics](http://localhost:8080/metrics)


```go

func registerGreeter(server *grpc.Server) {
    greeter.RegisterGreeterServer(server, &GreeterServer{})
    reflection.Register(server)
}
```


```shell script

go install \
    github.com/rookie-ninja/rk/cmd/rk@latest


go install \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway 
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 
    google.golang.org/protobuf/cmd/protoc-gen-go 
    google.golang.org/grpc/cmd/protoc-gen-go-grpc

brew install bufbuild/buf/buf

```
- https://docs.buf.build/installation


### Test gRPC services with GRPCURL

```shell script
brew install grpcurl
grpcurl -v -d '{"name":"bunyawat"}' -plaintext localhost:8080 api.v1.Greeter.Hello
```

### Test gRPC services with GRPCUI

```shell script
brew install grpcui
grpcui -plaintext localhost:8080
```

### Test gRPC services with EVANS

```shell script
brew tap ktr0731/evans
brew install evans
evans -p 8080 -r
-show service
-service Greeter
-call Hello
```

### Run gRPC client to call Customer service

```go
$ go run cmd/customer_client.go
```

### Use MySql database driver and Viper library

```go
$ go get github.com/go-sql-driver/mysql
$ go get github.com/spf13/viper
```

### Run with VIPER configuration loader

```go
$ DB_HOSTNAME="localhost:3307" go run .
```

### Create and Run Docker image

```shell script

$ docker image ls 
$ docker run -e DB_HOSTNAME='host.docker.internal:3306' --publish 8080:8080 rk-demo-app 
$ docker container ls  
$ docker inspect rk-demo-grpc-app-1 | grep Gateway
```

### Use sqlc library

```go
$ go install github.com/kyleconroy/sqlc/cmd/sqlc@latest
$ sqlc generate
```
- https://docs.sqlc.dev/en/latest/overview/install.html

### Use Buf to generate V2 API

```shell script

$ buf ls-files
$ buf generate --path api/v1 --template api/v1/buf.gen.yaml
$ buf generate --path api/v2 --template api/v2/buf.gen.yaml
```

### Use rk-boot plugin for MySql database GORM library

```go
$ go get github.com/go-sql-driver/mysql
```


### gRPC: Enable TLS/SSL

```go
$ go install github.com/cloudflare/cfssl/cmd/...@latest
$ cfssl print-defaults config > ca-config.json
$ cfssl print-defaults csr > ca-csr.json
$ cfssl gencert -initca ca-csr.json | cfssljson -bare ca -
$ cfssl gencert -config ca-config.json -ca ca.pem -ca-key ca-key.pem -profile www ca-csr.json | cfssljson -bare server

grpcurl -v -d '{"name":"bunyawat"}' -insecure localhost:8080 api.v2.Greeter.Hello
```
- https://github.com/cloudflare/cfssl

### Override boot.yaml

```go
$ RK_MYSQL_0_USER=bunyawat123 go run .
```

### Use rk-boot plugin for MongoDB database 

```go
$ go get github.com/rookie-ninja/rk-db/mongodb
$ go get github.com/rs/xid
```
