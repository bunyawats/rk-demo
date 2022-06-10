package service

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

const (
	driverName = "mysql"

	// "test:test@tcp(127.0.0.1:3306)/test"
	dataSourceNameTemplate = "%v:%v@tcp(%v)/%v"

	selectAllCustomer      = "SELECT fname, lname, age FROM customer"
	insertNewCustomer      = "INSERT INTO customer(fname, lname, age) VALUES(?, ?, ?)"
	updateExistingCustomer = "UPDATE customer SET fname=?,  lname=?,  age=? WHERE cusid=?"
	deleteExistingCustomer = "DELETE FROM customer WHERE cusid=?"
)

type (
	DbService struct {
		db *sql.DB
	}

	DbConnCfg struct {
		DbUsername string
		DbPassword string
		DbHost     string
		DbName     string
	}

	CustomerRecord struct {
		CusId int32
		Fname string
		Lname string
		Age   int
	}
)

func NewDbService(dbCfg *DbConnCfg) (*DbService, error) {

	dataSourceName := fmt.Sprintf(dataSourceNameTemplate,
		dbCfg.DbUsername,
		dbCfg.DbPassword,
		dbCfg.DbHost,
		dbCfg.DbName)
	log.Print("dataSourceName: ", dataSourceName)

	d, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}

	return &DbService{db: d}, nil
}

func (s *DbService) Close() {
	log.Println("closing database")
	s.db.Close()
}

func (s *DbService) SelectAll() ([]*CustomerRecord, error) {

	fmt.Println("Call DbService.SelectAll")

	// Execute the query
	results, err := s.db.Query(selectAllCustomer)
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

	insertSmt, err := s.db.Prepare(insertNewCustomer)
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

	updateSmt, err := s.db.Prepare(updateExistingCustomer)
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

	deleteSmt, err := s.db.Prepare(deleteExistingCustomer)
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
