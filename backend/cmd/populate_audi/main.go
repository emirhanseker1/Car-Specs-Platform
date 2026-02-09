package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	_ "modernc.org/sqlite"
)

// Data structures to hold the seeding info
type SeedTrim struct {
	Name           string
	StartYear      int
	EndYear        int
	IsFacelift     bool // Trim level override or redundant if Gen is split
	EngineType     string
	PowerPS        int
	TorqueNM       int
	Transmission   string
	TransCode      string
	Accel          string
	FuelType       string
	DriveTrain     string
	Cylinders      int
	DisplacementCC int
}

type SeedGeneration struct {
	Code       string
	Name       string
	StartYear  int
	EndYear    int
	IsFacelift bool
	Trims      []SeedTrim
}

type SeedModel struct {
	Name        string
	Generations []SeedGeneration
}

func main() {
	// 1. Connect usage
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		if _, err := os.Stat("vehicles.db"); err == nil {
			dbPath = "vehicles.db"
		} else if _, err := os.Stat("../../vehicles.db"); err == nil {
			dbPath = "../../vehicles.db"
		} else {
			dbPath = "vehicles.db"
		}
	}

	fmt.Printf("📂 Using Database: %s\n", dbPath)
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 2. Ensure Brand "Audi"
	audiID := getOrCreateBrand(db, "Audi", "Germany", "https://upload.wikimedia.org/wikipedia/commons/thumb/9/92/Audi-Logo_2016.svg/2560px-Audi-Logo_2016.svg.png")
	fmt.Printf("✅ Brand 'Audi' ID: %d\n", audiID)

	// 3. Define Data
	modelsData := []SeedModel{
		{
			Name: "A3",
			Generations: []SeedGeneration{
				{
					Code:       "8P",
					Name:       "8P (2003-2012)",
					StartYear:  2003,
					EndYear:    2012,
					IsFacelift: false,
					Trims: []SeedTrim{
						{Name: "1.6", StartYear: 2003, EndYear: 2012, EngineType: "4 Silindirli, Atmosferik", PowerPS: 102, TorqueNM: 148, Transmission: "6 İleri Tiptronic", TransCode: "09G / AQ250", Accel: "11.9", FuelType: "Benzin"},
						{Name: "1.6 FSI", StartYear: 2004, EndYear: 2007, EngineType: "4 Silindirli, FSI", PowerPS: 115, TorqueNM: 155, Transmission: "6 İleri Tiptronic", TransCode: "09G", Accel: "10.9", FuelType: "Benzin"},
						{Name: "1.4 TFSI", StartYear: 2007, EndYear: 2012, EngineType: "4 Silindirli, Turbo Benzinli", PowerPS: 125, TorqueNM: 200, Transmission: "7 İleri S tronic", TransCode: "DQ200", Accel: "9.3", FuelType: "Benzin", IsFacelift: true},
						{Name: "1.2 TFSI", StartYear: 2010, EndYear: 2012, EngineType: "4 Silindirli, Turbo Benzinli", PowerPS: 105, TorqueNM: 175, Transmission: "7 İleri S tronic", TransCode: "DQ200", Accel: "10.8", FuelType: "Benzin", IsFacelift: true},
						{Name: "1.9 TDI", StartYear: 2003, EndYear: 2009, EngineType: "4 Silindirli, Turbo Dizel (PD)", PowerPS: 105, TorqueNM: 250, Transmission: "6 İleri S tronic", TransCode: "DQ250", Accel: "11.4", FuelType: "Dizel"},
						{Name: "1.6 TDI", StartYear: 2009, EndYear: 2012, EngineType: "4 Silindirli, Turbo Dizel (CR)", PowerPS: 105, TorqueNM: 250, Transmission: "7 İleri S tronic", TransCode: "DQ200", Accel: "11.7", FuelType: "Dizel", IsFacelift: true},
						{Name: "2.0 TDI", StartYear: 2003, EndYear: 2012, EngineType: "4 Silindirli, Turbo Dizel", PowerPS: 140, TorqueNM: 320, Transmission: "6 İleri S tronic", TransCode: "DQ250", Accel: "9.0", FuelType: "Dizel"},
					},
				},
				// 8V Split
				{
					Code:       "8V",
					Name:       "8V (2012-2016)",
					StartYear:  2012,
					EndYear:    2016,
					IsFacelift: false,
					Trims: []SeedTrim{
						{Name: "1.2 TFSI", StartYear: 2012, EndYear: 2014, PowerPS: 105, TorqueNM: 175, Transmission: "7 İleri S tronic", TransCode: "DQ200", Accel: "10.3", FuelType: "Benzin", EngineType: "4 Silindirli, Turbo Benzinli"},
						{Name: "1.2 TFSI (110hp)", StartYear: 2014, EndYear: 2016, PowerPS: 110, TorqueNM: 175, Transmission: "7 İleri S tronic", TransCode: "DQ200", Accel: "10.3", FuelType: "Benzin", EngineType: "4 Silindirli, Turbo Benzinli"},
						{Name: "1.4 TFSI", StartYear: 2012, EndYear: 2014, PowerPS: 122, TorqueNM: 200, Transmission: "7 İleri S tronic", TransCode: "DQ200", Accel: "9.3", FuelType: "Benzin", EngineType: "4 Silindirli, Turbo Benzinli"},
						{Name: "1.4 TFSI (125hp)", StartYear: 2014, EndYear: 2016, PowerPS: 125, TorqueNM: 200, Transmission: "7 İleri S tronic", TransCode: "DQ200", Accel: "9.3", FuelType: "Benzin", EngineType: "4 Silindirli, Turbo Benzinli"},
						{Name: "1.4 TFSI COD", StartYear: 2012, EndYear: 2016, PowerPS: 140, TorqueNM: 250, Transmission: "7 İleri S tronic", TransCode: "DQ200", Accel: "8.4", FuelType: "Benzin", EngineType: "4 Silindirli, Turbo Benzinli (COD)"},
						{Name: "1.6 TDI", StartYear: 2012, EndYear: 2014, PowerPS: 105, TorqueNM: 250, Transmission: "7 İleri S tronic", TransCode: "DQ200", Accel: "10.7", FuelType: "Dizel", EngineType: "4 Silindirli, Turbo Dizel"},
						{Name: "1.6 TDI (110hp)", StartYear: 2014, EndYear: 2016, PowerPS: 110, TorqueNM: 250, Transmission: "7 İleri S tronic", TransCode: "DQ200", Accel: "10.7", FuelType: "Dizel", EngineType: "4 Silindirli, Turbo Dizel"},
					},
				},
				{
					Code:       "8V",
					Name:       "8V Makyajlı (2016-2020)",
					StartYear:  2016,
					EndYear:    2020,
					IsFacelift: true,
					Trims: []SeedTrim{
						{Name: "1.0 TFSI (30 TFSI)", StartYear: 2016, EndYear: 2020, PowerPS: 116, TorqueNM: 200, Transmission: "7 İleri S tronic", TransCode: "DQ200", Accel: "9.9", FuelType: "Benzin", EngineType: "3 Silindirli, Turbo Benzinli"},
						{Name: "1.4 TFSI COD", StartYear: 2016, EndYear: 2017, PowerPS: 150, TorqueNM: 250, Transmission: "7 İleri S tronic", TransCode: "DQ200", Accel: "8.2", FuelType: "Benzin", EngineType: "4 Silindirli, Turbo Benzinli"},
						{Name: "1.5 TFSI (35 TFSI)", StartYear: 2017, EndYear: 2020, PowerPS: 150, TorqueNM: 250, Transmission: "7 İleri S tronic", TransCode: "DQ200", Accel: "8.2", FuelType: "Benzin", EngineType: "4 Silindirli, Turbo Benzinli (COD)"},
						{Name: "1.6 TDI (30 TDI)", StartYear: 2016, EndYear: 2017, PowerPS: 110, TorqueNM: 250, Transmission: "7 İleri S tronic", TransCode: "DQ200", Accel: "10.4", FuelType: "Dizel", EngineType: "4 Silindirli, Turbo Dizel"},
						{Name: "1.6 TDI (30 TDI 116hp)", StartYear: 2017, EndYear: 2020, PowerPS: 116, TorqueNM: 250, Transmission: "7 İleri S tronic", TransCode: "DQ200", Accel: "10.4", FuelType: "Dizel", EngineType: "4 Silindirli, Turbo Dizel"},
						{Name: "S3 Sportback / Sedan", StartYear: 2016, EndYear: 2020, PowerPS: 310, TorqueNM: 400, Transmission: "7 İleri S tronic", TransCode: "DQ381", Accel: "4.6", FuelType: "Benzin", EngineType: "2.0 TFSI (4 Silindir)", DriveTrain: "quattro"},
					},
				},
			},
		},
		{
			Name: "A4",
			Generations: []SeedGeneration{
				{
					Code:       "B7",
					Name:       "B7 (2004-2008)",
					StartYear:  2004,
					EndYear:    2008,
					IsFacelift: false,
					Trims: []SeedTrim{
						{Name: "1.6", StartYear: 2004, EndYear: 2008, EngineType: "4 Silindirli, Atmosferik Benzinli", PowerPS: 102, TorqueNM: 148, Transmission: "5 İleri Manuel", TransCode: "-", Accel: "12.6", FuelType: "Benzin"},
						{Name: "2.0 TDI", StartYear: 2004, EndYear: 2008, EngineType: "4 Silindirli, Turbo Dizel", PowerPS: 140, TorqueNM: 320, Transmission: "7 İleri Multitronic", TransCode: "01J / VL300", Accel: "9.7", FuelType: "Dizel"},
						{Name: "1.8 T", StartYear: 2004, EndYear: 2008, EngineType: "4 Silindirli, Turbo Benzinli", PowerPS: 163, TorqueNM: 225, Transmission: "7 İleri Multitronic", TransCode: "01J / VL300", Accel: "8.6", FuelType: "Benzin"},
						{Name: "2.0 TFSI quattro", StartYear: 2004, EndYear: 2008, EngineType: "4 Silindirli, Turbo Benzinli", PowerPS: 200, TorqueNM: 280, Transmission: "6 İleri Tiptronic", TransCode: "ZF 6HP", Accel: "7.7", FuelType: "Benzin", DriveTrain: "quattro"},
					},
				},
				// B8 Split
				{
					Code:       "B8",
					Name:       "B8 (2008-2012)",
					StartYear:  2008,
					EndYear:    2012,
					IsFacelift: false,
					Trims: []SeedTrim{
						{Name: "1.8 TFSI", StartYear: 2008, EndYear: 2012, EngineType: "4 Silindirli, Turbo Benzinli", PowerPS: 160, TorqueNM: 250, Transmission: "8 İleri Multitronic", TransCode: "VL381", Accel: "8.3", FuelType: "Benzin"},
						{Name: "2.0 TDI", StartYear: 2008, EndYear: 2012, EngineType: "4 Silindirli, Turbo Dizel", PowerPS: 143, TorqueNM: 320, Transmission: "8 İleri Multitronic", TransCode: "VL381", Accel: "9.1", FuelType: "Dizel"},
						{Name: "2.0 TFSI quattro", StartYear: 2008, EndYear: 2012, EngineType: "4 Silindirli, Turbo Benzinli", PowerPS: 211, TorqueNM: 350, Transmission: "7 İleri S tronic", TransCode: "DL501", Accel: "6.5", FuelType: "Benzin", DriveTrain: "quattro"},
					},
				},
				{
					Code:       "B8.5",
					Name:       "B8 (2012-2016)",
					StartYear:  2012,
					EndYear:    2016,
					IsFacelift: true,
					Trims: []SeedTrim{
						{Name: "1.8 TFSI", StartYear: 2012, EndYear: 2016, EngineType: "4 Silindirli, Turbo Benzinli", PowerPS: 170, TorqueNM: 320, Transmission: "8 İleri Multitronic", TransCode: "VL381", Accel: "8.3", FuelType: "Benzin"},
						{Name: "2.0 TDI", StartYear: 2012, EndYear: 2016, EngineType: "4 Silindirli, Turbo Dizel", PowerPS: 150, TorqueNM: 320, Transmission: "8 İleri Multitronic", TransCode: "VL381", Accel: "9.1", FuelType: "Dizel"},
						{Name: "2.0 TDI (177hp)", StartYear: 2012, EndYear: 2016, EngineType: "4 Silindirli, Turbo Dizel", PowerPS: 177, TorqueNM: 380, Transmission: "8 İleri Multitronic", TransCode: "VL381", Accel: "8.5", FuelType: "Dizel"}, // Est accel
						{Name: "2.0 TFSI quattro", StartYear: 2012, EndYear: 2016, EngineType: "4 Silindirli, Turbo Benzinli", PowerPS: 225, TorqueNM: 350, Transmission: "7 İleri S tronic", TransCode: "DL501", Accel: "6.5", FuelType: "Benzin", DriveTrain: "quattro"},
					},
				},
				// B9 Split
				{
					Code:       "B9",
					Name:       "B9 (2016-2020)",
					StartYear:  2016,
					EndYear:    2020,
					IsFacelift: false,
					Trims: []SeedTrim{
						{Name: "1.4 TFSI", StartYear: 2016, EndYear: 2019, EngineType: "4 Silindirli, Turbo Benzinli", PowerPS: 150, TorqueNM: 250, Transmission: "7 İleri S tronic", TransCode: "DL382", Accel: "8.5", FuelType: "Benzin"},
						{Name: "2.0 TDI (190)", StartYear: 2016, EndYear: 2020, EngineType: "4 Silindirli, Turbo Dizel", PowerPS: 190, TorqueNM: 400, Transmission: "7 İleri S tronic", TransCode: "DL382", Accel: "7.2", FuelType: "Dizel", DriveTrain: "quattro"},
					},
				},
				{
					Code:       "B9.5",
					Name:       "B9 (2020-Present)",
					StartYear:  2020,
					EndYear:    2024,
					IsFacelift: true,
					Trims: []SeedTrim{
						{Name: "40 TDI", StartYear: 2020, EndYear: 2024, EngineType: "4 Silindirli, Turbo Dizel + Mild Hybrid", PowerPS: 204, TorqueNM: 400, Transmission: "7 İleri S tronic", TransCode: "DL382+", Accel: "6.9", FuelType: "Dizel (MHEV)", DriveTrain: "quattro"},
						{Name: "45 TFSI", StartYear: 2020, EndYear: 2024, EngineType: "4 Silindirli, Turbo Benzinli + Mild Hybrid", PowerPS: 265, TorqueNM: 370, Transmission: "7 İleri S tronic", TransCode: "DL382+", Accel: "5.5", FuelType: "Benzin (MHEV)", DriveTrain: "quattro"},
					},
				},
			},
		},
	}

	// 4. Processing Loop
	for _, mData := range modelsData {
		modelID := getOrCreateModel(db, audiID, mData.Name)
		fmt.Printf("🚙 Model '%s' ID: %d\n", mData.Name, modelID)

		for _, gData := range mData.Generations {
			genID := getOrCreateGeneration(db, modelID, gData.Code, gData.Name, gData.StartYear, gData.EndYear, gData.IsFacelift)
			fmt.Printf("  📅 Generation '%s' ID: %d\n", gData.Name, genID)

			for _, tData := range gData.Trims {
				// Inherit generation IsFacelift if not specified
				isFL := tData.IsFacelift
				if gData.IsFacelift {
					isFL = true
				}
				tData.IsFacelift = isFL
				createTrim(db, genID, modelID, tData)
				fmt.Printf("    ✨ Trim '%s' Added\n", tData.Name)
			}
		}
	}

	fmt.Println("🎉 All Done!")
}

