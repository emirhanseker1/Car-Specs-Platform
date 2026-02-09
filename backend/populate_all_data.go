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
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	fmt.Println("=== Populating VW Golf ===")
	populateVWGolf(db)

	fmt.Println("\n=== Populating VW Passat ===")
	populateVWPassat(db)

	fmt.Println("\n=== Populating BMW 1 Series ===")
	populateBMW1Series(db)

	fmt.Println("\n=== Populating BMW 3 Series ===")
	populateBMW3Series(db)

	fmt.Println("\n=== Populating BMW 5 Series ===")
	populateBMW5Series(db)

	fmt.Println("\n=== Splitting Audi 8V into 8V1/8V2 ===")
	splitAudi8V(db)

	// Verify
	fmt.Println("\n=== Verification ===")
	rows, _ := db.Query(`
		SELECT b.name, m.name, COUNT(g.id) as gen_count
		FROM brands b
		LEFT JOIN models m ON m.brand_id = b.id
		LEFT JOIN generations g ON g.model_id = m.id
		GROUP BY b.id, m.id
		ORDER BY b.name, m.name
	`)
	defer rows.Close()
	for rows.Next() {
		var brand, model sql.NullString
		var count int
		rows.Scan(&brand, &model, &count)
		if model.Valid {
			fmt.Printf("  %s %s: %d generations\n", brand.String, model.String, count)
		}
	}
}

func getVWBrandID(db *sql.DB) int64 {
	var id int64
	db.QueryRow(`SELECT id FROM brands WHERE name = 'Volkswagen'`).Scan(&id)
	return id
}

func getBMWBrandID(db *sql.DB) int64 {
	var id int64
	db.QueryRow(`SELECT id FROM brands WHERE name = 'BMW'`).Scan(&id)
	return id
}

func populateVWGolf(db *sql.DB) {
	brandID := getVWBrandID(db)

	// Check if Golf model exists
	var modelID int64
	err := db.QueryRow(`SELECT id FROM models WHERE name = 'Golf' AND brand_id = ?`, brandID).Scan(&modelID)
	if err != nil {
		fmt.Println("  Golf model not found, skipping")
		return
	}

	// Check if generations exist
	var count int
	db.QueryRow(`SELECT COUNT(*) FROM generations WHERE model_id = ?`, modelID).Scan(&count)
	if count > 0 {
		fmt.Printf("  Golf already has %d generations\n", count)
		return
	}

	// Insert generations
	generations := []struct {
		name, code, imageURL string
		startYear, endYear   int
		isFacelift           bool
	}{
		{"Golf 8", "MK8", "/images/vehicles/volkswagen/mk8/vw-golf-8.png", 2020, 2024, false},
		{"Golf 7.5 (Makyajlı)", "MK7.5", "/images/vehicles/volkswagen/mk7-5/vw-golf-7-5.png", 2017, 2020, true},
		{"Golf 7", "MK7", "/images/vehicles/volkswagen/mk7/vw-golf-7.png", 2012, 2016, false},
		{"Golf 6", "MK6", "/images/vehicles/volkswagen/mk6/vw-golf-6.png", 2009, 2012, false},
		{"Golf 5", "MK5", "/images/vehicles/volkswagen/mk5/vw-golf-5.png", 2004, 2008, false},
	}

	for _, g := range generations {
		_, err := db.Exec(`INSERT INTO generations (model_id, name, code, start_year, end_year, image_url, is_facelift) VALUES (?, ?, ?, ?, ?, ?, ?)`,
			modelID, g.name, g.code, g.startYear, g.endYear, g.imageURL, g.isFacelift)
		if err != nil {
			fmt.Printf("  ✗ Error inserting %s: %v\n", g.code, err)
		} else {
			fmt.Printf("  ✓ Inserted %s\n", g.code)
		}
	}
}

func populateVWPassat(db *sql.DB) {
	brandID := getVWBrandID(db)

	var modelID int64
	err := db.QueryRow(`SELECT id FROM models WHERE name = 'Passat' AND brand_id = ?`, brandID).Scan(&modelID)
	if err != nil {
		fmt.Println("  Passat model not found, skipping")
		return
	}

	var count int
	db.QueryRow(`SELECT COUNT(*) FROM generations WHERE model_id = ?`, modelID).Scan(&count)
	if count > 0 {
		fmt.Printf("  Passat already has %d generations\n", count)
		return
	}

	generations := []struct {
		name, code, imageURL string
		startYear, endYear   int
		isFacelift           bool
	}{
		{"Passat B8.5 (Makyajlı)", "B8.5", "/images/vehicles/volkswagen/b8-5/vw-passat-b8-5.png", 2019, 2024, true},
		{"Passat B8", "B8", "/images/vehicles/volkswagen/b8/vw-passat-b8.png", 2015, 2019, false},
		{"Passat B7", "B7", "/images/vehicles/volkswagen/b7/vw-passat-b7.png", 2011, 2015, false},
		{"Passat B6", "B6", "/images/vehicles/volkswagen/b6/vw-passat-b6.png", 2005, 2010, false},
		{"Passat B5.5 (Makyajlı)", "B5.5", "/images/vehicles/volkswagen/b5-5/vw-passat-b5-5.png", 2001, 2005, true},
	}

	for _, g := range generations {
		_, err := db.Exec(`INSERT INTO generations (model_id, name, code, start_year, end_year, image_url, is_facelift) VALUES (?, ?, ?, ?, ?, ?, ?)`,
			modelID, g.name, g.code, g.startYear, g.endYear, g.imageURL, g.isFacelift)
		if err != nil {
			fmt.Printf("  ✗ Error inserting %s: %v\n", g.code, err)
		} else {
			fmt.Printf("  ✓ Inserted %s\n", g.code)
		}
	}
}

