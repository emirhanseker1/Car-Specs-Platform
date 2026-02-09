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

	// Map of Code -> Relative Path (from /images/vehicles/audi/)
	updates := map[string]string{
		"8P":  "8p/audi-a3-8p-sportback.png",
		"8Y":  "8y/audi-a3-8y-sportback.png",
		"8V1": "8v1/audi-a3-8v-2013-2016.png",
		"8V2": "8v2/audi-a3-8v-sedan.png", // Verify this filename match
	}

	fmt.Println("🚀 Fixing Audi A3 Generation Images...")
	for code, relPath := range updates {
		fullPath := "/images/vehicles/audi/" + relPath
		res, err := db.Exec("UPDATE generations SET image_url = ? WHERE code = ?", fullPath, code)
		if err != nil {
			log.Printf("Failed to update %s: %v", code, err)
		} else {
			rows, _ := res.RowsAffected()
			if rows > 0 {
				fmt.Printf("✅ Updated %s -> %s\n", code, fullPath)
			} else {
				fmt.Printf("⚠️  Generation not found: %s\n", code)
			}
		}
	}
}
