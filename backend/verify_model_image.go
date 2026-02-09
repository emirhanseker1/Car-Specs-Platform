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

	var id int
	var name, imageUrl string

	err = db.QueryRow("SELECT id, name, image_url FROM models WHERE name = 'A4'").Scan(&id, &name, &imageUrl)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Model: %s (ID: %d)\nImage URL: %s\n", name, id, imageUrl)
}