func populateBMW1Series(db *sql.DB) {
	brandID := getBMWBrandID(db)

	// First ensure model exists
	var modelID int64
	err := db.QueryRow(`SELECT id FROM models WHERE name = '1 Serisi' AND brand_id = ?`, brandID).Scan(&modelID)
	if err != nil {
		// Create the model
		result, err := db.Exec(`INSERT INTO models (brand_id, name, image_url, body_style) VALUES (?, '1 Serisi', '/images/vehicles/bmw/f40/bmw-1-f40.png', 'Hatchback')`, brandID)
		if err != nil {
			fmt.Printf("  ✗ Error creating model: %v\n", err)
			return
		}
		modelID, _ = result.LastInsertId()
		fmt.Println("  ✓ Created 1 Serisi model")
	}

	var count int
	db.QueryRow(`SELECT COUNT(*) FROM generations WHERE model_id = ?`, modelID).Scan(&count)
	if count > 0 {
		fmt.Printf("  1 Serisi already has %d generations\n", count)
		return
	}

	generations := []struct {
		name, code, imageURL string
		startYear, endYear   int
		isFacelift           bool
	}{
		{"1 Serisi F40", "F40", "/images/vehicles/bmw/f40/bmw-1-f40.png", 2019, 2024, false},
		{"1 Serisi F20 (LCI)", "F20 LCI", "/images/vehicles/bmw/f20-lci/bmw-1-f20-lci.png", 2015, 2019, true},
		{"1 Serisi F20", "F20", "/images/vehicles/bmw/f20/bmw-1-f20.png", 2011, 2015, false},
		{"1 Serisi E87 (LCI)", "E87 LCI", "/images/vehicles/bmw/e87-lci/bmw-1-e87-lci.png", 2007, 2011, true},
		{"1 Serisi E87", "E87", "/images/vehicles/bmw/e87/bmw-1-e87.png", 2004, 2007, false},
	}

	for _, g := range generations {
		_, err := db.Exec(`INSERT INTO generations (model_id, name, code, start_year, end_year, image_url, is_facelift) VALUES (?, ?, ?, ?, ?, ?, ?)`,
			modelID, g.name, g.code, g.startYear, g.endYear, g.imageURL, g.isFacelift)
		if err != nil {
			fmt.Printf("  ✗ Error inserting %s: %v\n", g.code, err)
		} else {
			fmt.Printf("  ✓ Inserted %s\n", g.code)
		}
	}
}

func populateBMW3Series(db *sql.DB) {
	brandID := getBMWBrandID(db)

	var modelID int64
	err := db.QueryRow(`SELECT id FROM models WHERE name = '3 Serisi' AND brand_id = ?`, brandID).Scan(&modelID)
	if err != nil {
		result, err := db.Exec(`INSERT INTO models (brand_id, name, image_url, body_style) VALUES (?, '3 Serisi', '/images/vehicles/bmw/g20/bmw-3-g20.png', 'Sedan')`, brandID)
		if err != nil {
			fmt.Printf("  ✗ Error creating model: %v\n", err)
			return
		}
		modelID, _ = result.LastInsertId()
		fmt.Println("  ✓ Created 3 Serisi model")
	}

	var count int
	db.QueryRow(`SELECT COUNT(*) FROM generations WHERE model_id = ?`, modelID).Scan(&count)
	if count > 0 {
		fmt.Printf("  3 Serisi already has %d generations\n", count)
		return
	}

	generations := []struct {
		name, code, imageURL string
		startYear, endYear   int
		isFacelift           bool
	}{
		{"3 Serisi G20 (LCI)", "G20 LCI", "/images/vehicles/bmw/g20-lci/bmw-3-g20-lci.png", 2022, 2024, true},
		{"3 Serisi G20", "G20", "/images/vehicles/bmw/g20/bmw-3-g20.png", 2019, 2022, false},
		{"3 Serisi F30 (LCI)", "F30 LCI", "/images/vehicles/bmw/f30-lci/bmw-3-f30-lci.png", 2015, 2019, true},
		{"3 Serisi F30", "F30", "/images/vehicles/bmw/f30/bmw-3-f30.png", 2012, 2015, false},
		{"3 Serisi E90 (LCI)", "E90 LCI", "/images/vehicles/bmw/e90-lci/bmw-3-e90-lci.png", 2008, 2012, true},
		{"3 Serisi E90", "E90", "/images/vehicles/bmw/e90/bmw-3-e90.png", 2005, 2008, false},
	}

	for _, g := range generations {
		_, err := db.Exec(`INSERT INTO generations (model_id, name, code, start_year, end_year, image_url, is_facelift) VALUES (?, ?, ?, ?, ?, ?, ?)`,
			modelID, g.name, g.code, g.startYear, g.endYear, g.imageURL, g.isFacelift)
		if err != nil {
			fmt.Printf("  ✗ Error inserting %s: %v\n", g.code, err)
		} else {
			fmt.Printf("  ✓ Inserted %s\n", g.code)
		}
	}
}

