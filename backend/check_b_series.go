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

	query := `
		SELECT b.name as brand, m.name as model, g.code, g.start_year, g.end_year
		FROM generations g
		JOIN models m ON g.model_id = m.id
		JOIN brands b ON m.brand_id = b.id
		WHERE (b.name = 'Audi' AND m.name = 'A4')
		   OR (b.name = 'Volkswagen' AND m.name = 'Passat')
		ORDER BY b.name, g.start_year
	`

	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	fmt.Printf("%-15s | %-10s | %-10s | %s\n", "Brand", "Model", "Code", "Years")
	fmt.Println("-------------------------------------------------------------")

	for rows.Next() {
		var brand, model, code string
		var start int
		var end sql.NullInt64
		if err := rows.Scan(&brand, &model, &code, &start, &end); err != nil {
			log.Fatal(err)
		}

		endStr := "Present"
		if end.Valid {
			endStr = fmt.Sprintf("%d", end.Int64)
		}
		fmt.Printf("%-15s | %-10s | %-10s | %d-%s\n", brand, model, code, start, endStr)
	}
}
