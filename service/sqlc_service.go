package service

import (
	"context"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	repository "github.com/rookie-ninja/rk-demo/repository/gen"
	"log"
)

type (
	SQLcService struct {
		queries *repository.Queries
		ctx     context.Context
	}
)

func NewSQLcService(db *sql.DB, ctx context.Context) *SQLcService {
	return &SQLcService{
		queries: repository.New(db),
		ctx:     ctx,
	}
}

func (s *SQLcService) SelectAll() ([]repository.ListCustomersRow, error) {
	cusList, err := s.queries.ListCustomers(s.ctx)
	if err != nil {
		log.Print(err.Error())
	}
	for _, customer := range cusList {
		log.Println("SQLC Customer: ", customer.Fname.String, customer.Lname.String, customer.Age.Int32)
	}
	return cusList, err
}
