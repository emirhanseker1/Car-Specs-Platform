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

	query := `
		SELECT m.name, g.code, g.image_url
		FROM generations g
		JOIN models m ON g.model_id = m.id
		WHERE m.name = 'Passat' AND g.code = 'B6'
	`

	var model, code string
	var imageUrl sql.NullString

	err = db.QueryRow(query).Scan(&model, &code, &imageUrl)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Model: %s, Code: %s, ImageURL: %s\n", model, code, imageUrl.String)
}