func getOrCreateBrand(db *sql.DB, name, country, logo string) int64 {
	var id int64
	err := db.QueryRow("SELECT id FROM brands WHERE name = ?", name).Scan(&id)
	if err == sql.ErrNoRows {
		res, err := db.Exec("INSERT INTO brands (name, country, logo_url, created_at, updated_at) VALUES (?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)", name, country, logo)
		if err != nil {
			log.Fatalf("Failed to create brand: %v", err)
		}
		id, _ = res.LastInsertId()
	} else if err != nil {
		log.Fatal(err)
	}
	return id
}

func getOrCreateModel(db *sql.DB, brandID int64, name string) int64 {
	var id int64
	err := db.QueryRow("SELECT id FROM models WHERE brand_id = ? AND name = ?", brandID, name).Scan(&id)
	if err == sql.ErrNoRows {
		res, err := db.Exec("INSERT INTO models (brand_id, name, created_at, updated_at) VALUES (?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)", brandID, name)
		if err != nil {
			log.Fatalf("Failed to create model: %v", err)
		}
		id, _ = res.LastInsertId()
	} else if err != nil {
		log.Fatal(err)
	}
	return id
}

func getOrCreateGeneration(db *sql.DB, modelID int64, code, name string, start, end int, isFacelift bool) int64 {
	var id int64
	// Match by code AND is_facelift to ensure splits are stored as separate generations
	err := db.QueryRow("SELECT id FROM generations WHERE model_id = ? AND code = ? AND is_facelift = ?", modelID, code, isFacelift).Scan(&id)

	endYearVal := sql.NullInt64{}
	if end != 0 {
		endYearVal.Int64 = int64(end)
		endYearVal.Valid = true
	}

	if err == sql.ErrNoRows {
		// Use is_facelift column
		query := "INSERT INTO generations (model_id, code, name, start_year, end_year, is_facelift, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)"
		res, err := db.Exec(query, modelID, code, name, start, endYearVal, isFacelift)
		if err != nil {
			log.Fatalf("Failed to create generation: %v", err)
		}
		id, _ = res.LastInsertId()
	} else if err != nil {
		log.Fatal(err)
	}
	return id
}

