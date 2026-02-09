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
		log.Fatal(err)
	}
	defer db.Close()

	fmt.Println("🛠️ Migrating: Adding image_url to models table...")

	// Check if column exists first to avoid error
	// (Simple way: try running it, ignore specific "duplicate column" error, or checking pragma)
	// For simplicity in this one-off, we'll try to run it.

	_, err = db.Exec("ALTER TABLE models ADD COLUMN image_url TEXT DEFAULT ''")
	if err != nil {
		fmt.Printf("ℹ️ Result: %v (This is fine if column already exists)\n", err)
	} else {
		fmt.Println("✅ Column 'image_url' added successfully.")
	}
}
