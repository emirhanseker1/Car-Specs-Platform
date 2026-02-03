package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

func main() {
	// Open database
	db, err := sql.Open("sqlite", "./vehicles.db")
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}
	defer db.Close()

	// Enable foreign keys
	_, err = db.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		log.Fatal("Failed to enable foreign keys:", err)
	}

	// Read migration file
	sqlBytes, err := os.ReadFile("./migrations/008_fix_audi_a3_data.sql")
	if err != nil {
		log.Fatal("Failed to read migration file:", err)
	}

	// Execute migration
	fmt.Println("Executing migration 008_fix_audi_a3_data.sql...")
	_, err = db.Exec(string(sqlBytes))
	if err != nil {
		log.Fatal("Failed to execute migration:", err)
	}

	fmt.Println("✅ Migration executed successfully!")

	// Verify results
	var genCount, trimCount int

	err = db.QueryRow("SELECT COUNT(*) FROM generations WHERE model_id = 1").Scan(&genCount)
	if err != nil {
		log.Fatal("Failed to count generations:", err)
	}

	err = db.QueryRow("SELECT COUNT(*) FROM trims WHERE generation_id IN (SELECT id FROM generations WHERE model_id = 1)").Scan(&trimCount)
	if err != nil {
		log.Fatal("Failed to count trims:", err)
	}

	fmt.Printf("\n📊 Results:\n")
	fmt.Printf("   - A3 Generations: %d\n", genCount)
	fmt.Printf("   - A3 Motor Trims: %d\n", trimCount)
	fmt.Println("\n✅ All data loaded successfully!")
}
