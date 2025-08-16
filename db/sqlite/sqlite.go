package db

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

func NewSQLiteConn() *sql.DB {
	dbPath := "db/sqlite/kumo.db"

	dbConn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatal("Error opening database:", err)
	}

	err = dbConn.Ping()
	if err != nil {
		log.Fatal("Error pinging database:", err)
	}

	return dbConn
}
