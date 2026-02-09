package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

func main() {
	db, err := sql.Open("sqlite", "./vehicles.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Enable foreign keys
	_, err = db.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		log.Fatal("Failed to enable foreign keys:", err)
	}

	fmt.Println("🚀 Starting A4 Data Population...")

	// Execute 013 script
	fmt.Println("Applying 013_populate_audi_a4_detailed.sql...")
	sqlBytes, err := os.ReadFile("./migrations/013_populate_audi_a4_detailed.sql")
	if err != nil {
		log.Fatal("Failed to read migration file:", err)
	}

	_, err = db.Exec(string(sqlBytes))
	if err != nil {
		log.Fatal("Failed to execute 013 script:", err)
	}

	fmt.Println("✅ A4 Data Population Complete!")

	// Verify
	var genCount, trimCount int
	// Get count of A4 generations
	db.QueryRow(`
		SELECT COUNT(*) FROM generations 
		WHERE model_id = (SELECT id FROM models WHERE name = 'A4')
	`).Scan(&genCount)

	// Get count of A4 trims
	db.QueryRow(`
		SELECT COUNT(*) FROM trims 
		JOIN generations ON trims.generation_id = generations.id 
		JOIN models ON generations.model_id = models.id 
		WHERE models.name = 'A4'
	`).Scan(&trimCount)

	fmt.Printf("Final Stats: %d Generations, %d Trims for A4.\n", genCount, trimCount)
}
