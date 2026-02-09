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

	rows, err := db.Query("SELECT id, name, logo_url FROM brands")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	fmt.Println("ID | Name | Logo URL")
	fmt.Println("--------------------------------")
	for rows.Next() {
		var id int
		var name string
		var logoUrl sql.NullString
		if err := rows.Scan(&id, &name, &logoUrl); err != nil {
			log.Fatal(err)
		}
		img := "NULL"
		if logoUrl.Valid {
			img = logoUrl.String
		}
		fmt.Printf("%d | %s | %s\n", id, name, img)
	}
}
