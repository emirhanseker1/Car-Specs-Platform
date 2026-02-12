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

	fmt.Println("BMW Model ve Nesil Bilgileri")
	fmt.Println("=============================\n")

	// Get BMW brand ID
	var brandID int64
	err = db.QueryRow("SELECT id FROM brands WHERE name = 'BMW'").Scan(&brandID)
	if err != nil {
		log.Fatalf("BMW markası bulunamadı: %v", err)
	}
	fmt.Printf("BMW Brand ID: %d\n\n", brandID)

	// Get all BMW models
	rows, err := db.Query(`
		SELECT id, name FROM models WHERE brand_id = ? ORDER BY name
	`, brandID)
	if err != nil {
		log.Fatalf("Modeller alınamadı: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var modelID int64
		var modelName string
		err := rows.Scan(&modelID, &modelName)
		if err != nil {
			log.Printf("Hata: %v", err)
			continue
		}

		fmt.Printf("Model: %s (ID: %d)\n", modelName, modelID)
		fmt.Println("Nesiller:")

		// Get generations for this model
		genRows, err := db.Query(`
			SELECT id, code, name, start_year, end_year
			FROM generations
			WHERE model_id = ?
			ORDER BY start_year DESC
		`, modelID)
		if err != nil {
			log.Printf("  Nesiller alınamadı: %v", err)
			continue
		}

		for genRows.Next() {
			var genID int64
			var code string
			var name sql.NullString
			var startYear int
			var endYear sql.NullInt64

			err := genRows.Scan(&genID, &code, &name, &startYear, &endYear)
			if err != nil {
				log.Printf("  Hata: %v", err)
				continue
			}

			nameStr := ""
			if name.Valid {
				nameStr = name.String
			}

			endYearStr := "Günümüz"
			if endYear.Valid {
				endYearStr = fmt.Sprintf("%d", endYear.Int64)
			}

			fmt.Printf("  - %s (%s): %d - %s (ID: %d)\n", code, nameStr, startYear, endYearStr, genID)
		}
		genRows.Close()
		fmt.Println()
	}
}
