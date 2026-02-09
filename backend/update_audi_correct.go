package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

// Trim represents a vehicle trim with all specifications
type Trim struct {
	Name             string
	PowerHP          int
	TorqueNM         int
	TransmissionType string
	TransmissionCode string
	FuelType         string
	Drivetrain       string
	Acceleration     float64
	EngineType       string
	Cylinders        int
	Year             int
	StartYear        int
	EndYear          *int
}

func main() {
	db, err := sql.Open("sqlite", "./vehicles.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	fmt.Println("=== UPDATING AUDI DATA WITH CORRECT SPECIFICATIONS ===\n")

	// Get Audi brand ID
	var audiID int64
	err = db.QueryRow("SELECT id FROM brands WHERE name = 'Audi'").Scan(&audiID)
	if err != nil {
		log.Fatal("Audi not found:", err)
	}

	// =========================================
	// AUDI A4
	// =========================================
	fmt.Println("=== AUDI A4 ===")

	var a4ID int64
	err = db.QueryRow("SELECT id FROM models WHERE name = 'A4' AND brand_id = ?", audiID).Scan(&a4ID)
	if err != nil {
		log.Fatal("A4 not found:", err)
	}

	// Clear existing A4 data
	db.Exec("DELETE FROM trims WHERE model_id = ?", a4ID)
	db.Exec("DELETE FROM generations WHERE model_id = ?", a4ID)

	// A4 Generations with correct naming
	a4Gens := []struct {
		code       string
		name       string
		startYear  int
		endYear    *int
		isFacelift bool
	}{
		{"B9", "Audi A4 B9 (2016-Günümüz)", 2016, nil, false},
		{"B8", "Audi A4 B8 (2008-2016)", 2008, intPtr(2016), false},
		{"B7", "Audi A4 B7 (2004-2008)", 2004, intPtr(2008), false},
		{"B6", "Audi A4 B6 (2001-2004)", 2001, intPtr(2004), false},
		{"B5", "Audi A4 B5 (1995-2001)", 1995, intPtr(2001), false},
	}

	a4GenIDs := make(map[string]int64)
	for _, g := range a4Gens {
		res, _ := db.Exec(`INSERT INTO generations (model_id, code, name, start_year, end_year, is_facelift) VALUES (?, ?, ?, ?, ?, ?)`,
			a4ID, g.code, g.name, g.startYear, g.endYear, g.isFacelift)
		id, _ := res.LastInsertId()
		a4GenIDs[g.code] = id
		fmt.Printf("  ✓ %s (ID: %d)\n", g.code, id)
	}

	// A4 B9 Trims (2016-Günümüz)
	if id, ok := a4GenIDs["B9"]; ok {
		insertTrims(db, id, a4ID, []Trim{
			// Pre-facelift (2016-2019)
			{Name: "1.4 TFSI", PowerHP: 150, TorqueNM: 250, TransmissionType: "S tronic", TransmissionCode: "DL382", FuelType: "Benzin", Drivetrain: "Önden Çekiş", Acceleration: 8.5, EngineType: "Turbo", Cylinders: 4, StartYear: 2016, EndYear: intPtr(2019)},
			{Name: "2.0 TDI", PowerHP: 190, TorqueNM: 400, TransmissionType: "S tronic", TransmissionCode: "DL382", FuelType: "Dizel", Drivetrain: "Önden Çekiş", Acceleration: 7.5, EngineType: "Turbo Dizel", Cylinders: 4, StartYear: 2016},
			{Name: "2.0 TDI quattro", PowerHP: 190, TorqueNM: 400, TransmissionType: "S tronic", TransmissionCode: "DL382", FuelType: "Dizel", Drivetrain: "quattro", Acceleration: 7.2, EngineType: "Turbo Dizel", Cylinders: 4, StartYear: 2016},
			// Post-facelift (2020+) - Mild Hybrid
			{Name: "40 TDI quattro", PowerHP: 204, TorqueNM: 400, TransmissionType: "S tronic", TransmissionCode: "DL382+", FuelType: "Dizel + MHEV", Drivetrain: "quattro", Acceleration: 6.9, EngineType: "Turbo Dizel + Mild Hybrid", Cylinders: 4, StartYear: 2020},
			{Name: "45 TFSI quattro", PowerHP: 265, TorqueNM: 370, TransmissionType: "S tronic", TransmissionCode: "DL382+", FuelType: "Benzin + MHEV", Drivetrain: "quattro", Acceleration: 5.5, EngineType: "Turbo + Mild Hybrid", Cylinders: 4, StartYear: 2020},
		})
	}

	// A4 B8 Trims (2008-2016)
	if id, ok := a4GenIDs["B8"]; ok {
		insertTrims(db, id, a4ID, []Trim{
			{Name: "1.8 TFSI", PowerHP: 160, TorqueNM: 250, TransmissionType: "Multitronic", TransmissionCode: "VL381", FuelType: "Benzin", Drivetrain: "Önden Çekiş", Acceleration: 8.3, EngineType: "Turbo", Cylinders: 4, StartYear: 2008, EndYear: intPtr(2012)},
			{Name: "1.8 TFSI (Makyajlı)", PowerHP: 170, TorqueNM: 320, TransmissionType: "Multitronic", TransmissionCode: "VL381", FuelType: "Benzin", Drivetrain: "Önden Çekiş", Acceleration: 8.3, EngineType: "Turbo", Cylinders: 4, StartYear: 2012, EndYear: intPtr(2016)},
			{Name: "2.0 TDI", PowerHP: 143, TorqueNM: 320, TransmissionType: "Multitronic", TransmissionCode: "VL381", FuelType: "Dizel", Drivetrain: "Önden Çekiş", Acceleration: 9.1, EngineType: "Turbo Dizel", Cylinders: 4, StartYear: 2008, EndYear: intPtr(2012)},
			{Name: "2.0 TDI (Makyajlı)", PowerHP: 150, TorqueNM: 320, TransmissionType: "Multitronic", TransmissionCode: "VL381", FuelType: "Dizel", Drivetrain: "Önden Çekiş", Acceleration: 9.1, EngineType: "Turbo Dizel", Cylinders: 4, StartYear: 2012, EndYear: intPtr(2016)},
			{Name: "2.0 TDI 177 PS", PowerHP: 177, TorqueNM: 380, TransmissionType: "Multitronic", TransmissionCode: "VL381", FuelType: "Dizel", Drivetrain: "Önden Çekiş", Acceleration: 8.0, EngineType: "Turbo Dizel", Cylinders: 4, StartYear: 2012, EndYear: intPtr(2016)},
			{Name: "2.0 TFSI quattro", PowerHP: 211, TorqueNM: 350, TransmissionType: "S tronic", TransmissionCode: "DL501", FuelType: "Benzin", Drivetrain: "quattro", Acceleration: 6.5, EngineType: "Turbo", Cylinders: 4, StartYear: 2008, EndYear: intPtr(2012)},
			{Name: "2.0 TFSI quattro (Makyajlı)", PowerHP: 225, TorqueNM: 350, TransmissionType: "S tronic", TransmissionCode: "DL501", FuelType: "Benzin", Drivetrain: "quattro", Acceleration: 6.5, EngineType: "Turbo", Cylinders: 4, StartYear: 2012, EndYear: intPtr(2016)},
		})
	}

	// A4 B7 Trims (2004-2008)
	if id, ok := a4GenIDs["B7"]; ok {
		insertTrims(db, id, a4ID, []Trim{
			{Name: "1.6", PowerHP: 102, TorqueNM: 148, TransmissionType: "Tiptronic", TransmissionCode: "", FuelType: "Benzin", Drivetrain: "Önden Çekiş", Acceleration: 12.6, EngineType: "Atmosferik", Cylinders: 4, StartYear: 2004},
			{Name: "1.8 T", PowerHP: 163, TorqueNM: 225, TransmissionType: "Multitronic", TransmissionCode: "01J/VL300", FuelType: "Benzin", Drivetrain: "Önden Çekiş", Acceleration: 8.6, EngineType: "Turbo", Cylinders: 4, StartYear: 2004},
			{Name: "2.0 TDI", PowerHP: 140, TorqueNM: 320, TransmissionType: "Multitronic", TransmissionCode: "01J/VL300", FuelType: "Dizel", Drivetrain: "Önden Çekiş", Acceleration: 9.7, EngineType: "Turbo Dizel (PD)", Cylinders: 4, StartYear: 2004},
			{Name: "2.0 TFSI quattro", PowerHP: 200, TorqueNM: 280, TransmissionType: "Tiptronic", TransmissionCode: "ZF 6HP", FuelType: "Benzin", Drivetrain: "quattro", Acceleration: 7.7, EngineType: "Turbo", Cylinders: 4, StartYear: 2004},
		})
	}

	// =========================================
	// AUDI A3
	// =========================================
	fmt.Println("\n=== AUDI A3 ===")

	var a3ID int64
	err = db.QueryRow("SELECT id FROM models WHERE name = 'A3' AND brand_id = ?", audiID).Scan(&a3ID)
	if err != nil {
		log.Fatal("A3 not found:", err)
	}

	// Clear existing A3 data
	db.Exec("DELETE FROM trims WHERE model_id = ?", a3ID)
	db.Exec("DELETE FROM generations WHERE model_id = ?", a3ID)

	// A3 Generations
	a3Gens := []struct {
		code       string
		name       string
		startYear  int
		endYear    *int
		isFacelift bool
	}{
		{"8Y", "Audi A3 8Y (2020-Günümüz)", 2020, nil, false},
		{"8V.5", "Audi A3 8V Makyajlı (2016-2020)", 2016, intPtr(2020), true},
		{"8V", "Audi A3 8V (2012-2016)", 2012, intPtr(2016), false},
		{"8P.5", "Audi A3 8P Makyajlı (2008-2012)", 2008, intPtr(2012), true},
		{"8P", "Audi A3 8P (2003-2008)", 2003, intPtr(2008), false},
	}

	a3GenIDs := make(map[string]int64)
	for _, g := range a3Gens {
		res, _ := db.Exec(`INSERT INTO generations (model_id, code, name, start_year, end_year, is_facelift) VALUES (?, ?, ?, ?, ?, ?)`,
			a3ID, g.code, g.name, g.startYear, g.endYear, g.isFacelift)
		id, _ := res.LastInsertId()
		a3GenIDs[g.code] = id
		fmt.Printf("  ✓ %s (ID: %d)\n", g.code, id)
	}

	// A3 8Y (2020-Günümüz)
	if id, ok := a3GenIDs["8Y"]; ok {
		insertTrims(db, id, a3ID, []Trim{
			{Name: "30 TFSI (1.0)", PowerHP: 110, TorqueNM: 200, TransmissionType: "S tronic", TransmissionCode: "DQ200", FuelType: "Benzin", Drivetrain: "Önden Çekiş", Acceleration: 10.6, EngineType: "Turbo", Cylinders: 3, StartYear: 2020},
			{Name: "35 TFSI (1.5)", PowerHP: 150, TorqueNM: 250, TransmissionType: "S tronic", TransmissionCode: "DQ200", FuelType: "Benzin", Drivetrain: "Önden Çekiş", Acceleration: 8.4, EngineType: "Turbo", Cylinders: 4, StartYear: 2020},
			{Name: "30 TFSI (1.5) Makyajlı", PowerHP: 116, TorqueNM: 220, TransmissionType: "S tronic", TransmissionCode: "DQ200", FuelType: "Benzin", Drivetrain: "Önden Çekiş", Acceleration: 9.9, EngineType: "Turbo", Cylinders: 4, StartYear: 2024},
			{Name: "S3", PowerHP: 310, TorqueNM: 400, TransmissionType: "S tronic", TransmissionCode: "DQ381", FuelType: "Benzin", Drivetrain: "quattro", Acceleration: 4.8, EngineType: "Turbo", Cylinders: 4, StartYear: 2020},
			{Name: "S3 Makyajlı", PowerHP: 333, TorqueNM: 420, TransmissionType: "S tronic", TransmissionCode: "DQ381", FuelType: "Benzin", Drivetrain: "quattro", Acceleration: 4.8, EngineType: "Turbo", Cylinders: 4, StartYear: 2024},
			{Name: "RS3", PowerHP: 400, TorqueNM: 500, TransmissionType: "S tronic", TransmissionCode: "DQ500", FuelType: "Benzin", Drivetrain: "quattro", Acceleration: 3.8, EngineType: "Turbo", Cylinders: 5, StartYear: 2021},
		})
	}

	// A3 8V.5 Makyajlı (2016-2020)
	if id, ok := a3GenIDs["8V.5"]; ok {
		insertTrims(db, id, a3ID, []Trim{
			{Name: "1.0 TFSI (30 TFSI)", PowerHP: 116, TorqueNM: 200, TransmissionType: "S tronic", TransmissionCode: "DQ200", FuelType: "Benzin", Drivetrain: "Önden Çekiş", Acceleration: 9.9, EngineType: "Turbo", Cylinders: 3, StartYear: 2016},
			{Name: "1.4 TFSI COD", PowerHP: 150, TorqueNM: 250, TransmissionType: "S tronic", TransmissionCode: "DQ200", FuelType: "Benzin", Drivetrain: "Önden Çekiş", Acceleration: 8.2, EngineType: "Turbo (CoD)", Cylinders: 4, StartYear: 2016},
			{Name: "1.5 TFSI (35 TFSI)", PowerHP: 150, TorqueNM: 250, TransmissionType: "S tronic", TransmissionCode: "DQ200", FuelType: "Benzin", Drivetrain: "Önden Çekiş", Acceleration: 8.2, EngineType: "Turbo (CoD)", Cylinders: 4, StartYear: 2017},
			{Name: "1.6 TDI (30 TDI)", PowerHP: 116, TorqueNM: 250, TransmissionType: "S tronic", TransmissionCode: "DQ200", FuelType: "Dizel", Drivetrain: "Önden Çekiş", Acceleration: 10.4, EngineType: "Turbo Dizel", Cylinders: 4, StartYear: 2017},
			{Name: "S3", PowerHP: 310, TorqueNM: 400, TransmissionType: "S tronic", TransmissionCode: "DQ381", FuelType: "Benzin", Drivetrain: "quattro", Acceleration: 4.6, EngineType: "Turbo", Cylinders: 4, StartYear: 2017},
		})
	}

	// A3 8V (2012-2016)
	if id, ok := a3GenIDs["8V"]; ok {
		insertTrims(db, id, a3ID, []Trim{
			{Name: "1.2 TFSI", PowerHP: 110, TorqueNM: 175, TransmissionType: "S tronic", TransmissionCode: "DQ200", FuelType: "Benzin", Drivetrain: "Önden Çekiş", Acceleration: 10.3, EngineType: "Turbo", Cylinders: 4, StartYear: 2012},
			{Name: "1.4 TFSI", PowerHP: 125, TorqueNM: 200, TransmissionType: "S tronic", TransmissionCode: "DQ200", FuelType: "Benzin", Drivetrain: "Önden Çekiş", Acceleration: 9.3, EngineType: "Turbo", Cylinders: 4, StartYear: 2012},
			{Name: "1.4 TFSI COD", PowerHP: 140, TorqueNM: 250, TransmissionType: "S tronic", TransmissionCode: "DQ200", FuelType: "Benzin", Drivetrain: "Önden Çekiş", Acceleration: 8.4, EngineType: "Turbo (CoD)", Cylinders: 4, StartYear: 2013},
			{Name: "1.6 TDI", PowerHP: 110, TorqueNM: 250, TransmissionType: "S tronic", TransmissionCode: "DQ200", FuelType: "Dizel", Drivetrain: "Önden Çekiş", Acceleration: 10.7, EngineType: "Turbo Dizel", Cylinders: 4, StartYear: 2012},
			{Name: "S3", PowerHP: 300, TorqueNM: 380, TransmissionType: "S tronic", TransmissionCode: "DQ250", FuelType: "Benzin", Drivetrain: "quattro", Acceleration: 5.0, EngineType: "Turbo", Cylinders: 4, StartYear: 2013},
		})
	}

	// A3 8P.5 Makyajlı (2008-2012)
	if id, ok := a3GenIDs["8P.5"]; ok {
		insertTrims(db, id, a3ID, []Trim{
			{Name: "1.4 TFSI", PowerHP: 125, TorqueNM: 200, TransmissionType: "S tronic", TransmissionCode: "DQ200", FuelType: "Benzin", Drivetrain: "Önden Çekiş", Acceleration: 9.3, EngineType: "Turbo", Cylinders: 4, StartYear: 2008},
			{Name: "1.2 TFSI", PowerHP: 105, TorqueNM: 175, TransmissionType: "S tronic", TransmissionCode: "DQ200", FuelType: "Benzin", Drivetrain: "Önden Çekiş", Acceleration: 10.8, EngineType: "Turbo", Cylinders: 4, StartYear: 2010},
			{Name: "1.6 TDI", PowerHP: 105, TorqueNM: 250, TransmissionType: "S tronic", TransmissionCode: "DQ200", FuelType: "Dizel", Drivetrain: "Önden Çekiş", Acceleration: 11.7, EngineType: "Turbo Dizel (CR)", Cylinders: 4, StartYear: 2009},
			{Name: "2.0 TDI", PowerHP: 140, TorqueNM: 320, TransmissionType: "S tronic", TransmissionCode: "DQ250", FuelType: "Dizel", Drivetrain: "Önden Çekiş", Acceleration: 9.0, EngineType: "Turbo Dizel (CR)", Cylinders: 4, StartYear: 2008},
		})
	}

	// A3 8P (2003-2008)
	if id, ok := a3GenIDs["8P"]; ok {
		insertTrims(db, id, a3ID, []Trim{
			{Name: "1.6", PowerHP: 102, TorqueNM: 148, TransmissionType: "Tiptronic", TransmissionCode: "09G/AQ250", FuelType: "Benzin", Drivetrain: "Önden Çekiş", Acceleration: 11.9, EngineType: "Atmosferik", Cylinders: 4, StartYear: 2003},
			{Name: "1.6 FSI", PowerHP: 115, TorqueNM: 155, TransmissionType: "Tiptronic", TransmissionCode: "09G", FuelType: "Benzin", Drivetrain: "Önden Çekiş", Acceleration: 10.9, EngineType: "FSI", Cylinders: 4, StartYear: 2004, EndYear: intPtr(2007)},
			{Name: "1.9 TDI", PowerHP: 105, TorqueNM: 250, TransmissionType: "S tronic", TransmissionCode: "DQ250", FuelType: "Dizel", Drivetrain: "Önden Çekiş", Acceleration: 11.4, EngineType: "Turbo Dizel (PD)", Cylinders: 4, StartYear: 2003, EndYear: intPtr(2009)},
			{Name: "2.0 TDI", PowerHP: 140, TorqueNM: 320, TransmissionType: "S tronic", TransmissionCode: "DQ250", FuelType: "Dizel", Drivetrain: "Önden Çekiş", Acceleration: 9.0, EngineType: "Turbo Dizel (PD/CR)", Cylinders: 4, StartYear: 2004},
		})
	}

	// Verification
	fmt.Println("\n=== VERIFICATION ===")

	fmt.Println("\nAudi A4:")
	rows, _ := db.Query(`SELECT g.code, COUNT(t.id) FROM generations g LEFT JOIN trims t ON t.generation_id = g.id WHERE g.model_id = ? GROUP BY g.id ORDER BY g.start_year DESC`, a4ID)
	for rows.Next() {
		var code string
		var count int
		rows.Scan(&code, &count)
		fmt.Printf("  %s: %d trims\n", code, count)
	}
	rows.Close()

	fmt.Println("\nAudi A3:")
	rows, _ = db.Query(`SELECT g.code, COUNT(t.id) FROM generations g LEFT JOIN trims t ON t.generation_id = g.id WHERE g.model_id = ? GROUP BY g.id ORDER BY g.start_year DESC`, a3ID)
	for rows.Next() {
		var code string
		var count int
		rows.Scan(&code, &count)
		fmt.Printf("  %s: %d trims\n", code, count)
	}
	rows.Close()

	fmt.Println("\n✓ Audi data updated with correct specifications!")
}

func insertTrims(db *sql.DB, genID, modelID int64, trims []Trim) {
	for _, t := range trims {
		year := t.StartYear
		if year == 0 {
			year = 2020
		}
		_, err := db.Exec(`
			INSERT INTO trims (
				generation_id, model_id, name, power_hp, torque_nm, 
				transmission_type, transmission_code, fuel_type, drivetrain, 
				acceleration_0_100, engine_type, cylinders, year, start_year, end_year,
				seating_capacity, currency, market
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 5, 'TRY', 'TR')
		`, genID, modelID, t.Name, t.PowerHP, t.TorqueNM,
			t.TransmissionType, t.TransmissionCode, t.FuelType, t.Drivetrain,
			t.Acceleration, t.EngineType, t.Cylinders, year, t.StartYear, t.EndYear)
		if err != nil {
			fmt.Printf("    ✗ %s: %v\n", t.Name, err)
		}
	}
}

func intPtr(i int) *int {
	return &i
}
