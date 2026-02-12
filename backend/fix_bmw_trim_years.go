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

	fmt.Println("BMW Trim Yıllarını Düzeltme")
	fmt.Println("============================\n")

	// First, show the problem
	fmt.Println("Sorunlu Trim'ler (2020 yılıyla başlayanlar):")
	fmt.Println("--------------------------------------------")

	query := `
		SELECT b.name, m.name, g.code, g.start_year, t.name, t.year
		FROM trims t
		JOIN generations g ON t.generation_id = g.id
		JOIN models m ON g.model_id = m.id
		JOIN brands b ON m.brand_id = b.id
		WHERE b.name = 'BMW' AND t.year = 2020
		ORDER BY m.name, g.start_year DESC, t.name
	`

	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}

	problemCount := 0
	for rows.Next() {
		var brand, model, genCode, trimName string
		var genStartYear, trimYear int
		err := rows.Scan(&brand, &model, &genCode, &genStartYear, &trimName, &trimYear)
		if err != nil {
			log.Printf("Hata: %v", err)
			continue
		}

		fmt.Printf("  %s %s %s (%d): %s -> Trim yılı: %d (olması gereken: %d)\n",
			brand, model, genCode, genStartYear, trimName, trimYear, genStartYear)
		problemCount++
	}
	rows.Close()

	fmt.Printf("\nToplam %d sorunlu trim bulundu.\n\n", problemCount)

	// Fix: Update trim years to match their generation's start year
	fmt.Println("Düzeltme uygulanıyor...")
	fmt.Println("----------------------")

	result, err := db.Exec(`
		UPDATE trims 
		SET year = (
			SELECT g.start_year 
			FROM generations g 
			WHERE g.id = trims.generation_id
		)
		WHERE generation_id IN (
			SELECT g.id 
			FROM generations g
			JOIN models m ON g.model_id = m.id
			JOIN brands b ON m.brand_id = b.id
			WHERE b.name = 'BMW'
		)
		AND year = 2020
	`)

	if err != nil {
		log.Fatalf("Güncelleme hatası: %v", err)
	}

	rowsAffected, _ := result.RowsAffected()
	fmt.Printf("✅ %d trim güncellendi\n\n", rowsAffected)

	// Verify the fix
	fmt.Println("Doğrulama:")
	fmt.Println("----------")

	verifyQuery := `
		SELECT m.name, g.code, g.start_year, COUNT(t.id) as trim_count
		FROM trims t
		JOIN generations g ON t.generation_id = g.id
		JOIN models m ON g.model_id = m.id
		JOIN brands b ON m.brand_id = b.id
		WHERE b.name = 'BMW'
		GROUP BY g.id
		ORDER BY m.name, g.start_year DESC
	`

	rows, err = db.Query(verifyQuery)
	if err != nil {
		log.Fatal(err)
	}

	currentModel := ""
	for rows.Next() {
		var model, genCode string
		var genStartYear, trimCount int
		err := rows.Scan(&model, &genCode, &genStartYear, &trimCount)
		if err != nil {
			log.Printf("Hata: %v", err)
			continue
		}

		if currentModel != model {
			if currentModel != "" {
				fmt.Println()
			}
			fmt.Printf("BMW %s:\n", model)
			currentModel = model
		}

		fmt.Printf("  %s (%d): %d trim\n", genCode, genStartYear, trimCount)
	}
	rows.Close()

	// Check if any trims still have 2020
	var stillWrong int
	db.QueryRow(`
		SELECT COUNT(*) 
		FROM trims t
		JOIN generations g ON t.generation_id = g.id
		JOIN models m ON g.model_id = m.id
		JOIN brands b ON m.brand_id = b.id
		WHERE b.name = 'BMW' AND t.year = 2020 AND g.start_year != 2020
	`).Scan(&stillWrong)

	fmt.Printf("\n\nKalan sorunlu trim sayısı: %d\n", stillWrong)

	if stillWrong == 0 {
		fmt.Println("✅ Tüm BMW trim yılları başarıyla düzeltildi!")
	}
}
