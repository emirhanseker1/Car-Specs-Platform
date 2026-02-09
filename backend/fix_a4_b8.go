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

	fmt.Println("=== FIXING A4 B8: B8 (2008-2011) and B8.5 (2011-2015) ===\n")

	// Update B8 to end in 2011
	_, err = db.Exec("UPDATE generations SET name = 'Audi A4 B8 (2008-2011)', end_year = 2011 WHERE code = 'B8' AND model_id = 7")
	if err != nil {
		log.Fatal("Failed to update B8:", err)
	}
	fmt.Println("✓ Updated B8 to 2008-2011")

	// Check if B8.5 already exists
	var existingB85 int64
	err = db.QueryRow("SELECT id FROM generations WHERE code = 'B8.5' AND model_id = 7").Scan(&existingB85)
	if err == nil {
		fmt.Printf("B8.5 already exists with ID: %d\n", existingB85)
	} else {
		// Insert B8.5
		res, err := db.Exec("INSERT INTO generations (model_id, code, name, start_year, end_year, is_facelift) VALUES (7, 'B8.5', 'Audi A4 B8.5 (2011-2015)', 2011, 2015, 1)")
		if err != nil {
			log.Fatal("Failed to insert B8.5:", err)
		}
		existingB85, _ = res.LastInsertId()
		fmt.Printf("✓ Created B8.5 (ID: %d)\n", existingB85)
	}

	// Get B8 ID
	var b8ID int64
	err = db.QueryRow("SELECT id FROM generations WHERE code = 'B8' AND model_id = 7").Scan(&b8ID)
	if err != nil {
		log.Fatal("B8 not found:", err)
	}

	// Move 2011+ trims (makyajlı) from B8 to B8.5
	res, err := db.Exec("UPDATE trims SET generation_id = ? WHERE generation_id = ? AND start_year >= 2011", existingB85, b8ID)
	if err != nil {
		log.Fatal("Failed to move trims:", err)
	}
	n, _ := res.RowsAffected()
	fmt.Printf("✓ Moved %d trims from B8 to B8.5\n", n)

	// Also update B7 end year to 2008 (not 2008)
	db.Exec("UPDATE generations SET end_year = 2008 WHERE code = 'B7' AND model_id = 7")

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
