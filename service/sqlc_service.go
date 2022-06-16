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

func (s *SQLcService) InsertNewCustomer(cus *CustomerRecord) (int32, error) {

	result, err := s.queries.CreateCustomer(s.ctx, repository.CreateCustomerParams{
		Fname: sql.NullString{String: cus.Fname, Valid: true},
		Lname: sql.NullString{String: cus.Lname, Valid: true},
		Age:   sql.NullInt32{Int32: int32(cus.Age), Valid: true},
	})
	if err != nil {
		log.Println("InsertNewCustomer: ", err.Error())
		return -1, err
	}
	cusID, err := result.LastInsertId()
	if err != nil {
		log.Println("InsertNewCustomer: ", err.Error())
		return -1, err
	}
	return int32(cusID), nil
}

func (s *SQLcService) UpdateCustomer(cus *CustomerRecord) (int32, error) {
	result, err := s.queries.UpdateCustomer(s.ctx, repository.UpdateCustomerParams{
		Fname: sql.NullString{String: cus.Fname, Valid: true},
		Lname: sql.NullString{String: cus.Lname, Valid: true},
		Age:   sql.NullInt32{Int32: int32(cus.Age), Valid: true},
		Cusid: cus.CusId,
	})
	if err != nil {
		log.Println("UpdateCustomer : ", err.Error())
		return -1, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println("UpdateCustomer: ", err.Error())
		return -1, err
	}
	return int32(rowsAffected), nil
}

func (s *SQLcService) DeleteCustomer(cusId int32) (int32, error) {
	err := s.queries.DeleteCustomer(s.ctx, cusId)
	if err != nil {
		log.Println("DeleteCustomer : ", err.Error())
		return 0, err
	}
	return 1, nil
}
