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

	rows, err := db.Query("SELECT id, name, code, image_url FROM generations")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	fmt.Println("ID | Name | Code | Image URL")
	fmt.Println("--------------------------------")
	for rows.Next() {
		var id int
		var name, code string
		var imageUrl sql.NullString
		if err := rows.Scan(&id, &name, &code, &imageUrl); err != nil {
			log.Fatal(err)
		}
		img := "NULL"
		if imageUrl.Valid {
			img = imageUrl.String
		}
		fmt.Printf("%d | %s | %s | %s\n", id, name, code, img)
	}
}
