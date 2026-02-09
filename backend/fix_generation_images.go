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

	// Fix generation images with correct filenames
	fmt.Println("Fixing generation image paths...")

	updates := []struct {
		code     string
		modelID  int
		imageURL string
	}{
		{"8Y", 1, "/images/vehicles/audi/8y/audi-a3-8y-sportback.png"},
		{"8V", 1, "/images/vehicles/audi/8v1/audi-a3-8v-2013-2016.png"},
		{"8P", 1, "/images/vehicles/audi/8p/audi-a3-8p-sportback.png"},
		{"8L", 1, "/images/vehicles/audi/8p/audi-a3-8p-sportback.png"}, // Use 8P as placeholder for 8L
	}

	for _, u := range updates {
		_, err := db.Exec(`UPDATE generations SET image_url = ? WHERE code = ? AND model_id = ?`, u.imageURL, u.code, u.modelID)
		if err != nil {
			log.Printf("Warning: Failed to update %s: %v", u.code, err)
		} else {
			fmt.Printf("✓ Updated %s: %s\n", u.code, u.imageURL)
		}
	}

	// Verify
	fmt.Println("\nVerifying generation images:")
	rows, err := db.Query(`SELECT g.code, g.image_url FROM generations g WHERE g.model_id = 1`)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var code string
			var imageURL sql.NullString
			rows.Scan(&code, &imageURL)
			if imageURL.Valid {
				fmt.Printf("  %s: %s\n", code, imageURL.String)
			} else {
				fmt.Printf("  %s: (no image)\n", code)
			}
		}
	}
}
