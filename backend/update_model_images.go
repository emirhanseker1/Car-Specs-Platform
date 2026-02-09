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
		"A3":       "/images/models/audi/audi-a3.png",
		"A4":       "/images/models/audi/audi-a4.png",
		"Golf":     "/images/models/vw/vw-golf.png",
		"Passat":   "/images/models/vw/vw-passat.png",
		"1 Serisi": "/images/models/bmw/bmw-1-series.png",
		"3 Serisi": "/images/models/bmw/bmw-3-series.png",
		"5 Serisi": "/images/models/bmw/bmw-5-series.png",
	}

	fmt.Println("🚀 Updating Model Images...")
	for name, url := range updates {
		res, err := db.Exec("UPDATE models SET image_url = ? WHERE name = ?", url, name)
		if err != nil {
			log.Printf("Failed to update %s: %v", name, err)
		} else {
			rows, _ := res.RowsAffected()
			if rows > 0 {
				fmt.Printf("✅ Updated %s -> %s\n", name, url)
			} else {
				fmt.Printf("⚠️  Model not found: %s\n", name)
			}
		}
	}
}
