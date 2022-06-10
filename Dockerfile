
FROM golang:latest

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

CMD [ "/rk-demo-app" ]
