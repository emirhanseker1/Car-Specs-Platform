package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./car_specs.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Read and execute migration
	sqlBytes, err := os.ReadFile("./migrations/008_fix_audi_a3_data.sql")
	if err != nil {
		log.Fatal("Failed to read migration file:", err)
	}

	_, err = db.Exec(string(sqlBytes))
	if err != nil {
		log.Fatal("Failed to execute migration:", err)
	}

	fmt.Println("✅ Migration 008_fix_audi_a3_data.sql executed successfully!")

	// Verify results
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM trims WHERE generation_id IN (SELECT id FROM generations WHERE model_id = 1)").Scan(&count)
	if err != nil {
		log.Fatal("Failed to count trims:", err)
	}

	fmt.Printf("✅ Total A3 trims in database: %d\n", count)
}
