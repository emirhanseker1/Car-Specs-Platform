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

	rows, err := db.Query("SELECT id, name, image_url FROM models")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	fmt.Println("ID | Name | Image URL")
	fmt.Println("--------------------------------")
	for rows.Next() {
		var id int
		var name string
		var imageUrl sql.NullString
		if err := rows.Scan(&id, &name, &imageUrl); err != nil {
			log.Fatal(err)
		}
		img := "NULL"
		if imageUrl.Valid {
			img = imageUrl.String
		}
		fmt.Printf("%d | %s | %s\n", id, name, img)
	}
}
