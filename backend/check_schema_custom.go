package main

import (
	"database/sql"
	"fmt"
	"log"
	"path/filepath"
	"runtime"

	_ "modernc.org/sqlite"
)

func main() {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	dbPath := filepath.Join(basepath, "vehicles.db")

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("PRAGMA table_info(generations)")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	fmt.Println("Generations Table Schema:")
	for rows.Next() {
		var cid int
		var name, typeStr string
		var notnull, pk int
		var dfltValue interface{} // Use interface{} for nullable

		err = rows.Scan(&cid, &name, &typeStr, &notnull, &dfltValue, &pk)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("- %s (%s)\n", name, typeStr)
	}
}
