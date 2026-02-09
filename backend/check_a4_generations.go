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

	modelID := 2 // A4
	fmt.Printf("\nAnalyzing generations for Model ID: %d (A4)\n", modelID)

	rows, err := db.Query("SELECT id, code, name, image_url FROM generations WHERE model_id = ?", modelID)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var code, gName string
		var img sql.NullString
		rows.Scan(&id, &code, &gName, &img)

		imgVal := "<NULL>"
		if img.Valid {
			imgVal = img.String
		}

		fmt.Printf("ID: %d | Code: %-4s | Name: %-20s | Img: %s\n", id, code, gName, imgVal)
	}
}
