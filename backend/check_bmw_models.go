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

	rows, err := db.Query("SELECT name FROM models WHERE name LIKE '%1 Series%' OR name LIKE '%1 Serisi%'")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	fmt.Println("--- BMW 1 Series Models ---")
	for rows.Next() {
		var name string
		rows.Scan(&name)
		fmt.Printf("'%s'\n", name)
	}
}
