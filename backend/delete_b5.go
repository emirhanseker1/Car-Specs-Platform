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

	// Enable foreign keys
	db.Exec("PRAGMA foreign_keys = ON")

	fmt.Println("🧹 Removing Audi A4 B5 Generation...")

	// Find ID
	var id int
	err = db.QueryRow("SELECT id FROM generations WHERE code = 'B5' AND model_id = (SELECT id FROM models WHERE name = 'A4')").Scan(&id)
	if err != nil {
		fmt.Println("B5 already gone or not found.")
		return
	}

	res, err := db.Exec("DELETE FROM generations WHERE id = ?", id)
	if err != nil {
		log.Fatal(err)
	}

	n, _ := res.RowsAffected()
	fmt.Printf("✅ Deleted %d generation(s) (B5).\n", n)
}
