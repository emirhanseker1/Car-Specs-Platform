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

	// Fix 1: Update brand logo paths from /images/brands/ to /images/logos/
	fmt.Println("Fixing brand logo paths...")

	// Update Audi logo
	_, err = db.Exec(`UPDATE brands SET logo_url = '/images/logos/audi-logo.png' WHERE LOWER(name) = 'audi'`)
	if err != nil {
		log.Printf("Warning: Failed to update Audi logo: %v", err)
	}

	// Update BMW logo
	_, err = db.Exec(`UPDATE brands SET logo_url = '/images/logos/bmw.png' WHERE LOWER(name) = 'bmw'`)
	if err != nil {
		log.Printf("Warning: Failed to update BMW logo: %v", err)
	}

	// Update VW logo
	_, err = db.Exec(`UPDATE brands SET logo_url = '/images/logos/volkswagen.png' WHERE LOWER(name) = 'volkswagen'`)
	if err != nil {
		log.Printf("Warning: Failed to update VW logo: %v", err)
	}

	fmt.Println("✓ Brand logos updated")

	// Fix 2: Update Audi A3 generation images
	fmt.Println("Updating generation images...")

	// 8Y generation
	_, err = db.Exec(`UPDATE generations SET image_url = '/images/vehicles/audi/8y/8y.png' WHERE code = '8Y' AND model_id = 1`)
	if err != nil {
		log.Printf("Warning: Failed to update 8Y image: %v", err)
	}

	// 8V generation (assuming 8v1 folder has the main image)
	_, err = db.Exec(`UPDATE generations SET image_url = '/images/vehicles/audi/8v1/8v1.png' WHERE code = '8V' AND model_id = 1`)
	if err != nil {
		log.Printf("Warning: Failed to update 8V image: %v", err)
	}

	// 8P generation
	_, err = db.Exec(`UPDATE generations SET image_url = '/images/vehicles/audi/8p/8p.png' WHERE code = '8P' AND model_id = 1`)
	if err != nil {
		log.Printf("Warning: Failed to update 8P image: %v", err)
	}

	// 8L generation - check if folder exists, if not use placeholder
	_, err = db.Exec(`UPDATE generations SET image_url = '/images/vehicles/audi/a3/a3.png' WHERE code = '8L' AND model_id = 1`)
	if err != nil {
		log.Printf("Warning: Failed to update 8L image: %v", err)
	}

	fmt.Println("✓ Generation images updated")

	// Verify updates
	fmt.Println("\nVerifying brand logos:")
	rows, err := db.Query(`SELECT name, logo_url FROM brands`)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var name string
			var logoURL sql.NullString
			rows.Scan(&name, &logoURL)
			if logoURL.Valid {
				fmt.Printf("  %s: %s\n", name, logoURL.String)
			} else {
				fmt.Printf("  %s: (no logo)\n", name)
			}
		}
	}

	fmt.Println("\nVerifying generation images:")
	rows2, err := db.Query(`SELECT g.code, g.image_url FROM generations g WHERE g.model_id = 1`)
	if err == nil {
		defer rows2.Close()
		for rows2.Next() {
			var code string
			var imageURL sql.NullString
			rows2.Scan(&code, &imageURL)
			if imageURL.Valid {
				fmt.Printf("  %s: %s\n", code, imageURL.String)
			} else {
				fmt.Printf("  %s: (no image)\n", code)
			}
		}
	}
}
