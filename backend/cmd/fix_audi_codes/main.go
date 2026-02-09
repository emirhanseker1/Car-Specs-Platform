package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
	_ "modernc.org/sqlite"
)

func main() {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		if _, err := os.Stat("vehicles.db"); err == nil {
			dbPath = "vehicles.db"
		} else if _, err := os.Stat("../../vehicles.db"); err == nil {
			dbPath = "../../vehicles.db"
		} else {
			dbPath = "vehicles.db"
		}
	}

	fmt.Printf("📂 Opening DB: %s\n", dbPath)
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 1. Get A3 Model ID
	var a3ID int64
	err = db.QueryRow("SELECT id FROM models WHERE name = 'A3'").Scan(&a3ID)
	if err != nil {
		log.Fatalf("A3 not found: %v", err)
	}

	// 2. Update 8V Codes
	// 8V Pre-Facelift (is_facelift = 0) -> 8V1
	res, err := db.Exec("UPDATE generations SET code = '8V1' WHERE model_id = ? AND code = '8V' AND is_facelift = 0", a3ID)
	if err != nil {
		log.Fatal(err)
	}
	aff, _ := res.RowsAffected()
	fmt.Printf("Updated 8V Pre-Facelift to 8V1: %d rows\n", aff)

	// 8V Facelift (is_facelift = 1) -> 8V2
	res, err = db.Exec("UPDATE generations SET code = '8V2' WHERE model_id = ? AND code = '8V' AND is_facelift = 1", a3ID)
	if err != nil {
		log.Fatal(err)
	}
	aff, _ = res.RowsAffected()
	fmt.Printf("Updated 8V Facelift to 8V2: %d rows\n", aff)

	// Verify
	rows, _ := db.Query("SELECT code, name FROM generations WHERE model_id = ?", a3ID)
	defer rows.Close()
	fmt.Println("--- Current A3 Generations ---")
	for rows.Next() {
		var c, n string
		rows.Scan(&c, &n)
		fmt.Printf("[%s] %s\n", c, n)
	}
}
