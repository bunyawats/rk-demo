package repository

import (
	"database/sql"
	"fmt"
	rkmongo "github.com/rookie-ninja/rk-db/mongodb"
	rkmysql "github.com/rookie-ninja/rk-db/mysql"
	rkentry "github.com/rookie-ninja/rk-entry/v2/entry"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

const (
	configName = "ssc-config"

	dbUsername = "DB_USERNAME"
	dbPassword = "DB_PASSWORD"
	dbHost     = "DB_HOSTNAME"
	dbName     = "DB_NAME"

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

func getConfigString(name string) string {
	return rkentry.GlobalAppCtx.GetConfigEntry(configName).GetString(name)
}

func NewDbConnectionEnv(dbCfg *DbConnCfg) *sql.DB {

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

func NewMySqlDbConnectionRKDB() *sql.DB {

	mysqlEntry := rkmysql.GetMySqlEntry("test-db")
	testDb := mysqlEntry.GetDB("test")
	dbConn, err := testDb.DB()
	if err != nil {
		log.Fatal(err.Error())
	}
	err = dbConn.Ping()
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Println("Init Gorm connection success")

	return dbConn
}

func NewMongoDbConnectionRKDB() *mongo.Database {

	db := rkmongo.GetMongoDB("my-mongo", "test")

	return db
}
