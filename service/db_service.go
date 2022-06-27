package service

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	rkentry "github.com/rookie-ninja/rk-entry/v2/entry"
)

const (
	selectAllCustomer      = "SELECT fname, lname, age FROM customer"
	insertNewCustomer      = "INSERT INTO customer(fname, lname, age) VALUES(?, ?, ?)"
	updateExistingCustomer = "UPDATE customer SET fname=?,  lname=?,  age=? WHERE cusid=?"
	deleteExistingCustomer = "DELETE FROM customer WHERE cusid=?"
)

type (
	DbService struct {
		DbConn func() *sql.DB
	}

	CustomerRecord struct {
		CusId int32
		Fname string
		Lname string
		Age   int
	}
)

func (s *DbService) SelectAll() ([]*CustomerRecord, error) {

	// Execute the query
	results, err := s.DbConn().Query(selectAllCustomer)
	if err != nil {
		return nil, err
	}
	defer results.Close()

	customerList := make([]*CustomerRecord, 0)
	for results.Next() {
		var cus CustomerRecord
		err = results.Scan(&cus.Fname, &cus.Lname, &cus.Age)
		if err != nil {
			return nil, err
		}
		customerList = append(customerList, &cus)
	}
	logger := rkentry.GlobalAppCtx.GetLoggerEntry("my-logger")
	logger.Info(fmt.Sprintf("Call DbService.SelectAll length: ", len(customerList)))
	return customerList, nil

}

func (s *DbService) InsertNewCustomer(cus *CustomerRecord) (int32, error) {

	fmt.Println("Call DbService.InsertNewCustomer")

	insertSmt, err := s.DbConn().Prepare(insertNewCustomer)
	if err != nil {
		return -1, err
	}
	defer insertSmt.Close()

	results, err := insertSmt.Exec(
		cus.Fname,
		cus.Lname,
		cus.Age,
	)
	if err != nil {
		return -1, err
	}

	cusId, err := results.LastInsertId()
	if err != nil {
		return -1, err
	}

	return int32(cusId), nil
}

func (s *DbService) UpdateCustomer(cus *CustomerRecord) (int32, error) {

	fmt.Println("Call DbService.UpdateCustomer")

	updateSmt, err := s.DbConn().Prepare(updateExistingCustomer)
	if err != nil {
		return -1, err
	}
	defer updateSmt.Close()

	results, err := updateSmt.Exec(
		cus.Fname,
		cus.Lname,
		cus.Age,
		cus.CusId,
	)
	if err != nil {
		return -1, err
	}

	rowsAffected, err := results.RowsAffected()
	if err != nil {
		return -1, err
	}

	return int32(rowsAffected), nil
}

func (s *DbService) DeleteCustomer(cusId int32) (int32, error) {

	fmt.Println("Call DbService.DeleteCustomer")

	deleteSmt, err := s.DbConn().Prepare(deleteExistingCustomer)
	if err != nil {
		return -1, err
	}
	defer deleteSmt.Close()

	results, err := deleteSmt.Exec(cusId)
	if err != nil {
		return -1, err
	}

	rowsAffected, err := results.RowsAffected()
	if err != nil {
		return -1, err
	}

	return int32(rowsAffected), nil
}
