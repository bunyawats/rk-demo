package service

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	repository "github.com/rookie-ninja/rk-demo/repository/gen"
	rkentry "github.com/rookie-ninja/rk-entry/v2/entry"
)

type (
	SQLcService struct {
		conn func() *sql.DB
		ctx  context.Context
	}
)

var (
	logger *rkentry.LoggerEntry
)

func NewSQLcService(conFun func() *sql.DB, ctx context.Context) *SQLcService {
	return &SQLcService{
		conn: conFun,
		ctx:  ctx,
	}
}

func (s *SQLcService) SelectAll() ([]repository.ListCustomersRow, error) {

	logger = rkentry.GlobalAppCtx.GetLoggerEntry("my-logger")
	logger.Info("Call SQLcService.SelectAll")

	queries := repository.New(s.conn())
	cusList, err := queries.ListCustomers(s.ctx)
	if err != nil {
		logger.Info(fmt.Sprintf("SelectAll : ", err.Error()))
	}
	logger.Info(fmt.Sprintf("Call SQLcService.SelectAll length: %v", len(cusList)))
	return cusList, err
}

func (s *SQLcService) InsertNewCustomer(cus *CustomerRecord) (int32, error) {

	logger = rkentry.GlobalAppCtx.GetLoggerEntry("my-logger")
	queries := repository.New(s.conn())
	result, err := queries.CreateCustomer(s.ctx, repository.CreateCustomerParams{
		Fname: sql.NullString{String: cus.Fname, Valid: true},
		Lname: sql.NullString{String: cus.Lname, Valid: true},
		Age:   sql.NullInt32{Int32: int32(cus.Age), Valid: true},
	})
	if err != nil {
		logger.Info(fmt.Sprintf("InsertNewCustomer: %v", err.Error()))
		return -1, err
	}
	cusID, err := result.LastInsertId()
	if err != nil {
		logger.Info(fmt.Sprintf("InsertNewCustomer: %v", err.Error()))
		return -1, err
	}
	return int32(cusID), nil
}

func (s *SQLcService) UpdateCustomer(cus *CustomerRecord) (int32, error) {

	logger = rkentry.GlobalAppCtx.GetLoggerEntry("my-logger")
	queries := repository.New(s.conn())
	result, err := queries.UpdateCustomer(s.ctx, repository.UpdateCustomerParams{
		Fname: sql.NullString{String: cus.Fname, Valid: true},
		Lname: sql.NullString{String: cus.Lname, Valid: true},
		Age:   sql.NullInt32{Int32: int32(cus.Age), Valid: true},
		Cusid: cus.CusId,
	})
	if err != nil {
		logger.Info(fmt.Sprintf("UpdateCustomer : %v", err.Error()))
		return -1, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.Info(fmt.Sprintf("UpdateCustomer: %v", err.Error()))
		return -1, err
	}
	return int32(rowsAffected), nil
}

func (s *SQLcService) DeleteCustomer(cusId int32) (int32, error) {

	logger = rkentry.GlobalAppCtx.GetLoggerEntry("my-logger")
	queries := repository.New(s.conn())
	err := queries.DeleteCustomer(s.ctx, cusId)
	if err != nil {
		logger.Info(fmt.Sprintf("DeleteCustomer : %v", err.Error()))
		return 0, err
	}
	return 1, nil
}
