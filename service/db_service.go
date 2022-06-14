package service

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

const (
	selectAllCustomer      = "SELECT fname, lname, age FROM customer"
	insertNewCustomer      = "INSERT INTO customer(fname, lname, age) VALUES(?, ?, ?)"
	updateExistingCustomer = "UPDATE customer SET fname=?,  lname=?,  age=? WHERE cusid=?"
	deleteExistingCustomer = "DELETE FROM customer WHERE cusid=?"
)

type (
	DbService struct {
		DB *sql.DB
	}

	CustomerRecord struct {
		CusId int32
		Fname string
		Lname string
		Age   int
	}
)

func (s *DbService) Close() {
	log.Println("closing database")
	s.DB.Close()
}

func (s *DbService) SelectAll() ([]*CustomerRecord, error) {

	fmt.Println("Call DbService.SelectAll")

	// Execute the query
	results, err := s.DB.Query(selectAllCustomer)
	if err != nil {
		return nil, err
	}
	defer results.Close()

	customerList := make([]*CustomerRecord, 0)
	for results.Next() {
		var cus CustomerRecord
		// for each row, scan the result into our tag composite object
		err = results.Scan(&cus.Fname, &cus.Lname, &cus.Age)
		if err != nil {
			return nil, err
		}
		// and then print out the tag's Name attribute
		log.Println(cus.Fname, cus.Lname, cus.Age)

		customerList = append(customerList, &cus)
	}
	return customerList, nil

}

func (s *DbService) InsertNewCustomer(cus *CustomerRecord) (int32, error) {

	fmt.Println("Call DbService.InsertNewCustomer")

	insertSmt, err := s.DB.Prepare(insertNewCustomer)
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

	updateSmt, err := s.DB.Prepare(updateExistingCustomer)
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

	deleteSmt, err := s.DB.Prepare(deleteExistingCustomer)
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