func createTrim(db *sql.DB, genID, modelID int64, t SeedTrim) {
	// Check if already exists to avoid duplicates
	var exists int
	err := db.QueryRow("SELECT count(*) FROM trims WHERE generation_id = ? AND name = ? AND power_hp = ?", genID, t.Name, t.PowerPS).Scan(&exists)
	if err != nil {
		log.Fatal(err)
	}
	if exists > 0 {
		return // Skip
	}

	query := `
	INSERT INTO trims (
		generation_id, model_id, name, year, start_year, end_year, is_facelift, market,
		engine_type, fuel_type, power_hp, torque_nm, 
		transmission_type, transmission_code, acceleration_0_100, drivetrain,
		cylinders, displacement_cc,
		created_at, updated_at
	) VALUES (
		?, ?, ?, ?, ?, ?, ?, ?,
		?, ?, ?, ?, 
		?, ?, ?, ?,
		?, ?,
		CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
	)
	`

	drive := t.DriveTrain
	if drive == "" {
		drive = "Önden Çekiş"
	}

	accelVal := sql.NullFloat64{}
	if t.Accel != "" {
		// handle ~ or text
		clean := strings.ReplaceAll(t.Accel, "~", "")
		clean = strings.ReplaceAll(clean, ",", ".") // 10,3 -> 10.3
		var f float64
		fmt.Sscanf(clean, "%f", &f)
		accelVal.Float64 = f
		accelVal.Valid = true
	}

	var endYear interface{} = t.EndYear
	if t.EndYear == 0 {
		endYear = nil // Null
	}

	_, err = db.Exec(query,
		genID, modelID, t.Name, t.StartYear, t.StartYear, endYear, t.IsFacelift, "TR",
		t.EngineType, t.FuelType, t.PowerPS, t.TorqueNM,
		t.Transmission, t.TransCode, accelVal, drive,
		t.Cylinders, t.DisplacementCC,
	)
	if err != nil {
		log.Printf("❌ Failed to insert trim %s: %v", t.Name, err)
	}
}
