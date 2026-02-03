package main

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

func main() {
	// Try opening with standard relative path which works for server
	db, err := sql.Open("sqlite", "./car_specs.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Debug: List tables
	rows, err := db.Query("SELECT name FROM sqlite_master WHERE type='table'")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	log.Println("Found tables:")
	for rows.Next() {
		var name string
		rows.Scan(&name)
		log.Println("- " + name)
	}

	// 1. Add column
	_, err = db.Exec("ALTER TABLE generations ADD COLUMN image_url TEXT")
	if err != nil {
		log.Printf("Column might already exist: %v", err)
	} else {
		log.Println("✅ Added image_url column to generations table")
	}

	// 2. Update rows
	updates := map[string]string{
		"8P": "/images/generations/audi-a3-8p-sportback.png",
		"8V": "/images/generations/audi-a3-8v-sedan.png",
		"8Y": "/images/generations/audi-a3-8y-sportback.png",
	}

	for code, imageURL := range updates {
		result, err := db.Exec(
			"UPDATE generations SET image_url = ? WHERE code = ? AND model_id = 1",
			imageURL, code,
		)
		if err != nil {
			log.Printf("❌ Failed to update %s: %v", code, err)
			continue
		}
		rows, _ := result.RowsAffected()
		log.Printf("Updated %s: %d rows", code, rows)
	}
}
