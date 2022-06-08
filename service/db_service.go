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

	selectAllSql = "SELECT fname, lname, age FROM customer"
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
	results, err := s.db.Query(selectAllSql)
	if err != nil {
		return nil, err
	}

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
