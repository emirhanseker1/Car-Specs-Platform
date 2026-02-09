package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

func main() {
	db, err := sql.Open("sqlite", "./vehicles.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	file := "./migrations/012_update_a3_images.sql"
	fmt.Printf("Running %s...\n", file)
	sqlBytes, err := os.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(string(sqlBytes))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Migration applied. Verifying A3 generations...")
	rows, err := db.Query("SELECT id, name, image_url FROM generations WHERE model_id = 1")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name string
		var imageUrl sql.NullString
		rows.Scan(&id, &name, &imageUrl)
		url := "NULL"
		if imageUrl.Valid {
			url = imageUrl.String
		}
		fmt.Printf("ID: %d | %s | Image: %s\n", id, name, url)
	}
}
