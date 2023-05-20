package models

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sandromai/go-http-server/types"
)

var dbConnection *sql.DB

func getDBInstance() (*sql.DB, *types.AppError) {
	if dbConnection != nil {
		if err := dbConnection.Ping(); err == nil {
			return dbConnection, nil
		}
	}

	dataSourceName := fmt.Sprintf(
		"%v:%v@tcp(%v:%v)/%v",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_DATABASE"),
	)

	dbConnectionPool, err := sql.Open(
		"mysql",
		dataSourceName,
	)

	if err != nil {
		return nil, &types.AppError{
			StatusCode: 500,
			Message:    "Failed to connect to database.",
		}
	}

	if err = dbConnectionPool.Ping(); err != nil {
		return nil, &types.AppError{
			StatusCode: 500,
			Message:    "Error connecting to database.",
		}
	}

	dbConnectionPool.SetMaxIdleConns(15)
	dbConnectionPool.SetMaxOpenConns(25)
	dbConnectionPool.SetConnMaxIdleTime(1 * time.Second)
	dbConnectionPool.SetConnMaxLifetime(30 * time.Second)

	dbConnection = dbConnectionPool

	return dbConnection, nil
}
