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

	fmt.Println("=== ADDING AUDI A4 DATA ===\n")

	// Get Audi brand ID
	var audiID int
	err = db.QueryRow("SELECT id FROM brands WHERE name = 'Audi'").Scan(&audiID)
	if err != nil {
		log.Fatal("Audi brand not found:", err)
	}
	fmt.Printf("Audi brand ID: %d\n", audiID)

	// Create A4 model
	res, err := db.Exec(`INSERT INTO models (brand_id, name, body_style) VALUES (?, 'A4', 'Sedan')`, audiID)
	if err != nil {
		log.Fatal("Failed to create A4 model:", err)
	}
	a4ID, _ := res.LastInsertId()
	fmt.Printf("Created Audi A4 model with ID: %d\n", a4ID)

	// Add A4 generations (B5 to B9)
	generations := []struct {
		code       string
		name       string
		startYear  int
		endYear    *int
		isFacelift bool
		imageURL   string
	}{
		{"B9", "Audi A4 B9 (2020-Günümüz)", 2020, nil, true, "/images/vehicles/audi/a4/b9-facelift.png"},
		{"B9-Pre", "Audi A4 B9 (2016-2020)", 2016, intPtr(2020), false, "/images/vehicles/audi/a4/b9.png"},
		{"B8.5", "Audi A4 B8.5 (Makyajlı) (2012-2016)", 2012, intPtr(2016), true, "/images/vehicles/audi/a4/b8-facelift.png"},
		{"B8", "Audi A4 B8 (2008-2012)", 2008, intPtr(2012), false, "/images/vehicles/audi/a4/b8.png"},
		{"B7", "Audi A4 B7 (2004-2008)", 2004, intPtr(2008), false, "/images/vehicles/audi/a4/b7.png"},
		{"B6", "Audi A4 B6 (2001-2004)", 2001, intPtr(2004), false, "/images/vehicles/audi/a4/b6.png"},
		{"B5", "Audi A4 B5 (1994-2001)", 1994, intPtr(2001), false, "/images/vehicles/audi/a4/b5.png"},
	}

	genIDs := make(map[string]int64)
	for _, g := range generations {
		var endYear interface{} = nil
		if g.endYear != nil {
			endYear = *g.endYear
		}
		res, err := db.Exec(`
			INSERT INTO generations (model_id, code, name, start_year, end_year, is_facelift, image_url)
			VALUES (?, ?, ?, ?, ?, ?, ?)
		`, a4ID, g.code, g.name, g.startYear, endYear, g.isFacelift, g.imageURL)
		if err != nil {
			fmt.Printf("  ✗ %s: %v\n", g.code, err)
			continue
		}
		id, _ := res.LastInsertId()
		genIDs[g.code] = id
		fmt.Printf("  ✓ %s (ID: %d)\n", g.code, id)
	}

	// Add trims for each generation
	fmt.Println("\nAdding trims...")

	// B9 (current gen)
	if id, ok := genIDs["B9"]; ok {
		addA4Trims(db, id, a4ID, []trimData{
			{"30 TFSI", 150, 250, "S tronic", "Benzin", "Önden Çekiş", 8.9},
			{"35 TFSI", 150, 250, "S tronic", "Benzin", "Önden Çekiş", 8.4},
			{"40 TFSI Quattro", 190, 320, "S tronic", "Benzin", "4WD", 7.3},
			{"45 TFSI Quattro", 265, 370, "S tronic", "Benzin", "4WD", 5.8},
			{"S4 3.0 TFSI Quattro", 354, 500, "Tiptronic", "Benzin", "4WD", 4.7},
			{"30 TDI", 136, 300, "S tronic", "Dizel", "Önden Çekiş", 9.6},
			{"35 TDI", 163, 380, "S tronic", "Dizel", "Önden Çekiş", 8.5},
			{"40 TDI Quattro", 190, 400, "S tronic", "Dizel", "4WD", 7.4},
		})
	}

	// B9-Pre
	if id, ok := genIDs["B9-Pre"]; ok {
		addA4Trims(db, id, a4ID, []trimData{
			{"1.4 TFSI", 150, 250, "S tronic", "Benzin", "Önden Çekiş", 8.5},
			{"2.0 TFSI", 190, 320, "S tronic", "Benzin", "Önden Çekiş", 7.3},
			{"2.0 TFSI Quattro", 252, 370, "S tronic", "Benzin", "4WD", 5.8},
			{"S4 3.0 TFSI Quattro", 354, 500, "Tiptronic", "Benzin", "4WD", 4.7},
			{"2.0 TDI", 150, 320, "S tronic", "Dizel", "Önden Çekiş", 8.4},
			{"2.0 TDI Quattro", 190, 400, "S tronic", "Dizel", "4WD", 7.4},
			{"3.0 TDI Quattro", 272, 600, "Tiptronic", "Dizel", "4WD", 5.3},
		})
	}

	// B8.5
	if id, ok := genIDs["B8.5"]; ok {
		addA4Trims(db, id, a4ID, []trimData{
			{"1.8 TFSI", 170, 320, "Multitronic", "Benzin", "Önden Çekiş", 7.9},
			{"2.0 TFSI", 211, 350, "S tronic", "Benzin", "Önden Çekiş", 7.0},
			{"2.0 TFSI Quattro", 225, 350, "S tronic", "Benzin", "4WD", 6.4},
			{"S4 3.0 TFSI Quattro", 333, 440, "S tronic", "Benzin", "4WD", 5.0},
			{"2.0 TDI", 143, 320, "Multitronic", "Dizel", "Önden Çekiş", 8.6},
			{"2.0 TDI Quattro", 177, 380, "S tronic", "Dizel", "4WD", 7.5},
			{"3.0 TDI Quattro", 245, 500, "S tronic", "Dizel", "4WD", 5.9},
		})
	}

	// B8
	if id, ok := genIDs["B8"]; ok {
		addA4Trims(db, id, a4ID, []trimData{
			{"1.8 TFSI", 160, 250, "Multitronic", "Benzin", "Önden Çekiş", 8.4},
			{"2.0 TFSI", 211, 350, "S tronic", "Benzin", "Önden Çekiş", 7.0},
			{"2.0 TFSI Quattro", 211, 350, "S tronic", "Benzin", "4WD", 6.9},
			{"S4 3.0 TFSI Quattro", 333, 440, "S tronic", "Benzin", "4WD", 5.1},
			{"2.0 TDI", 143, 320, "Multitronic", "Dizel", "Önden Çekiş", 8.8},
			{"2.0 TDI Quattro", 170, 350, "S tronic", "Dizel", "4WD", 7.8},
			{"3.0 TDI Quattro", 240, 500, "S tronic", "Dizel", "4WD", 6.1},
		})
	}

	// B7
	if id, ok := genIDs["B7"]; ok {
		addA4Trims(db, id, a4ID, []trimData{
			{"1.6", 102, 148, "Multitronic", "Benzin", "Önden Çekiş", 11.9},
			{"1.8T", 163, 225, "Multitronic", "Benzin", "Önden Çekiş", 8.5},
			{"2.0 TFSI", 200, 280, "Multitronic", "Benzin", "Önden Çekiş", 7.5},
			{"2.0 TFSI Quattro", 200, 280, "Tiptronic", "Benzin", "4WD", 7.4},
			{"S4 4.2 V8 Quattro", 344, 410, "Tiptronic", "Benzin", "4WD", 5.6},
			{"1.9 TDI", 115, 285, "Multitronic", "Dizel", "Önden Çekiş", 10.4},
			{"2.0 TDI", 140, 320, "Multitronic", "Dizel", "Önden Çekiş", 9.5},
			{"2.5 TDI Quattro", 180, 370, "Tiptronic", "Dizel", "4WD", 8.1},
		})
	}

	// B6
	if id, ok := genIDs["B6"]; ok {
		addA4Trims(db, id, a4ID, []trimData{
			{"1.6", 102, 148, "Multitronic", "Benzin", "Önden Çekiş", 12.0},
			{"1.8T", 150, 210, "Multitronic", "Benzin", "Önden Çekiş", 9.1},
			{"2.0", 130, 195, "Multitronic", "Benzin", "Önden Çekiş", 10.2},
			{"2.0 FSI", 150, 200, "Multitronic", "Benzin", "Önden Çekiş", 9.2},
			{"3.0 V6 Quattro", 220, 300, "Tiptronic", "Benzin", "4WD", 7.1},
			{"S4 4.2 V8 Quattro", 344, 410, "Tiptronic", "Benzin", "4WD", 5.6},
			{"1.9 TDI", 130, 310, "Multitronic", "Dizel", "Önden Çekiş", 9.8},
			{"2.5 TDI Quattro", 180, 370, "Tiptronic", "Dizel", "4WD", 8.1},
		})
	}

	// B5
	if id, ok := genIDs["B5"]; ok {
		addA4Trims(db, id, a4ID, []trimData{
			{"1.6", 101, 140, "Otomatik", "Benzin", "Önden Çekiş", 12.5},
			{"1.8", 125, 173, "Otomatik", "Benzin", "Önden Çekiş", 10.5},
			{"1.8T", 150, 210, "Otomatik", "Benzin", "Önden Çekiş", 8.9},
			{"1.8T Quattro", 180, 235, "Tiptronic", "Benzin", "4WD", 7.8},
			{"2.8 V6 Quattro", 193, 280, "Tiptronic", "Benzin", "4WD", 7.7},
			{"S4 2.7 BiTurbo Quattro", 265, 400, "Tiptronic", "Benzin", "4WD", 5.9},
			{"1.9 TDI", 110, 235, "Otomatik", "Dizel", "Önden Çekiş", 11.1},
			{"2.5 TDI Quattro", 150, 310, "Tiptronic", "Dizel", "4WD", 9.0},
		})
	}

	// Verify
	fmt.Println("\n=== VERIFICATION ===")
	var totalTrims int
	db.QueryRow("SELECT COUNT(*) FROM trims WHERE model_id = ?", a4ID).Scan(&totalTrims)
	fmt.Printf("Total Audi A4 trims: %d\n", totalTrims)

	rows, _ := db.Query(`
		SELECT g.code, COUNT(t.id) 
		FROM generations g 
		LEFT JOIN trims t ON t.generation_id = g.id 
		WHERE g.model_id = ?
		GROUP BY g.id
		ORDER BY g.start_year DESC
	`, a4ID)
	for rows.Next() {
		var code string
		var count int
		rows.Scan(&code, &count)
		fmt.Printf("  %s: %d trims\n", code, count)
	}
	rows.Close()
}

type trimData struct {
	name         string
	powerHP      int
	torqueNM     int
	transmission string
	fuel         string
	drivetrain   string
	accel        float64
}

func addA4Trims(db *sql.DB, genID, modelID int64, trims []trimData) {
	for _, t := range trims {
		_, err := db.Exec(`
			INSERT INTO trims (generation_id, model_id, name, power_hp, torque_nm, transmission_type, fuel_type, drivetrain, acceleration_0_100, year, seating_capacity, currency, market)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, 2020, 5, 'TRY', 'TR')
		`, genID, modelID, t.name, t.powerHP, t.torqueNM, t.transmission, t.fuel, t.drivetrain, t.accel)
		if err != nil {
			fmt.Printf("    ✗ %s: %v\n", t.name, err)
		} else {
			fmt.Printf("    ✓ %s\n", t.name)
		}
	}
}

func intPtr(i int) *int {
	return &i
}
