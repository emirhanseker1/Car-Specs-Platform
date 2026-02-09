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

	fmt.Println("=== FIXING A4 GENERATIONS: B9 (2016-2020) and B9.5 (2020+) ===\n")

	// Update B9 to end in 2020
	_, err = db.Exec("UPDATE generations SET name = 'Audi A4 B9 (2016-2020)', end_year = 2020 WHERE code = 'B9' AND model_id = 7")
	if err != nil {
		log.Fatal("Failed to update B9:", err)
	}
	fmt.Println("✓ Updated B9 to 2016-2020")

	// Check if B9.5 already exists
	var existingB95 int64
	err = db.QueryRow("SELECT id FROM generations WHERE code = 'B9.5' AND model_id = 7").Scan(&existingB95)
	if err == nil {
		fmt.Printf("B9.5 already exists with ID: %d\n", existingB95)
	} else {
		// Insert B9.5
		res, err := db.Exec("INSERT INTO generations (model_id, code, name, start_year, end_year, is_facelift) VALUES (7, 'B9.5', 'Audi A4 B9.5 (2020-Günümüz)', 2020, NULL, 1)")
		if err != nil {
			log.Fatal("Failed to insert B9.5:", err)
		}
		existingB95, _ = res.LastInsertId()
		fmt.Printf("✓ Created B9.5 (ID: %d)\n", existingB95)
	}

	// Get B9 ID
	var b9ID int64
	err = db.QueryRow("SELECT id FROM generations WHERE code = 'B9' AND model_id = 7").Scan(&b9ID)
	if err != nil {
		log.Fatal("B9 not found:", err)
	}

	// Move 2020+ trims from B9 to B9.5
	res, err := db.Exec("UPDATE trims SET generation_id = ? WHERE generation_id = ? AND start_year >= 2020", existingB95, b9ID)
	if err != nil {
		log.Fatal("Failed to move trims:", err)
	}
	n, _ := res.RowsAffected()
	fmt.Printf("✓ Moved %d trims from B9 to B9.5\n", n)

	// Verify
	fmt.Println("\n=== VERIFICATION ===")
	rows, _ := db.Query(`
		SELECT g.code, g.name, g.start_year, g.end_year, COUNT(t.id) as trim_count
		FROM generations g 
		LEFT JOIN trims t ON t.generation_id = g.id 
		WHERE g.model_id = 7
		GROUP BY g.id
		ORDER BY g.start_year DESC
	`)
	for rows.Next() {
		var code, name string
		var startYear int
		var endYear sql.NullInt64
		var count int
		rows.Scan(&code, &name, &startYear, &endYear, &count)
		endStr := "Günümüz"
		if endYear.Valid {
			endStr = fmt.Sprintf("%d", endYear.Int64)
		}
		fmt.Printf("  %s: %d-%s (%d trims)\n", code, startYear, endStr, count)
	}
	rows.Close()
}
