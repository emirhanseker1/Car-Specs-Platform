package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
	_ "modernc.org/sqlite"
)

func main() {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "vehicles.db"
	}
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Update 8V1 details
	// Set Name = '8V (Makyajsız)' and EndYear = 2016
	res, err := db.Exec("UPDATE generations SET name = '8V (Makyajsız)', end_year = 2016 WHERE code = '8V1'")
	if err != nil {
		log.Fatal(err)
	}

	aff, _ := res.RowsAffected()
	fmt.Printf("Updated 8V1 Details: %d rows\n", aff)
}
