
##
## Build
##
FROM golang:latest AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY *.go ./
COPY *.env ./
COPY *.yaml ./
COPY api ./api
COPY service ./service

RUN go build -o /rk-demo-app

##
## Deploy
##
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /rk-demo-app /rk-demo-app
COPY *.env ./
COPY *.yaml ./

EXPOSE 8080

USER root:root

CMD [ "/rk-demo-app" ]