func populateBMW5Series(db *sql.DB) {
	brandID := getBMWBrandID(db)

	var modelID int64
	err := db.QueryRow(`SELECT id FROM models WHERE name = '5 Serisi' AND brand_id = ?`, brandID).Scan(&modelID)
	if err != nil {
		result, err := db.Exec(`INSERT INTO models (brand_id, name, image_url, body_style) VALUES (?, '5 Serisi', '/images/vehicles/bmw/g60/bmw-5-g60.png', 'Sedan')`, brandID)
		if err != nil {
			fmt.Printf("  ✗ Error creating model: %v\n", err)
			return
		}
		modelID, _ = result.LastInsertId()
		fmt.Println("  ✓ Created 5 Serisi model")
	}

	var count int
	db.QueryRow(`SELECT COUNT(*) FROM generations WHERE model_id = ?`, modelID).Scan(&count)
	if count > 0 {
		fmt.Printf("  5 Serisi already has %d generations\n", count)
		return
	}

	generations := []struct {
		name, code, imageURL string
		startYear, endYear   int
		isFacelift           bool
	}{
		{"5 Serisi G60", "G60", "/images/vehicles/bmw/g60/bmw-5-g60.png", 2023, 2024, false},
		{"5 Serisi G30 (LCI)", "G30 LCI", "/images/vehicles/bmw/g30-lci/bmw-5-g30-lci.png", 2020, 2023, true},
		{"5 Serisi G30", "G30", "/images/vehicles/bmw/g30/bmw-5-g30.png", 2017, 2020, false},
		{"5 Serisi F10 (LCI)", "F10 LCI", "/images/vehicles/bmw/f10-lci/bmw-5-f10-lci.png", 2013, 2017, true},
		{"5 Serisi F10", "F10", "/images/vehicles/bmw/f10/bmw-5-f10.png", 2010, 2013, false},
		{"5 Serisi E60 (LCI)", "E60 LCI", "/images/vehicles/bmw/e60-lci/bmw-5-e60-lci.png", 2007, 2010, true},
		{"5 Serisi E60", "E60", "/images/vehicles/bmw/e60/bmw-5-e60.png", 2003, 2007, false},
	}

	for _, g := range generations {
		_, err := db.Exec(`INSERT INTO generations (model_id, name, code, start_year, end_year, image_url, is_facelift) VALUES (?, ?, ?, ?, ?, ?, ?)`,
			modelID, g.name, g.code, g.startYear, g.endYear, g.imageURL, g.isFacelift)
		if err != nil {
			fmt.Printf("  ✗ Error inserting %s: %v\n", g.code, err)
		} else {
			fmt.Printf("  ✓ Inserted %s\n", g.code)
		}
	}
}

func splitAudi8V(db *sql.DB) {
	// Check if 8V1 already exists
	var count int
	db.QueryRow(`SELECT COUNT(*) FROM generations WHERE code = '8V1'`).Scan(&count)
	if count > 0 {
		fmt.Println("  8V1/8V2 split already done")
		return
	}

	// Get Audi A3 model ID
	var modelID int64
	err := db.QueryRow(`SELECT id FROM models WHERE name = 'A3'`).Scan(&modelID)
	if err != nil {
		fmt.Printf("  A3 model not found: %v\n", err)
		return
	}

	// Update existing 8V to 8V1
	_, err = db.Exec(`UPDATE generations SET code = '8V1', name = 'Tip 8V - 3. Nesil (2013-2016) | Makyajsız', start_year = 2013, end_year = 2016, image_url = '/images/vehicles/audi/8v1/audi-a3-8v-2013-2016.png', is_facelift = 0 WHERE code = '8V' AND model_id = ?`, modelID)
	if err != nil {
		fmt.Printf("  ✗ Error updating 8V to 8V1: %v\n", err)
	} else {
		fmt.Println("  ✓ Updated 8V to 8V1")
	}

	// Insert 8V2
	_, err = db.Exec(`INSERT INTO generations (model_id, name, code, start_year, end_year, image_url, is_facelift) VALUES (?, 'Tip 8V - 3. Nesil (2016-2020) | Makyajlı', '8V2', 2016, 2020, '/images/vehicles/audi/8v2/audi-a3-8v-2016-2020.png', 1)`, modelID)
	if err != nil {
		fmt.Printf("  ✗ Error inserting 8V2: %v\n", err)
	} else {
		fmt.Println("  ✓ Inserted 8V2")
	}
}
