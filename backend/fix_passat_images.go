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

	fmt.Println("=== FIXING VW PASSAT IMAGE PATHS ===\n")

	// Correct file names based on actual files:
	// b8-5: vw-passat-b8.5.png (not b8-5)
	// b7: vw-passat-b7.png
	// b6: cover.png
	// b5-5: EMPTY

	passatImages := map[string]string{
		"B8.5": "/images/vehicles/volkswagen/b8-5/vw-passat-b8.5.png",
		"B8":   "/images/vehicles/volkswagen/b8/vw-passat-b8.png",
		"B7":   "/images/vehicles/volkswagen/b7/vw-passat-b7.png",
		"B6":   "/images/vehicles/volkswagen/b6/cover.png",
		"B5.5": "/images/vehicles/volkswagen/b7/vw-passat-b7.png", // fallback to B7
	}

	fmt.Println("Updating Passat images:")
	for code, imgPath := range passatImages {
		res, err := db.Exec("UPDATE generations SET image_url = ? WHERE code = ? AND model_id = 3", imgPath, code)
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
	rows, _ := db.Query(`SELECT code, image_url FROM generations WHERE model_id = 3 ORDER BY start_year DESC`)
	for rows.Next() {
		var code string
		var imageURL sql.NullString
		rows.Scan(&code, &imageURL)
		img := "NULL"
		if imageURL.Valid {
			img = imageURL.String
		}
		fmt.Printf("  Passat %s: %s\n", code, img)
	}
	rows.Close()
}
