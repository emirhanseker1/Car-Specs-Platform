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

	updates := map[string]string{
		"Audi":       "/images/logos/audi-logo.png",
		"BMW":        "/images/logos/bmw.png",
		"Volkswagen": "/images/logos/volkswagen.png",
	}

	fmt.Println("🚀 Updating Brand Logos...")
	for name, url := range updates {
		res, err := db.Exec("UPDATE brands SET logo_url = ? WHERE name = ?", url, name)
		if err != nil {
			log.Printf("Failed to update %s: %v", name, err)
		} else {
			rows, _ := res.RowsAffected()
			if rows > 0 {
				fmt.Printf("✅ Updated %s -> %s\n", name, url)
			} else {
				fmt.Printf("⚠️  Brand not found: %s\n", name)
			}
		}
	}
}
