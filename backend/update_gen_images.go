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

	fmt.Println("=== UPDATING GENERATION IMAGES ===\n")

	// A4 generations image mapping
	a4Images := map[string]string{
		"B9.5": "/images/vehicles/audi/b9-5/b9.5.png",
		"B9":   "/images/vehicles/audi/b9/b9.png",
		"B8.5": "/images/vehicles/audi/b8-5/b8.5.png",
		"B8":   "/images/vehicles/audi/b8/b8.png",
		"B7":   "/images/vehicles/audi/b7/b7.png",
		"B6":   "/images/vehicles/audi/b7/b7.png", // fallback
		"B5":   "/images/vehicles/audi/b7/b7.png", // fallback
	}

	// A3 generations image mapping
	a3Images := map[string]string{
		"8Y":   "/images/vehicles/audi/8y/audi-a3-8y-sportback.png",
		"8V.5": "/images/vehicles/audi/8v2/8v2.png",
		"8V":   "/images/vehicles/audi/8v1/8v1.png",
		"8P.5": "/images/vehicles/audi/8p/8p.png",
		"8P":   "/images/vehicles/audi/8p/8p.png",
	}

	// Update A4 (model_id = 7)
	fmt.Println("Audi A4:")
	for code, imgPath := range a4Images {
		res, err := db.Exec("UPDATE generations SET image_url = ? WHERE code = ? AND model_id = 7", imgPath, code)
		if err != nil {
			fmt.Printf("  ✗ %s: %v\n", code, err)
			continue
		}
		n, _ := res.RowsAffected()
		if n > 0 {
			fmt.Printf("  ✓ %s -> %s\n", code, imgPath)
		}
	}

	// Update A3 (model_id = 1)
	fmt.Println("\nAudi A3:")
	for code, imgPath := range a3Images {
		res, err := db.Exec("UPDATE generations SET image_url = ? WHERE code = ? AND model_id = 1", imgPath, code)
		if err != nil {
			fmt.Printf("  ✗ %s: %v\n", code, err)
			continue
		}
		n, _ := res.RowsAffected()
		if n > 0 {
			fmt.Printf("  ✓ %s -> %s\n", code, imgPath)
		}
	}

	// Verify
	fmt.Println("\n=== VERIFICATION ===")
	rows, _ := db.Query(`
		SELECT m.name, g.code, g.image_url 
		FROM generations g 
		JOIN models m ON m.id = g.model_id 
		WHERE m.name IN ('A3', 'A4')
		ORDER BY m.name, g.start_year DESC
	`)
	for rows.Next() {
		var model, code string
		var imageURL sql.NullString
		rows.Scan(&model, &code, &imageURL)
		img := "NULL"
		if imageURL.Valid {
			img = imageURL.String
		}
		fmt.Printf("  %s %s: %s\n", model, code, img)
	}
	rows.Close()

	fmt.Println("\n✓ Generation images updated!")
}
