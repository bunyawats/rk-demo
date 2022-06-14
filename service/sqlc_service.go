package service

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	repository "github.com/rookie-ninja/rk-demo/repository/gen"
)

type (
	SQLcService struct {
		queries *repository.Queries
	}
)

func NewSQLcService(db *sql.DB) *SQLcService {
	return &SQLcService{
		queries: repository.New(db),
	}
}
