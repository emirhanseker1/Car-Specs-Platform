package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

func main() {
	db, err := sql.Open("sqlite", "./vehicles.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	fmt.Println("Applying 022_fix_brand_logos.sql...")
	sqlBytes, err := os.ReadFile("./migrations/022_fix_brand_logos.sql")
	if err != nil {
		log.Fatal("Failed to read migration file:", err)
	}

	_, err = db.Exec(string(sqlBytes))
	if err != nil {
		log.Fatal("Failed to execute script:", err)
	}

	fmt.Println("✅ Brand Logos Updated Successfully!")
}
