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

	// Check Models
	fmt.Println("--- Models ---")
	rows, err := db.Query("SELECT name, image_url FROM models WHERE name LIKE '%A3%'")
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var name, url string
		rows.Scan(&name, &url)
		fmt.Printf("Model: %s -> %s\n", name, url)
	}
	rows.Close()

	// Check Generations
	fmt.Println("\n--- Generations ---")
	rows, err = db.Query(`
		SELECT g.code, g.image_url 
		FROM generations g
		JOIN models m ON g.model_id = m.id
		WHERE m.name LIKE '%A3%'
	`)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var code string
		var url sql.NullString
		rows.Scan(&code, &url)
		val := "NULL"
		if url.Valid {
			val = url.String
		}
		fmt.Printf("Gen: %s -> %s\n", code, val)
	}
	rows.Close()
}
