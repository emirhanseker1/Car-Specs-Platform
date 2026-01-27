package db

import (
	"database/sql"
	_ "embed"
	"log"

	_ "github.com/glebarez/go-sqlite"
)

//go:embed schema.sql
var schema string

func InitDB(dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dataSourceName)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	if _, err := db.Exec(schema); err != nil {
		return nil, err
	}

	log.Println("Database initialized and schema migrated.")
	return db, nil
}
