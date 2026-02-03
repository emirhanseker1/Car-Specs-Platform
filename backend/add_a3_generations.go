package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Load database path from environment
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./car_specs.db"
	}

	// Open database
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	log.Println("🚀 Running Audi A3 Generations Migration...")

	// Execute migration
	migrationSQL := `
INSERT INTO generations (model_id, code, name, start_year, end_year, is_current) VALUES
(1, '8L', 'Typ 8L', 1996, 2003, 0),
(1, '8P', 'Typ 8P', 2003, 2012, 0),
(1, '8V', 'Typ 8V', 2012, 2020, 0),
(1, '8Y', 'Typ 8Y', 2020, NULL, 1);
`

	_, err = db.Exec(migrationSQL)
	if err != nil {
		log.Fatalf("❌ Migration failed: %v", err)
	}

	log.Println("✅ Successfully added Audi A3 generations!")

	// Verify
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM generations WHERE model_id = 1").Scan(&count)
	if err != nil {
		log.Printf("⚠️  Warning: Could not verify: %v", err)
	} else {
		log.Printf("📊 Total generations for Audi A3: %d", count)
	}
}
