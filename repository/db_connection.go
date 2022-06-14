package repository

import (
	"database/sql"
	"fmt"
	"log"
)

const (
	driverName = "mysql"

	// "test:test@tcp(127.0.0.1:3306)/test"
	dataSourceNameTemplate = "%v:%v@tcp(%v)/%v"
)

type (
	DbConnCfg struct {
		DbUsername string
		DbPassword string
		DbHost     string
		DbName     string
	}
)

func NewDbConnection(dbCfg *DbConnCfg) *sql.DB {

	dataSourceName := fmt.Sprintf(dataSourceNameTemplate,
		dbCfg.DbUsername,
		dbCfg.DbPassword,
		dbCfg.DbHost,
		dbCfg.DbName)
	log.Print("dataSourceName: ", dataSourceName)

	dbConn, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		log.Fatal(err)
	}
	err = dbConn.Ping()
	if err != nil {
		log.Fatal(err.Error())
	}

	return dbConn
}
