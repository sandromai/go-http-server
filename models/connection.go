package models

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func getDBConnection() (*sql.DB, error) {
	dataSourceName := fmt.Sprintf(
		"%v:%v@tcp(%v:%v)/%v",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_DATABASE"),
	)

	dbConnection, err := sql.Open(
		"mysql",
		dataSourceName,
	)

	if err != nil {
		return nil, err
	}

	return dbConnection, nil
}
