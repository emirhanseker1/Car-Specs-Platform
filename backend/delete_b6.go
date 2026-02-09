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

	// Enable foreign keys
	db.Exec("PRAGMA foreign_keys = ON")

	fmt.Println("🚀 Removing Audi A4 B6 Generation and related trims...")

	// 1. Get Model ID for Audi A4
	var modelID int
	err = db.QueryRow(`
		SELECT m.id FROM models m
		JOIN brands b ON m.brand_id = b.id
		WHERE b.name = 'Audi' AND m.name = 'A4'
	`).Scan(&modelID)
	if err != nil {
		log.Fatalf("❌ Audi A4 model not found: %v", err)
	}

	// 2. Get Generation ID for B6
	var genID int
	err = db.QueryRow("SELECT id FROM generations WHERE model_id = ? AND code = 'B6'", modelID).Scan(&genID)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("✅ A4 B6 generation not found (already deleted?).")
			return
		}
		log.Fatalf("❌ Failed to query generation B6: %v", err)
	}

	// 3. Delete Trims
	// Since FKs are ON and maybe ON DELETE CASCADE isn't set, we delete trims first.
	res, err := db.Exec("DELETE FROM trims WHERE generation_id = ?", genID)
	if err != nil {
		log.Fatalf("❌ Failed to delete trims: %v", err)
	}
	count, _ := res.RowsAffected()
	fmt.Printf("✅ Deleted %d trims for A4 B6.\n", count)

	// 4. Delete Generation
	res, err = db.Exec("DELETE FROM generations WHERE id = ?", genID)
	if err != nil {
		log.Fatalf("❌ Failed to delete generation: %v", err)
	}
	count, _ = res.RowsAffected()
	fmt.Printf("✅ Deleted A4 B6 generation (ID: %d).\n", genID)
}
