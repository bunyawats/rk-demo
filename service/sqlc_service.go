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
		log.Print("SelectAll : ", err.Error())
	}
	for _, customer := range cusList {
		log.Println("SQLC Customer: ", customer.Fname.String, customer.Lname.String, customer.Age.Int32)
	}
	return cusList, err
}

func (s *SQLcService) InsertNewCustomer() int32 {

	result, err := s.queries.CreateCustomer(s.ctx, repository.CreateCustomerParams{
		Fname: sql.NullString{String: "Bunyawat_999", Valid: true},
		Lname: sql.NullString{String: "Singchai_999", Valid: true},
		Age:   sql.NullInt32{Int32: 999, Valid: true},
	})
	if err != nil {
		log.Println("InsertNewCustomer: ", err.Error())
		return -1
	}
	cusID, err := result.LastInsertId()
	if err != nil {
		log.Println("InsertNewCustomer: ", err.Error())
		return -1
	}
	return int32(cusID)
}

func (s *SQLcService) UpdateCustomer(cusId int32) {
	_, err := s.queries.UpdateCustomer(s.ctx, repository.UpdateCustomerParams{
		Fname: sql.NullString{String: "Bunyawat_888", Valid: true},
		Lname: sql.NullString{String: "Singchai_888", Valid: true},
		Age:   sql.NullInt32{Int32: 888, Valid: true},
		Cusid: cusId,
	})
	if err != nil {
		log.Println("UpdateCustomer : ", err.Error())
	}
}

func (s *SQLcService) DeleteCustomer(cusId int32) {
	err := s.queries.DeleteCustomer(s.ctx, cusId)
	if err != nil {
		log.Println("DeleteCustomer : ", err.Error())
	}
}
