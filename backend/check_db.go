package main

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"runtime"

	_ "modernc.org/sqlite"
)

func main() {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	dbPath := filepath.Join(basepath, "vehicles.db")
	fmt.Printf("Checking DB at: %s\n", dbPath)

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, name, transmission_type FROM trims WHERE generation_id = 1")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	fmt.Println("--- TRIMS ---")
	for rows.Next() {
		var id int
		var name, trans string
		rows.Scan(&id, &name, &trans)
		fmt.Printf("ID: %d | Name: %s | Trans: %s\n", id, name, trans)
	}
}
