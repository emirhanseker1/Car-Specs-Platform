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

	fmt.Println("=== UPDATING VW GENERATION IMAGES ===\n")

	// Get VW brand ID and model IDs
	var vwID int64
	err = db.QueryRow("SELECT id FROM brands WHERE name = 'Volkswagen'").Scan(&vwID)
	if err != nil {
		log.Fatal("Volkswagen not found:", err)
	}

	var golfID, passatID int64
	db.QueryRow("SELECT id FROM models WHERE name = 'Golf' AND brand_id = ?", vwID).Scan(&golfID)
	db.QueryRow("SELECT id FROM models WHERE name = 'Passat' AND brand_id = ?", vwID).Scan(&passatID)

	fmt.Printf("VW ID: %d, Golf ID: %d, Passat ID: %d\n\n", vwID, golfID, passatID)

	// Golf generations image mapping
	golfImages := map[string]string{
		"MK8":   "/images/vehicles/volkswagen/mk8/vw-golf-8.png",
		"MK7.5": "/images/vehicles/volkswagen/mk7-5/vw-golf-7-5.png",
		"MK7":   "/images/vehicles/volkswagen/mk7/vw-golf-7.png",
		"MK6":   "/images/vehicles/volkswagen/mk6/vw-golf-6.png",
		"MK5":   "/images/vehicles/volkswagen/mk5/vw-golf-5.png",
	}

	// Passat generations image mapping
	passatImages := map[string]string{
		"B8.5": "/images/vehicles/volkswagen/b8-5/vw-passat-b8-5.png",
		"B8":   "/images/vehicles/volkswagen/b8/vw-passat-b8.png",
		"B7":   "/images/vehicles/volkswagen/b7/vw-passat-b7.png",
		"B6":   "/images/vehicles/volkswagen/b6/vw-passat-b6.png",
		"B5.5": "/images/vehicles/volkswagen/b5-5/vw-passat-b5-5.png",
	}

	// Update Golf
	fmt.Println("Golf:")
	for code, imgPath := range golfImages {
		res, err := db.Exec("UPDATE generations SET image_url = ? WHERE code = ? AND model_id = ?", imgPath, code, golfID)
		if err != nil {
			fmt.Printf("  ✗ %s: %v\n", code, err)
			continue
		}
		n, _ := res.RowsAffected()
		if n > 0 {
			fmt.Printf("  ✓ %s -> %s\n", code, imgPath)
		} else {
			fmt.Printf("  - %s: no matching generation\n", code)
		}
	}

	// Update Passat
	fmt.Println("\nPassat:")
	for code, imgPath := range passatImages {
		res, err := db.Exec("UPDATE generations SET image_url = ? WHERE code = ? AND model_id = ?", imgPath, code, passatID)
		if err != nil {
			fmt.Printf("  ✗ %s: %v\n", code, err)
			continue
		}
		n, _ := res.RowsAffected()
		if n > 0 {
			fmt.Printf("  ✓ %s -> %s\n", code, imgPath)
		} else {
			fmt.Printf("  - %s: no matching generation\n", code)
		}
	}

	// Verify
	fmt.Println("\n=== VERIFICATION ===")
	rows, _ := db.Query(`
		SELECT m.name, g.code, g.image_url 
		FROM generations g 
		JOIN models m ON m.id = g.model_id 
		WHERE m.brand_id = ?
		ORDER BY m.name, g.start_year DESC
	`, vwID)
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

	fmt.Println("\n✓ VW generation images updated!")
}
