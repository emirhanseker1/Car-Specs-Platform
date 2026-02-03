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

	// Get trims table schema
	rows, err := db.Query("PRAGMA table_info(trims)")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	fmt.Println("=== TRIMS TABLE SCHEMA ===")
	for rows.Next() {
		var cid int
		var name, typ string
		var notnull, pk int
		var dfltValue sql.NullString

		err = rows.Scan(&cid, &name, &typ, &notnull, &dfltValue, &pk)
		if err != nil {
			log.Fatal(err)
		}

		required := ""
		if notnull == 1 {
			required = " NOT NULL"
		}
		fmt.Printf("%d. %s %s%s\n", cid, name, typ, required)
	}

	// Get generations table schema
	rows2, err := db.Query("PRAGMA table_info(generations)")
	if err != nil {
		log.Fatal(err)
	}
	defer rows2.Close()

	fmt.Println("\n=== GENERATIONS TABLE SCHEMA ===")
	for rows2.Next() {
		var cid int
		var name, typ string
		var notnull, pk int
		var dfltValue sql.NullString

		err = rows2.Scan(&cid, &name, &typ, &notnull, &dfltValue, &pk)
		if err != nil {
			log.Fatal(err)
		}

		required := ""
		if notnull == 1 {
			required = " NOT NULL"
		}
		fmt.Printf("%d. %s %s%s\n", cid, name, typ, required)
	}
}
