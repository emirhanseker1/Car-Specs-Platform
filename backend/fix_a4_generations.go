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

	fmt.Println("=== FIXING AUDI A4 GENERATION NAMES ===\n")

	// First, delete all existing A4 generations and their trims
	// We'll recreate them with correct naming

	// Get A4 model ID
	var a4ID int64
	err = db.QueryRow("SELECT id FROM models WHERE name = 'A4' AND brand_id = 1").Scan(&a4ID)
	if err != nil {
		log.Fatal("A4 model not found:", err)
	}
	fmt.Printf("A4 model ID: %d\n", a4ID)

	// Delete existing trims for A4
	res, err := db.Exec("DELETE FROM trims WHERE model_id = ?", a4ID)
	if err != nil {
		log.Fatal("Failed to delete trims:", err)
	}
	n, _ := res.RowsAffected()
	fmt.Printf("Deleted %d trims\n", n)

	// Delete existing generations for A4
	res, err = db.Exec("DELETE FROM generations WHERE model_id = ?", a4ID)
	if err != nil {
		log.Fatal("Failed to delete generations:", err)
	}
	n, _ = res.RowsAffected()
	fmt.Printf("Deleted %d generations\n", n)

	// Add correct generations
	fmt.Println("\nAdding correct generations...")

	generations := []struct {
		code       string
		name       string
		startYear  int
		endYear    *int
		isFacelift bool
	}{
		{"B9.5", "Audi A4 B9.5 (2020-Günümüz)", 2020, nil, true},
		{"B9", "Audi A4 B9 (2016-2020)", 2016, intPtr(2020), false},
		{"B8", "Audi A4 B8 (2008-2016)", 2008, intPtr(2016), false},
		{"B7", "Audi A4 B7 (2004-2008)", 2004, intPtr(2008), false},
		{"B6", "Audi A4 B6 (2001-2004)", 2001, intPtr(2004), false},
		{"B5", "Audi A4 B5 (1995-2001)", 1995, intPtr(2001), false},
	}

	genIDs := make(map[string]int64)
	for _, g := range generations {
		var endYear interface{} = nil
		if g.endYear != nil {
			endYear = *g.endYear
		}
		res, err := db.Exec(`
			INSERT INTO generations (model_id, code, name, start_year, end_year, is_facelift)
			VALUES (?, ?, ?, ?, ?, ?)
		`, a4ID, g.code, g.name, g.startYear, endYear, g.isFacelift)
		if err != nil {
			fmt.Printf("  ✗ %s: %v\n", g.code, err)
			continue
		}
		id, _ := res.LastInsertId()
		genIDs[g.code] = id
		fmt.Printf("  ✓ %s (ID: %d)\n", g.code, id)
	}

	// Add trims
	fmt.Println("\nAdding trims...")

	// B9.5 (current)
	if id, ok := genIDs["B9.5"]; ok {
		addTrims(db, id, a4ID, 2020, []trimData{
			{"30 TFSI", 150, 250, "S tronic", "Benzin", "Önden Çekiş", 8.9},
			{"35 TFSI", 150, 250, "S tronic", "Benzin", "Önden Çekiş", 8.4},
			{"40 TFSI Quattro", 204, 320, "S tronic", "Benzin", "4WD", 7.1},
			{"45 TFSI Quattro", 265, 370, "S tronic", "Benzin", "4WD", 5.8},
			{"S4 3.0 TDI Quattro", 347, 700, "Tiptronic", "Dizel", "4WD", 4.9},
			{"30 TDI", 136, 300, "S tronic", "Dizel", "Önden Çekiş", 9.6},
			{"35 TDI", 163, 380, "S tronic", "Dizel", "Önden Çekiş", 8.2},
			{"40 TDI Quattro", 204, 400, "S tronic", "Dizel", "4WD", 7.0},
		})
	}

	// B9
	if id, ok := genIDs["B9"]; ok {
		addTrims(db, id, a4ID, 2016, []trimData{
			{"1.4 TFSI", 150, 250, "S tronic", "Benzin", "Önden Çekiş", 8.5},
			{"2.0 TFSI", 190, 320, "S tronic", "Benzin", "Önden Çekiş", 7.3},
			{"2.0 TFSI Quattro", 252, 370, "S tronic", "Benzin", "4WD", 5.8},
			{"S4 3.0 TFSI Quattro", 354, 500, "Tiptronic", "Benzin", "4WD", 4.7},
			{"2.0 TDI", 150, 320, "S tronic", "Dizel", "Önden Çekiş", 8.4},
			{"2.0 TDI Quattro", 190, 400, "S tronic", "Dizel", "4WD", 7.4},
			{"3.0 TDI Quattro", 272, 600, "Tiptronic", "Dizel", "4WD", 5.3},
		})
	}

	// B8
	if id, ok := genIDs["B8"]; ok {
		addTrims(db, id, a4ID, 2008, []trimData{
			{"1.8 TFSI", 160, 250, "Multitronic", "Benzin", "Önden Çekiş", 8.4},
			{"2.0 TFSI", 211, 350, "S tronic", "Benzin", "Önden Çekiş", 7.0},
			{"2.0 TFSI Quattro", 211, 350, "S tronic", "Benzin", "4WD", 6.9},
			{"S4 3.0 TFSI Quattro", 333, 440, "S tronic", "Benzin", "4WD", 5.1},
			{"2.0 TDI", 143, 320, "Multitronic", "Dizel", "Önden Çekiş", 8.8},
			{"2.0 TDI Quattro", 177, 380, "S tronic", "Dizel", "4WD", 7.5},
			{"3.0 TDI Quattro", 245, 500, "S tronic", "Dizel", "4WD", 5.9},
		})
	}

	// B7
	if id, ok := genIDs["B7"]; ok {
		addTrims(db, id, a4ID, 2004, []trimData{
			{"1.6", 102, 148, "Multitronic", "Benzin", "Önden Çekiş", 11.9},
			{"1.8T", 163, 225, "Multitronic", "Benzin", "Önden Çekiş", 8.5},
			{"2.0 TFSI", 200, 280, "Multitronic", "Benzin", "Önden Çekiş", 7.5},
			{"2.0 TFSI Quattro", 200, 280, "Tiptronic", "Benzin", "4WD", 7.4},
			{"S4 4.2 V8 Quattro", 344, 410, "Tiptronic", "Benzin", "4WD", 5.6},
			{"1.9 TDI", 115, 285, "Multitronic", "Dizel", "Önden Çekiş", 10.4},
			{"2.0 TDI", 140, 320, "Multitronic", "Dizel", "Önden Çekiş", 9.5},
		})
	}

	// B6
	if id, ok := genIDs["B6"]; ok {
		addTrims(db, id, a4ID, 2001, []trimData{
			{"1.6", 102, 148, "Multitronic", "Benzin", "Önden Çekiş", 12.0},
			{"1.8T", 150, 210, "Multitronic", "Benzin", "Önden Çekiş", 9.1},
			{"2.0", 130, 195, "Multitronic", "Benzin", "Önden Çekiş", 10.2},
			{"3.0 V6 Quattro", 220, 300, "Tiptronic", "Benzin", "4WD", 7.1},
			{"S4 4.2 V8 Quattro", 344, 410, "Tiptronic", "Benzin", "4WD", 5.6},
			{"1.9 TDI", 130, 310, "Multitronic", "Dizel", "Önden Çekiş", 9.8},
		})
	}

	// B5
	if id, ok := genIDs["B5"]; ok {
		addTrims(db, id, a4ID, 1995, []trimData{
			{"1.6", 101, 140, "Otomatik", "Benzin", "Önden Çekiş", 12.5},
			{"1.8", 125, 173, "Otomatik", "Benzin", "Önden Çekiş", 10.5},
			{"1.8T", 150, 210, "Tiptronic", "Benzin", "Önden Çekiş", 8.9},
			{"1.8T Quattro", 180, 235, "Tiptronic", "Benzin", "4WD", 7.8},
			{"2.8 V6 Quattro", 193, 280, "Tiptronic", "Benzin", "4WD", 7.7},
			{"S4 2.7 BiTurbo Quattro", 265, 400, "Tiptronic", "Benzin", "4WD", 5.9},
			{"1.9 TDI", 110, 235, "Otomatik", "Dizel", "Önden Çekiş", 11.1},
		})
	}

	// Verify
	fmt.Println("\n=== VERIFICATION ===")
	rows, _ := db.Query(`
		SELECT g.code, g.name, g.start_year, COUNT(t.id) as trim_count
		FROM generations g 
		LEFT JOIN trims t ON t.generation_id = g.id 
		WHERE g.model_id = ?
		GROUP BY g.id
		ORDER BY g.start_year DESC
	`, a4ID)
	for rows.Next() {
		var code, name string
		var startYear, count int
		rows.Scan(&code, &name, &startYear, &count)
		fmt.Printf("  %s (%d): %d trims\n", code, startYear, count)
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

func addTrims(db *sql.DB, genID, modelID int64, year int, trims []trimData) {
	for _, t := range trims {
		_, err := db.Exec(`
			INSERT INTO trims (generation_id, model_id, name, power_hp, torque_nm, transmission_type, fuel_type, drivetrain, acceleration_0_100, year, seating_capacity, currency, market)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 5, 'TRY', 'TR')
		`, genID, modelID, t.name, t.powerHP, t.torqueNM, t.transmission, t.fuel, t.drivetrain, t.accel, year)
		if err != nil {
			fmt.Printf("    ✗ %s: %v\n", t.name, err)
		}
	}
}

func intPtr(i int) *int {
	return &i
}
