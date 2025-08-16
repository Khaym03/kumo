package db

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

func NewSQLiteConn() (*sql.DB, error) {
	dbPath := "db/kumo.db"

	dbConn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Println("Error opening database:", err)
		return nil, err
	}

	err = dbConn.Ping()
	if err != nil {
		log.Println("Error pinging database:", err)
		return nil, err
	}

	return dbConn, nil
}
