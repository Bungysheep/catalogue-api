package database

import (
	"database/sql"
	"log"

	"github.com/bungysheep/catalogue-api/pkg/configs"
)

var (
	// DbConnection - Database connection
	DbConnection *sql.DB
)

// CreateDbConnection - Creates connection to database
func CreateDbConnection() error {
	log.Printf("Creating database connection...")

	db, err := sql.Open("mysql", configs.MYSQLTESTCONNSTRING)
	if err != nil {
		return err
	}

	DbConnection = db

	return DbConnection.Ping()
}
