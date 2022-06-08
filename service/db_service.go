package service

import (
	"database/sql"
	"fmt"
	"log"
)

type (
	DbService struct {
		db *sql.DB
	}

	CustomerRecord struct {
		Fname string
		Lname string
		Age   int
	}
)

func NewDbService(d *sql.DB) *DbService {

	return &DbService{
		db: d,
	}
}

func (s *DbService) Close() {
	s.db.Close()
}

func (s *DbService) SelectAll() ([]*CustomerRecord, error) {

	fmt.Println("Call DbService.SelectAll")

	// Execute the query
	results, err := s.db.Query("SELECT fname, lname, age FROM customer")
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
