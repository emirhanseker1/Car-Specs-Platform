package main

import (
	"database/sql"
	"fmt"
	"log"
	"path/filepath"
	"runtime"

	_ "modernc.org/sqlite"
)

func main() {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	dbPath := filepath.Join(basepath, "vehicles.db")

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// 0. Ensure column exists
	_, err = db.Exec("ALTER TABLE trims ADD COLUMN transmission_code TEXT")
	if err != nil {
		fmt.Println("Column transmission_code likely already exists or error:", err)
	}

	var genID int64
	err = db.QueryRow("SELECT id FROM generations WHERE model_id = 1 AND code = '8Y'").Scan(&genID)
	if err != nil {
		log.Fatalf("Failed to find generation 8Y: %v", err)
	}

	// Delete existing trims for this generation to avoid duplicates
	_, err = db.Exec("DELETE FROM trims WHERE generation_id = ?", genID)
	if err != nil {
		log.Fatalf("Failed to delete existing trims: %v", err)
	}

	trims := []struct {
		Name         string
		EngineType   string
		Cylinders    int
		PowerHP      int
		TorqueNM     int
		Transmission string
		TransCode    string
		Accel        float64
		TopSpeed     int
		Drivetrain   string
		StartYear    int
		EndYear      any // int or nil
	}{
		{
			Name: "30 TFSI", EngineType: "1.0 Turbo Benzinli", Cylinders: 3, PowerHP: 110, TorqueNM: 200,
			Transmission: "7 Ileri S tronic", TransCode: "DQ200 (Kuru Kavrama)", Accel: 10.6, TopSpeed: 210, Drivetrain: "FWD",
			StartYear: 2020, EndYear: 2024,
		},
		{
			Name: "35 TFSI", EngineType: "1.5 Turbo Benzinli", Cylinders: 4, PowerHP: 150, TorqueNM: 250,
			Transmission: "7 Ileri S tronic", TransCode: "DQ200 (Kuru Kavrama)", Accel: 8.4, TopSpeed: 232, Drivetrain: "FWD",
			StartYear: 2020, EndYear: nil,
		},
		{
			Name: "Yeni 30 TFSI (Makyajlı)", EngineType: "1.5 Turbo Benzinli", Cylinders: 4, PowerHP: 116, TorqueNM: 220,
			Transmission: "7 Ileri S tronic", TransCode: "DQ200 (Kuru Kavrama)", Accel: 9.9, TopSpeed: 210, Drivetrain: "FWD",
			StartYear: 2024, EndYear: nil,
		},
		{
			Name: "S3 Sedan / Sportback", EngineType: "2.0 Turbo Benzinli", Cylinders: 4, PowerHP: 310, TorqueNM: 400,
			Transmission: "7 Ileri S tronic", TransCode: "DQ381 (Islak Kavrama)", Accel: 4.8, TopSpeed: 250, Drivetrain: "quattro",
			StartYear: 2020, EndYear: nil,
		},
		{
			Name: "RS3 Sportback / Sedan", EngineType: "2.5 Turbo Benzinli", Cylinders: 5, PowerHP: 400, TorqueNM: 500,
			Transmission: "7 Ileri S tronic", TransCode: "DQ500 (Islak Kavrama)", Accel: 3.8, TopSpeed: 250, Drivetrain: "quattro",
			StartYear: 2020, EndYear: nil,
		},
	}

	for _, t := range trims {
		res, err := db.Exec(`
			INSERT INTO trims (
				generation_id, model_id, name, year, start_year, end_year,
				engine_type, cylinders, power_hp, torque_nm, 
				transmission_type, transmission_code, acceleration_0_100, top_speed_kmh, drivetrain, market, created_at, updated_at
			) VALUES (?, 1, ?, 2020, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 'TR', DATETIME('now'), DATETIME('now'))
		`, genID, t.Name, t.StartYear, t.EndYear, t.EngineType, t.Cylinders, t.PowerHP, t.TorqueNM, t.Transmission, t.TransCode, t.Accel, t.TopSpeed, t.Drivetrain)

		if err != nil {
			log.Printf("Failed to insert trim %s: %v", t.Name, err)
			continue
		}

		trimID, _ := res.LastInsertId()
		fmt.Printf("Inserted: %s\n", t.Name)

		// Insert Specs (including redundantly in case we want details)
		specs := map[string]string{
			"Şanzıman Kodu": t.TransCode,
			"Motor Tipi":    t.EngineType,
		}

		for k, v := range specs {
			db.Exec("INSERT INTO specs (trim_id, category, name, value) VALUES (?, 'Technical', ?, ?)", trimID, k, v)
		}
	}
}
