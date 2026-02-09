package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

func main() {
	db, err := sql.Open("sqlite", "vehicles.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 1. Verify Generations
	fmt.Println("🔎 Checking Audi A3 Generations (8V codes)...")
	rows, err := db.Query(`
		SELECT g.code, g.name, g.start_year, g.end_year, g.is_facelift, g.image_url
		FROM generations g
		JOIN models m ON g.model_id = m.id
		JOIN brands b ON m.brand_id = b.id
		WHERE b.name = 'Audi' AND m.name = 'A3' AND g.code LIKE '8V%'
		ORDER BY g.code
	`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var code, name string
		var start, end int
		var isFL bool
		var imgURL sql.NullString
		rows.Scan(&code, &name, &start, &end, &isFL, &imgURL)
		fmt.Printf("  📌 [%s] %s (%d-%d) Facelift: %v\n     🖼️ Image: %s\n", code, name, start, end, isFL, imgURL.String)
	}

	// 2. Verify Trims count per generation code
	fmt.Println("\n🔎 Checking Trim Counts...")
	rows2, err := db.Query(`
		SELECT g.code, COUNT(t.id) as trim_count
		FROM trims t
		JOIN generations g ON t.generation_id = g.id
		JOIN models m ON g.model_id = m.id
		JOIN brands b ON m.brand_id = b.id
		WHERE b.name = 'Audi' AND m.name = 'A3' AND g.code LIKE '8V%'
		GROUP BY g.code
	`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows2.Close()

	for rows2.Next() {
		var code string
		var count int
		rows2.Scan(&code, &count)
		fmt.Printf("  🚙 %s Trims: %d\n", code, count)
	}

	// 3. Sample a few trims from 8V2
	fmt.Println("\n🔎 Sample Trims from 8V2 (Facelift)...")
	rows3, err := db.Query(`
		SELECT t.name, t.power_hp, t.year 
		FROM trims t
		JOIN generations g ON t.generation_id = g.id
		WHERE g.code = '8V2'
		LIMIT 5
	`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows3.Close()

	for rows3.Next() {
		var name string
		var hp, year int
		rows3.Scan(&name, &hp, &year)
		fmt.Printf("  🔹 %s (%d hp) - %d\n", name, hp, year)
	}
}
