package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

func main() {
	db, err := sql.Open("sqlite", "./vehicles.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Check if image_url column exists in generations table
	var count int
	err = db.QueryRow(`
		SELECT COUNT(*) 
		FROM pragma_table_info('generations') 
		WHERE name='image_url'
	`).Scan(&count)

	if err != nil {
		log.Fatalf("Failed to check column: %v", err)
	}

	if count > 0 {
		fmt.Println("✓ Column 'image_url' already exists in 'generations' table")
		return
	}

	// Add the missing column
	_, err = db.Exec(`ALTER TABLE generations ADD COLUMN image_url TEXT`)
	if err != nil {
		log.Fatalf("Failed to add column: %v", err)
	}

	fmt.Println("✓ Successfully added 'image_url' column to 'generations' table")

	// Verify
	err = db.QueryRow(`
		SELECT COUNT(*) 
		FROM pragma_table_info('generations') 
		WHERE name='image_url'
	`).Scan(&count)

	if err == nil && count > 0 {
		fmt.Println("✓ Column verified successfully")
	}
}
