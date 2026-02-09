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

	// Find the model ID for '3 Series' or containing 'G20'
	rows, err := db.Query(`
		SELECT m.name, g.code, g.name 
		FROM generations g
		JOIN models m ON g.model_id = m.id
		WHERE g.code LIKE '%G20%' OR g.name LIKE '%G20%' OR m.name LIKE '%3 Series%'
	`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	fmt.Println("--- BMW Generations ---")
	fmt.Println("Model | Code | Name")
	fmt.Println("---------------------")
	for rows.Next() {
		var mName, gCode string
		var gName sql.NullString
		rows.Scan(&mName, &gCode, &gName)
		fmt.Printf("%s | %s | %v\n", mName, gCode, gName.String)
	}
}
