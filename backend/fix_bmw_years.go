package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

type GenerationFix struct {
	Code      string
	StartYear int
	EndYear   *int
}

func main() {
	db, err := sql.Open("sqlite", "./vehicles.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	fmt.Println("BMW Nesil Yılları Düzeltme Scripti")
	fmt.Println("====================================\n")

	// Define correct generation years for BMW models
	bmw1SeriesFixes := []GenerationFix{
		{"E87", 2004, intPtr(2007)},
		{"E87 LCI", 2007, intPtr(2011)},
		{"F20", 2011, intPtr(2015)},
		{"F20 LCI", 2015, intPtr(2019)},
		{"F40", 2019, intPtr(2024)},
	}

	bmw3SeriesFixes := []GenerationFix{
		{"E90", 2005, intPtr(2008)},
		{"E90 LCI", 2008, intPtr(2012)},
		{"F30", 2012, intPtr(2015)},
		{"F30 LCI", 2015, intPtr(2019)},
		{"G20", 2019, intPtr(2022)},
		{"G20 LCI", 2022, intPtr(2024)},
	}

	bmw5SeriesFixes := []GenerationFix{
		{"E60", 2003, intPtr(2007)},
		{"E60 LCI", 2007, intPtr(2010)},
		{"F10", 2010, intPtr(2013)},
		{"F10 LCI", 2013, intPtr(2017)},
		{"G30", 2017, intPtr(2020)},
		{"G30 LCI", 2020, intPtr(2023)},
		{"G60", 2024, nil},
	}

	// First, show current state
	fmt.Println("Mevcut Durum:")
	fmt.Println("-------------")
	showCurrentBMWGenerations(db)

	// Apply fixes
	fmt.Println("\n\nDüzeltmeler Uygulanıyor:")
	fmt.Println("------------------------")

	totalFixed := 0
	totalFixed += applyFixes(db, "1 Serisi", bmw1SeriesFixes)
	totalFixed += applyFixes(db, "3 Serisi", bmw3SeriesFixes)
	totalFixed += applyFixes(db, "5 Serisi", bmw5SeriesFixes)

	// Show final state
	fmt.Println("\n\nSon Durum:")
	fmt.Println("----------")
	showCurrentBMWGenerations(db)

	fmt.Printf("\n\nToplam %d nesil güncellendi.\n", totalFixed)
}

func showCurrentBMWGenerations(db *sql.DB) {
	query := `
		SELECT m.name, g.code, g.start_year, g.end_year
		FROM generations g
		JOIN models m ON g.model_id = m.id
		JOIN brands b ON m.brand_id = b.id
		WHERE b.name = 'BMW'
		ORDER BY m.name, g.start_year DESC
	`

	rows, err := db.Query(query)
	if err != nil {
		log.Printf("Hata: %v", err)
		return
	}
	defer rows.Close()

	currentModel := ""
	for rows.Next() {
		var modelName, code string
		var startYear int
		var endYear sql.NullInt64

		err := rows.Scan(&modelName, &code, &startYear, &endYear)
		if err != nil {
			log.Printf("Hata: %v", err)
			continue
		}

		if currentModel != modelName {
			if currentModel != "" {
				fmt.Println()
			}
			fmt.Printf("BMW %s:\n", modelName)
			currentModel = modelName
		}

		endYearStr := "Günümüz"
		if endYear.Valid {
			endYearStr = fmt.Sprintf("%d", endYear.Int64)
		}
		fmt.Printf("  %s: %d - %s\n", code, startYear, endYearStr)
	}
}

func applyFixes(db *sql.DB, modelName string, fixes []GenerationFix) int {
	fmt.Printf("\nBMW %s:\n", modelName)

	// Get model ID
	var modelID int64
	err := db.QueryRow(`
		SELECT m.id 
		FROM models m 
		JOIN brands b ON m.brand_id = b.id 
		WHERE b.name = 'BMW' AND m.name = ?
	`, modelName).Scan(&modelID)

	if err != nil {
		log.Printf("  ❌ Model bulunamadı: %v", err)
		return 0
	}

	fixedCount := 0
	for _, fix := range fixes {
		// Check current values
		var currentStart int
		var currentEnd sql.NullInt64
		err := db.QueryRow(`
			SELECT start_year, end_year 
			FROM generations 
			WHERE model_id = ? AND code = ?
		`, modelID, fix.Code).Scan(&currentStart, &currentEnd)

		if err == sql.ErrNoRows {
			fmt.Printf("  ⚠️  %s: Bulunamadı\n", fix.Code)
			continue
		} else if err != nil {
			log.Printf("  ❌ %s: Hata: %v\n", fix.Code, err)
			continue
		}

		// Check if needs fixing
		needsFix := false
		if currentStart != fix.StartYear {
			needsFix = true
		}

		currentEndVal := int64(0)
		if currentEnd.Valid {
			currentEndVal = currentEnd.Int64
		}

		expectedEndVal := int64(0)
		if fix.EndYear != nil {
			expectedEndVal = int64(*fix.EndYear)
		}

		if (currentEnd.Valid != (fix.EndYear != nil)) || (currentEnd.Valid && currentEndVal != expectedEndVal) {
			needsFix = true
		}

		if !needsFix {
			fmt.Printf("  ✓ %s: Zaten doğru (%d - %v)\n", fix.Code, currentStart, formatEndYear(fix.EndYear))
			continue
		}

		// Apply fix
		_, err = db.Exec(`
			UPDATE generations 
			SET start_year = ?, end_year = ? 
			WHERE model_id = ? AND code = ?
		`, fix.StartYear, fix.EndYear, modelID, fix.Code)

		if err != nil {
			log.Printf("  ❌ %s: Güncelleme hatası: %v\n", fix.Code, err)
			continue
		}

		fmt.Printf("  ✓ %s: %d-%v → %d-%v\n",
			fix.Code,
			currentStart, formatEndYear(ptrFromNullInt(currentEnd)),
			fix.StartYear, formatEndYear(fix.EndYear))
		fixedCount++
	}

	return fixedCount
}

func intPtr(val int) *int {
	return &val
}

func ptrFromNullInt(n sql.NullInt64) *int {
	if !n.Valid {
		return nil
	}
	val := int(n.Int64)
	return &val
}

func formatEndYear(endYear *int) string {
	if endYear == nil {
		return "Günümüz"
	}
	return fmt.Sprintf("%d", *endYear)
}
