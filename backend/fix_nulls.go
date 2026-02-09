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

	// Fix NULL year values
	res, err := db.Exec("UPDATE trims SET year = COALESCE(year, start_year, 2020) WHERE year IS NULL")
	if err != nil {
		log.Fatal("year update failed:", err)
	}
	n, _ := res.RowsAffected()
	fmt.Printf("Fixed %d NULL year values\n", n)

	// Fix NULL seating_capacity values
	res, err = db.Exec("UPDATE trims SET seating_capacity = 5 WHERE seating_capacity IS NULL")
	if err != nil {
		log.Fatal("seating_capacity update failed:", err)
	}
	n, _ = res.RowsAffected()
	fmt.Printf("Fixed %d NULL seating_capacity values\n", n)

	// Fix NULL currency values
	res, err = db.Exec("UPDATE trims SET currency = 'TRY' WHERE currency IS NULL")
	if err != nil {
		log.Fatal("currency update failed:", err)
	}
	n, _ = res.RowsAffected()
	fmt.Printf("Fixed %d NULL currency values\n", n)

	// Verify
	var nullCount int
	db.QueryRow("SELECT COUNT(*) FROM trims WHERE year IS NULL").Scan(&nullCount)
	fmt.Printf("Remaining NULL year values: %d\n", nullCount)
}
