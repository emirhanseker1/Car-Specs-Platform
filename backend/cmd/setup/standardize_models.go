package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

// GenerationDefinition represents the desired state for a generation
type GenerationDefinition struct {
	Code      string
	Name      string
	StartYear int
	EndYear   *int   // nil for current
	ImageURL  string // can be empty
	Platform  string
}

// ModelDefinition represents a model and its generations
type ModelDefinition struct {
	BrandName   string
	ModelName   string
	ImageURL    string // New field
	Generations []GenerationDefinition
}

func main() {
	db, err := sql.Open("sqlite", "./vehicles.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Enable foreign keys
	db.Exec("PRAGMA foreign_keys = ON")

	// --- 1. DEFINE DATA ---
	models := []ModelDefinition{
		{
			BrandName: "Audi",
			ModelName: "A3",
			Generations: []GenerationDefinition{
				{Code: "8L", Name: "Tip 8L (1996-2003)", StartYear: 1996, EndYear: intPtr(2003), ImageURL: "", Platform: "PQ34"},
				{Code: "8P", Name: "Tip 8P (2003-2012)", StartYear: 2003, EndYear: intPtr(2012), ImageURL: "/images/generations/audi-a3-8p-sportback.png", Platform: "PQ35"},
				{Code: "8V", Name: "Tip 8V (2012-2020)", StartYear: 2012, EndYear: intPtr(2020), ImageURL: "/images/generations/audi-a3-8v-sedan.png", Platform: "MQB"},
				{Code: "8Y", Name: "Tip 8Y (2020-Günümüz)", StartYear: 2020, EndYear: nil, ImageURL: "/images/generations/audi-a3-8y-sportback.png", Platform: "MQB Evo"},
			},
		},
		{
			BrandName: "Audi",
			ModelName: "A4",
			ImageURL:  "/images/vehicles/audi/a4/a4model.png",
			Generations: []GenerationDefinition{
				{Code: "B7", Name: "Tip B7 (2005-2008)", StartYear: 2005, EndYear: intPtr(2008), ImageURL: "/images/vehicles/audi/a4/b7.png", Platform: "PL46"},
				{Code: "B8", Name: "Tip B8 (2008-2011)", StartYear: 2008, EndYear: intPtr(2011), ImageURL: "/images/vehicles/audi/a4/b8.png", Platform: "MLB"},
				{Code: "B8.5", Name: "Tip B8.5 (2011-2015)", StartYear: 2011, EndYear: intPtr(2015), ImageURL: "/images/vehicles/audi/a4/b8.5.png", Platform: "MLB"},
				{Code: "B9", Name: "Tip B9 (2016-2019)", StartYear: 2016, EndYear: intPtr(2019), ImageURL: "/images/vehicles/audi/a4/b9.png", Platform: "MLB Evo"},
				{Code: "B9.5", Name: "Tip B9.5 (2020-Günümüz)", StartYear: 2020, EndYear: nil, ImageURL: "/images/vehicles/audi/a4/b9.5.png", Platform: "MLB Evo"},
			},
		},
	}

	// --- 2. APPLY ----
	fmt.Println("🚀 Starting Model Standardization...")

	for _, m := range models {
		fmt.Printf("\nProcessing %s %s...\n", m.BrandName, m.ModelName)

		// Get Brand ID
		var brandID int
		err := db.QueryRow("SELECT id FROM brands WHERE name = ?", m.BrandName).Scan(&brandID)
		if err != nil {
			log.Printf("❌ Brand %s not found, skipping.", m.BrandName)
			continue
		}

		// Get Model ID
		var modelID int
		err = db.QueryRow("SELECT id FROM models WHERE name = ? AND brand_id = ?", m.ModelName, brandID).Scan(&modelID)
		if err != nil {
			// Create model if missing? (Skipping for now to be safe, assuming models exist)
			log.Printf("❌ Model %s not found, skipping.", m.ModelName)
			continue
		}

		// UPDATE MODEL IMAGE IF SET
		if m.ImageURL != "" {
			_, err := db.Exec("UPDATE models SET image_url = ? WHERE id = ?", m.ImageURL, modelID)
			if err != nil {
				log.Printf("❌ Failed to update model image: %v", err)
			} else {
				fmt.Printf("✅ Model image updated to %s\n", m.ImageURL)
			}
		}

		for _, gen := range m.Generations {
			// A. Upsert Generation
			fmt.Printf("  - Generation %s (%d-%v)... ", gen.Code, gen.StartYear, derefInt(gen.EndYear))

			var genID int
			err := db.QueryRow("SELECT id FROM generations WHERE code = ? AND model_id = ?", gen.Code, modelID).Scan(&genID)

			if err == sql.ErrNoRows {
				// Insert
				res, err := db.Exec(`
					INSERT INTO generations (model_id, code, name, start_year, end_year, image_url, platform, is_current)
					VALUES (?, ?, ?, ?, ?, ?, ?, ?)
				`, modelID, gen.Code, gen.Name, gen.StartYear, gen.EndYear, gen.ImageURL, gen.Platform, isCurrent(gen.EndYear))
				if err != nil {
					log.Printf("❌ Failed to insert: %v", err)
					continue
				}
				id, _ := res.LastInsertId()
				genID = int(id)
				fmt.Print("✅ Created. ")
			} else if err == nil {
				// Update
				_, err := db.Exec(`
					UPDATE generations 
					SET name=?, start_year=?, end_year=?, image_url=?, platform=?, is_current=?
					WHERE id=?
				`, gen.Name, gen.StartYear, gen.EndYear, gen.ImageURL, gen.Platform, isCurrent(gen.EndYear), genID)
				if err != nil {
					log.Printf("❌ Failed to update: %v", err)
					continue
				}
				fmt.Print("✅ Updated. ")
			} else {
				log.Printf("❌ DB Error: %v", err)
				continue
			}

			// B. Link Trims by Year
			// Logic: trims.year >= gen.StartYear AND (gen.EndYear IS NULL OR trims.year <= gen.EndYear)
			// AND model_id matches

			query := `
				UPDATE trims 
				SET generation_id = ?
				WHERE model_id = ? 
				AND year >= ?
			`
			args := []interface{}{genID, modelID, gen.StartYear}

			if gen.EndYear != nil {
				query += " AND year <= ?"
				args = append(args, *gen.EndYear)
			}

			// Be careful not to overwrite manual assignments if needed,
			// but here we want to ENFORCE the standard.
			// However, overlapping years (production changes mid-year) are tricky.
			// Simple logic: if a year matches multiple generations (e.g. 2012 is both end of 8P and start of 8V),
			// this script will overwrite with whichever generation runs last in the loop.
			// Since we defined order chronological, later gens will claim the overlap year.
			// E.g. A3 8V (2012 start) comes after 8P (2012 end). 8V will claim 2012 trims.
			// This is usually correct for "New Model Year".

			res, err := db.Exec(query, args...)
			if err != nil {
				log.Printf("❌ Failed to link trims: %v", err)
			} else {
				count, _ := res.RowsAffected()
				if count > 0 {
					fmt.Printf("Linked %d trims.", count)
				}
			}
			fmt.Println()
		}
	}

	// --- 3. CLEANUP ---
	fmt.Println("\n🧹 Cleaning up empty fallback generations...")
	// Delete generations where code starts with 'Y' and has NO trims
	res, err := db.Exec(`
		DELETE FROM generations 
		WHERE code LIKE 'Y%' 
		AND id NOT IN (SELECT DISTINCT generation_id FROM trims WHERE generation_id IS NOT NULL)
	`)
	if err != nil {
		log.Printf("❌ Cleanup failed: %v", err)
	} else {
		count, _ := res.RowsAffected()
		if count > 0 {
			fmt.Printf("✅ Removed %d unused fallback generations.\n", count)
		} else {
			fmt.Println("No unused generations found.")
		}
	}

	fmt.Println("\n✅ Standardization Complete!")
}

func intPtr(i int) *int {
	return &i
}

func derefInt(i *int) string {
	if i == nil {
		return "Present"
	}
	return fmt.Sprintf("%d", *i)
}

func isCurrent(endYear *int) bool {
	return endYear == nil
}
