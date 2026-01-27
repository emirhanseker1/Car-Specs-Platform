package main

import (
	"database/sql"
	stdjson "encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/emirh/car-specs/backend/internal/config"
	"github.com/emirh/car-specs/backend/internal/db"
	"github.com/emirh/car-specs/backend/internal/importer"
	"github.com/emirh/car-specs/backend/internal/json"
	"github.com/emirh/car-specs/backend/internal/models"
	"github.com/emirh/car-specs/backend/internal/scraper"
	"github.com/emirh/car-specs/backend/pkg/apininjas"
	"github.com/emirh/car-specs/backend/pkg/carquery"
	"github.com/emirh/car-specs/backend/pkg/utils"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func main() {
	// --- CarQuery API Demo ---
	// This runs on startup to demonstrate the requested functionality
	/*
		fmt.Println("--- CarQuery API Demo ---")
		client := carquery.NewClient()

		// 1. Get Years
		minYear, maxYear, err := client.GetYears()
		if err != nil {
			log.Printf("Error getting years: %v", err)
		} else {
			fmt.Printf("Available Years: %d - %d\n", minYear, maxYear)
		}

		// 2. Get Makes for 2020
		makes, err := client.GetMakes(2020)
		if err != nil {
			log.Printf("Error getting makes: %v", err)
		} else {
			fmt.Printf("Found %d makes for 2020. First 3: ", len(makes))
			for i, m := range makes {
				if i >= 3 {
					break
				}
				fmt.Printf("%s, ", m.MakeDisplay)
			}
			fmt.Println("...")
		}

		// 3. Get Models for BMW in 2020
		carModels, err := client.GetModels("bmw", 2020)
		if err != nil {
			log.Printf("Error getting BMW models: %v", err)
		} else {
			fmt.Printf("Found %d BMW models for 2020. Examples: ", len(carModels))
			for i, m := range carModels {
				if i >= 3 {
					break
				}
				fmt.Printf("%s, ", m.ModelName)
			}
			fmt.Println("...")
		}

		// 4. Get Trims for a specific model (e.g., '3 Series' in 2020)
		trims, err := client.GetTrims(carquery.TrimFilter{
			Make:  "bmw",
			Model: "3 Series",
			Year:  2020,
		})
		if err != nil {
			log.Printf("Error getting trims: %v", err)
		} else {
			fmt.Printf("Found %d trims for BMW 3 Series (2020).\n", len(trims))
			if len(trims) > 0 {
				t := trims[0]
				fmt.Printf("Sample Trim: %s %s (%s) - %s HP\n", t.ModelMakeDisplay, t.ModelName, t.ModelTrim, t.ModelPowerPS)
			}
		}
		fmt.Println("--- End Demo ---")
	*/
	// -------------------------------

	// Load Configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	scrape := flag.Bool("scrape", false, "Run the scraper on startup")
	importCSV := flag.String("import", "", "Import data from a CSV file into the database and exit")
	syncCarQuery := flag.Bool("sync-carquery", false, "Sync data from CarQuery API (v2)")
	syncNinjas := flag.Bool("sync-ninjas", false, "Sync data from API Ninjas Cars API")
	demoDetailed := flag.Bool("demo-detailed", false, "Fetch detailed 2024-2025 data for demo")
	seedDemo := flag.Bool("seed", false, "Seed demo data for testing")
	flag.Parse()

	// --- CLI MODES ---

	if *seedDemo {
		fmt.Println("--- Seeding Mode ---")
		gormDB, err := gorm.Open(sqlite.Open(cfg.DatabaseURL), &gorm.Config{})
		if err != nil {
			log.Fatalf("Failed to connect to GORM DB: %v", err)
		}

		// Manual migration fix
		gormDB.Exec("ALTER TABLE trims ADD COLUMN model_id INTEGER")
		gormDB.Exec("ALTER TABLE trims ADD COLUMN vehicle_id INTEGER")
		gormDB.Exec("ALTER TABLE trims ADD COLUMN year INTEGER")

		if err := gormDB.AutoMigrate(&models.Make{}, &models.Model{}, &models.Trim{}, &models.Spec{}); err != nil {
			log.Printf("AutoMigrate partial warning: %v", err)
		}

		if err := importer.SeedDemoData(gormDB); err != nil {
			log.Fatalf("Seed failed: %v", err)
		}
		fmt.Println("Seeding completed.")
		return
	}

	// Mode 3: Sync API Ninjas (New Source)
	if *syncNinjas {
		fmt.Println("--- Starting API Ninjas Sync ---")
		// Check removed, handled by LoadConfig

		// Init GORM
		targetDB := "vehicles.db"
		if envDB := os.Getenv("DATABASE_URL"); envDB != "" {
			targetDB = envDB
		}
		gormDB, err := gorm.Open(sqlite.Open(targetDB), &gorm.Config{})
		if err != nil {
			log.Fatalf("Failed to connect to GORM DB: %v", err)
		}

		// Schema Fixes (Manual migration helpers)
		// Clean start: Drop tables removed to ensure persistence
		// fmt.Println("Cleaning old data...")
		// gormDB.Migrator().DropTable(&models.Spec{}, &models.Trim{}, &models.Model{}, &models.Make{})

		// Migrate
		if err := gormDB.AutoMigrate(&models.Make{}, &models.Model{}, &models.Trim{}, &models.Spec{}, &models.Chronicle{}); err != nil {
			log.Printf("AutoMigrate partial warning: %v", err)
		}

		// Init Client
		client := apininjas.NewClient()

		// List of makes to sync (Demo)
		makesToSync := []string{"BMW", "Audi", "Toyota", "Honda", "Mercedes-Benz", "Ford"}

		if err := importer.SyncApiNinjasData(gormDB, client, makesToSync); err != nil {
			log.Fatalf("Ninjas Sync failed: %v", err)
		}
		fmt.Println("Ninjas Sync Completed Successfully.")
		return
	}

	// Mode 4: Demo Detailed Sync (2024-2025)
	if *demoDetailed {
		fmt.Println("--- Starting Demo Detailed Sync (2024-2025) ---")

		targetDB := "vehicles.db"
		if envDB := os.Getenv("DATABASE_URL"); envDB != "" {
			targetDB = envDB
		}
		gormDB, err := gorm.Open(sqlite.Open(targetDB), &gorm.Config{})
		if err != nil {
			log.Fatalf("Failed to connect to GORM DB: %v", err)
		}

		// Clean start for demo
		// fmt.Println("Cleaning old data for demo...")
		// gormDB.Migrator().DropTable(&models.Spec{}, &models.Trim{}, &models.Model{}, &models.Make{})
		if err := gormDB.AutoMigrate(&models.Make{}, &models.Model{}, &models.Trim{}, &models.Spec{}, &models.Chronicle{}); err != nil {
			log.Printf("AutoMigrate warning: %v", err)
		}

		client := apininjas.NewClient()
		makes := []string{"BMW", "Audi", "Mercedes-Benz", "Toyota", "Ford", "Honda", "Porsche", "Ferrari"}
		years := []int{2023, 2024, 2025}

		if err := importer.SyncDetailedYears(gormDB, client, makes, years); err != nil {
			log.Fatalf("Demo Sync failed: %v", err)
		}
		fmt.Println("Demo Data (2024-2025) Ready! Please run server normally.")
		return
	}

	// Existing main logic continues...

	// Mode 1: Sync CarQuery Data (Uses GORM) - Preserved if user wants to switch back
	if *syncCarQuery {
		fmt.Println("--- Starting CarQuery Sync Mode ---")
		// Open GORM connection
		gormDB, err := gorm.Open(sqlite.Open(cfg.DatabaseURL), &gorm.Config{})
		if err != nil {
			log.Fatalf("Failed to connect to GORM DB: %v", err)
		}

		// Auto Migrate the new models
		fmt.Println("Migrating GORM models...")
		// Attempt to manually add model_id if missing to help GORM (workaround for SQLite migration issue)
		// We use raw SQL on the underlying connection or via GORM
		gormDB.Exec("ALTER TABLE trims ADD COLUMN model_id INTEGER")
		gormDB.Exec("ALTER TABLE trims ADD COLUMN vehicle_id INTEGER")
		gormDB.Exec("ALTER TABLE trims ADD COLUMN year INTEGER")

		if err := gormDB.AutoMigrate(&models.Make{}, &models.Model{}, &models.Trim{}, &models.Spec{}); err != nil {
			log.Printf("AutoMigrate partial warning: %v", err)
			// Continue anyway as we might have fixed it manually
		}

		// Run Sync
		client := carquery.NewClient()
		startYear := 2024
		endYear := 2025
		fmt.Printf("Syncing years %d-%d...\n", startYear, endYear)

		if err := importer.SyncCarQueryData(gormDB, client, startYear, endYear); err != nil {
			log.Printf("Sync failed: %v", err)
		} else {
			fmt.Println("Sync completed successfully.")
		}
		return
	}

	// Mode 2: Server Mode (Existing Logic)
	database, err := db.InitDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	// Init GORM for mixed usage
	gormApp, err := gorm.Open(sqlite.Open("vehicles.db"), &gorm.Config{})
	if err != nil {
		log.Println("GORM init warning:", err)
	} else {
		gormApp.AutoMigrate(&models.Chronicle{}, &models.Trim{}, &models.Spec{})
	}

	if *importCSV != "" {
		if err := importer.ImportCSV(database, *importCSV); err != nil {
			log.Fatalf("Import failed: %v", err)
		}
		log.Printf("Import completed: %s", *importCSV)
		return
	}

	if *scrape {
		go func() {
			if err := scraper.ScrapeAll(database); err != nil {
				log.Printf("Scraper error: %v", err)
			}
		}()
	}

	mux := http.NewServeMux()

	// GET /health
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		json.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"}, nil)
	})

	// POST /api/admin/normalize-vehicles - Fix capitalization of all vehicle names
	mux.HandleFunc("POST /api/admin/normalize-vehicles", func(w http.ResponseWriter, r *http.Request) {
		// Use raw SQL to avoid GORM model parsing issues with Vehicle struct
		rows, err := database.Query("SELECT id, brand, model FROM vehicles")
		if err != nil {
			json.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()}, nil)
			return
		}
		defer rows.Close()

		type vehicleRow struct {
			ID    int64
			Brand string
			Model string
		}

		var vehicles []vehicleRow
		for rows.Next() {
			var v vehicleRow
			if err := rows.Scan(&v.ID, &v.Brand, &v.Model); err != nil {
				continue
			}
			vehicles = append(vehicles, v)
		}

		updated := 0
		for _, v := range vehicles {
			newBrand, newModel := utils.NormalizeVehicleName(v.Brand, v.Model)
			if newBrand != v.Brand || newModel != v.Model {
				_, err := database.Exec("UPDATE vehicles SET brand = ?, model = ? WHERE id = ?", newBrand, newModel, v.ID)
				if err != nil {
					log.Printf("Error updating vehicle %d: %v", v.ID, err)
					continue
				}
				updated++
				log.Printf("Updated: %s %s -> %s %s", v.Brand, v.Model, newBrand, newModel)
			}
		}

		json.WriteJSON(w, http.StatusOK, map[string]interface{}{
			"message": "Vehicle names normalized",
			"total":   len(vehicles),
			"updated": updated,
		}, nil)
	})

	// POST /api/admin/fix-duplicate-brands - Remove duplicate brand names from model field
	mux.HandleFunc("POST /api/admin/fix-duplicate-brands", func(w http.ResponseWriter, r *http.Request) {
		rows, err := database.Query("SELECT id, brand, model FROM vehicles")
		if err != nil {
			json.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()}, nil)
			return
		}
		defer rows.Close()

		type vehicleRow struct {
			ID    int64
			Brand string
			Model string
		}

		var vehicles []vehicleRow
		for rows.Next() {
			var v vehicleRow
			if err := rows.Scan(&v.ID, &v.Brand, &v.Model); err != nil {
				continue
			}
			vehicles = append(vehicles, v)
		}

		updated := 0
		for _, v := range vehicles {
			// Check if model starts with brand name (case-insensitive)
			brandLower := strings.ToLower(v.Brand)
			modelLower := strings.ToLower(v.Model)

			if strings.HasPrefix(modelLower, brandLower+" ") {
				// Remove brand prefix from model
				newModel := strings.TrimSpace(v.Model[len(v.Brand):])

				_, err := database.Exec("UPDATE vehicles SET model = ? WHERE id = ?", newModel, v.ID)
				if err != nil {
					log.Printf("Error updating vehicle %d: %v", v.ID, err)
					continue
				}
				updated++
				log.Printf("Fixed duplicate: %s %s -> %s %s", v.Brand, v.Model, v.Brand, newModel)
			}
		}

		json.WriteJSON(w, http.StatusOK, map[string]interface{}{
			"message": "Duplicate brand names fixed",
			"total":   len(vehicles),
			"updated": updated,
		}, nil)
	})

	// POST /api/admin/update-vehicle-image - Update single vehicle image
	mux.HandleFunc("POST /api/admin/update-vehicle-image", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			VehicleID int64  `json:"vehicle_id"`
			ImageURL  string `json:"image_url"`
		}
		if err := stdjson.NewDecoder(r.Body).Decode(&req); err != nil {
			json.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"}, nil)
			return
		}

		if _, err := database.Exec("UPDATE vehicles SET image_url = ? WHERE id = ?", req.ImageURL, req.VehicleID); err != nil {
			log.Printf("Error updating image for vehicle %d: %v", req.VehicleID, err)
			json.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()}, nil)
			return
		}

		json.WriteJSON(w, http.StatusOK, map[string]string{"status": "updated"}, nil)
		// ... (existing update-vehicle-image endpoint)
	})

	// POST /api/admin/fix-data
	mux.HandleFunc("POST /api/admin/fix-data", func(w http.ResponseWriter, r *http.Request) {
		// 1. Delete Fiats
		res, err := database.Exec("DELETE FROM vehicles WHERE brand = 'Fiat'")
		if err != nil {
			http.Error(w, "Error deleting Fiats: "+err.Error(), 500)
			return
		}
		deleted, _ := res.RowsAffected()

		// Cleanup orphans
		database.Exec("DELETE FROM trims WHERE vehicle_id NOT IN (SELECT id FROM vehicles)")
		database.Exec("DELETE FROM trim_powertrain_meta WHERE trim_id NOT IN (SELECT id FROM trims)")
		database.Exec("DELETE FROM vehicle_generation_meta WHERE vehicle_id NOT IN (SELECT id FROM vehicles)")

		// 2. Fix BMW Images
		rows, err := database.Query("SELECT id, model FROM vehicles WHERE brand = 'BMW'")
		if err != nil {
			http.Error(w, "Error querying BMWs: "+err.Error(), 500)
			return
		}
		defer rows.Close()

		updated := 0
		type Update struct {
			ID  int64
			Img string
		}
		var updates []Update

		for rows.Next() {
			var id int64
			var model string
			rows.Scan(&id, &model)
			modelLower := strings.ToLower(model)

			var img string
			if strings.Contains(modelLower, "850") || strings.Contains(modelLower, "8-series") {
				if strings.Contains(modelLower, "cabrio") {
					img = "/images/vehicles/bmw_8-series_cabrio.png"
				} else if strings.Contains(modelLower, "gran") {
					img = "/images/vehicles/bmw_8-series_gran_coupe.png"
				} else {
					img = "/images/vehicles/bmw_m850i_xdrive_coupe.png"
				}
			} else if strings.Contains(modelLower, "z4") {
				img = "/images/vehicles/bmw_z4_sdrive30i.png"
			} else if strings.Contains(modelLower, "4-series") || strings.Contains(modelLower, "430") {
				img = "/images/vehicles/bmw_430i_coupe.png"
			} else if strings.Contains(modelLower, "2-series") {
				img = "/images/vehicles/bmw_2-series_gran_coupe.png"
			} else if strings.Contains(modelLower, "ix3") {
				img = "/images/vehicles/bmw_ix3.png"
			} else if strings.Contains(modelLower, "x7") {
				img = "/images/vehicles/bmw_x7.png"
			} else {
				img = "/images/vehicles/bmw_m850i_xdrive_coupe.png"
			}
			updates = append(updates, Update{id, img})
		}
		rows.Close() // Close before update loop

		for _, u := range updates {
			_, err := database.Exec("UPDATE vehicles SET image_url = ? WHERE id = ?", u.Img, u.ID)
			if err == nil {
				updated++
			} else {
				log.Printf("Error updating vehicle %d: %v", u.ID, err)
			}
		}

		json.WriteJSON(w, http.StatusOK, map[string]interface{}{
			"status":        "success",
			"deleted_fiats": deleted,
			"updated_bmws":  updated,
		}, nil)
	})

	// GET /api/search?q=golf&fuel=Diesel&transmission=DSG&hpMin=100&hpMax=200&yearMin=2015&yearMax=2021
	mux.HandleFunc("GET /api/search", func(w http.ResponseWriter, r *http.Request) {
		type PowertrainSummary struct {
			FuelTypes         []string `json:"fuel_types,omitempty"`
			Transmissions     []string `json:"transmissions,omitempty"`
			MinHP             *int64   `json:"min_hp,omitempty"`
			MaxHP             *int64   `json:"max_hp,omitempty"`
			MinDisplacementCC *int64   `json:"min_displacement_cc,omitempty"`
			MaxDisplacementCC *int64   `json:"max_displacement_cc,omitempty"`
		}

		type SearchVehicle struct {
			models.Vehicle
			PowertrainSummary *PowertrainSummary `json:"powertrain_summary,omitempty"`
		}

		type Facets struct {
			FuelTypes     []string `json:"fuel_types,omitempty"`
			Transmissions []string `json:"transmissions,omitempty"`
			MinHP         *int64   `json:"min_hp,omitempty"`
			MaxHP         *int64   `json:"max_hp,omitempty"`
			MinYear       *int64   `json:"min_year,omitempty"`
			MaxYear       *int64   `json:"max_year,omitempty"`
		}

		type SearchResponse struct {
			Results []SearchVehicle `json:"results"`
			Facets  Facets          `json:"facets"`
		}

		q := strings.TrimSpace(r.URL.Query().Get("q"))
		brand := strings.TrimSpace(r.URL.Query().Get("brand"))
		fuel := strings.TrimSpace(r.URL.Query().Get("fuel"))
		transmission := strings.TrimSpace(r.URL.Query().Get("transmission"))

		parseInt64 := func(raw string) (*int64, bool) {
			raw = strings.TrimSpace(raw)
			if raw == "" {
				return nil, false
			}
			n, err := strconv.ParseInt(raw, 10, 64)
			if err != nil {
				return nil, false
			}
			return &n, true
		}

		hpMin, hasHPMin := parseInt64(r.URL.Query().Get("hpMin"))
		hpMax, hasHPMax := parseInt64(r.URL.Query().Get("hpMax"))
		yearMin, hasYearMin := parseInt64(r.URL.Query().Get("yearMin"))
		yearMax, hasYearMax := parseInt64(r.URL.Query().Get("yearMax"))

		args := make([]interface{}, 0)
		where := make([]string, 0)

		if brand != "" {
			where = append(where, "LOWER(v.brand) = LOWER(?)")
			args = append(args, brand)
		}

		if q != "" {
			likeQ := "%" + strings.ToLower(q) + "%"
			where = append(where, "(LOWER(v.brand) LIKE ? OR LOWER(v.model) LIKE ? OR LOWER(v.generation) LIKE ? OR EXISTS (SELECT 1 FROM trims t2 WHERE t2.vehicle_id = v.id AND LOWER(t2.name) LIKE ?))")
			args = append(args, likeQ, likeQ, likeQ, likeQ)
		}

		if fuel != "" {
			where = append(where, "EXISTS (SELECT 1 FROM trims tf JOIN trim_powertrain_meta pf ON pf.trim_id = tf.id WHERE tf.vehicle_id = v.id AND LOWER(pf.fuel_type) = LOWER(?))")
			args = append(args, fuel)
		}

		if transmission != "" {
			where = append(where, "EXISTS (SELECT 1 FROM trims tt JOIN trim_powertrain_meta pt ON pt.trim_id = tt.id WHERE tt.vehicle_id = v.id AND LOWER(pt.transmission_type) = LOWER(?))")
			args = append(args, transmission)
		}

		if hasHPMin {
			where = append(where, "EXISTS (SELECT 1 FROM trims th JOIN trim_powertrain_meta ph ON ph.trim_id = th.id WHERE th.vehicle_id = v.id AND ph.power_hp >= ?)")
			args = append(args, *hpMin)
		}
		if hasHPMax {
			where = append(where, "EXISTS (SELECT 1 FROM trims th JOIN trim_powertrain_meta ph ON ph.trim_id = th.id WHERE th.vehicle_id = v.id AND ph.power_hp <= ?)")
			args = append(args, *hpMax)
		}

		if hasYearMin {
			where = append(where, "(COALESCE(m.start_year, CASE WHEN v.generation GLOB '[0-9]*' THEN CAST(v.generation AS INTEGER) END) >= ?)")
			args = append(args, *yearMin)
		}
		if hasYearMax {
			where = append(where, "(COALESCE(m.end_year, m.start_year, CASE WHEN v.generation GLOB '[0-9]*' THEN CAST(v.generation AS INTEGER) END) <= ?)")
			args = append(args, *yearMax)
		}

		base := "SELECT v.id, v.brand, v.model, v.generation, v.image_url, v.link, m.start_year, m.end_year, m.is_facelift, m.market, " +
			"GROUP_CONCAT(DISTINCT pt.fuel_type), GROUP_CONCAT(DISTINCT pt.transmission_type), " +
			"MIN(pt.power_hp), MAX(pt.power_hp), MIN(pt.displacement_cc), MAX(pt.displacement_cc) " +
			"FROM vehicles v " +
			"LEFT JOIN vehicle_generation_meta m ON m.vehicle_id = v.id " +
			"LEFT JOIN trims t ON t.vehicle_id = v.id " +
			"LEFT JOIN trim_powertrain_meta pt ON pt.trim_id = t.id"

		if len(where) > 0 {
			base += " WHERE " + strings.Join(where, " AND ")
		}
		base += " GROUP BY v.id"

		rows, err := database.Query(base, args...)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var results []SearchVehicle
		var vehicleIDs []int64

		seenFuel := make(map[string]bool)
		seenTrans := make(map[string]bool)
		var facetsMinHP *int64
		var facetsMaxHP *int64
		var facetsMinYear *int64
		var facetsMaxYear *int64

		for rows.Next() {
			var v models.Vehicle
			var imgURL sql.NullString
			var startYear sql.NullInt64
			var endYear sql.NullInt64
			var isFacelift sql.NullInt64
			var market sql.NullString
			var fuelCSV sql.NullString
			var transCSV sql.NullString
			var minHP sql.NullInt64
			var maxHP sql.NullInt64
			var minCC sql.NullInt64
			var maxCC sql.NullInt64

			if err := rows.Scan(&v.ID, &v.Brand, &v.Model, &v.Generation, &imgURL, &v.Link, &startYear, &endYear, &isFacelift, &market, &fuelCSV, &transCSV, &minHP, &maxHP, &minCC, &maxCC); err != nil {
				continue
			}
			v.ImageURL = imgURL.String

			if startYear.Valid || endYear.Valid || isFacelift.Valid || market.Valid {
				meta := &models.VehicleGenerationMeta{VehicleID: v.ID}
				if startYear.Valid {
					meta.StartYear = &startYear.Int64
				}
				if endYear.Valid {
					meta.EndYear = &endYear.Int64
				}
				if isFacelift.Valid {
					b := isFacelift.Int64 != 0
					meta.IsFacelift = &b
				}
				if market.Valid {
					meta.Market = market.String
				}
				v.GenMeta = meta
			}

			splitCSV := func(s string) []string {
				parts := strings.Split(s, ",")
				out := make([]string, 0, len(parts))
				seen := make(map[string]bool)
				for _, p := range parts {
					pp := strings.TrimSpace(p)
					if pp == "" {
						continue
					}
					if seen[pp] {
						continue
					}
					seen[pp] = true
					out = append(out, pp)
				}
				sort.Strings(out)
				return out
			}

			summary := &PowertrainSummary{}
			if fuelCSV.Valid {
				summary.FuelTypes = splitCSV(fuelCSV.String)
				for _, f := range summary.FuelTypes {
					seenFuel[f] = true
				}
			}
			if transCSV.Valid {
				summary.Transmissions = splitCSV(transCSV.String)
				for _, t := range summary.Transmissions {
					seenTrans[t] = true
				}
			}
			if minHP.Valid {
				summary.MinHP = &minHP.Int64
				if facetsMinHP == nil || minHP.Int64 < *facetsMinHP {
					x := minHP.Int64
					facetsMinHP = &x
				}
			}
			if maxHP.Valid {
				summary.MaxHP = &maxHP.Int64
				if facetsMaxHP == nil || maxHP.Int64 > *facetsMaxHP {
					x := maxHP.Int64
					facetsMaxHP = &x
				}
			}
			if minCC.Valid {
				summary.MinDisplacementCC = &minCC.Int64
			}
			if maxCC.Valid {
				summary.MaxDisplacementCC = &maxCC.Int64
			}

			var derivedStart *int64
			if startYear.Valid {
				derivedStart = &startYear.Int64
			}
			if derivedStart != nil {
				if facetsMinYear == nil || *derivedStart < *facetsMinYear {
					x := *derivedStart
					facetsMinYear = &x
				}
				if facetsMaxYear == nil || *derivedStart > *facetsMaxYear {
					x := *derivedStart
					facetsMaxYear = &x
				}
			}
			if endYear.Valid {
				derivedEnd := endYear.Int64
				if facetsMinYear == nil || derivedEnd < *facetsMinYear {
					x := derivedEnd
					facetsMinYear = &x
				}
				if facetsMaxYear == nil || derivedEnd > *facetsMaxYear {
					x := derivedEnd
					facetsMaxYear = &x
				}
			}

			results = append(results, SearchVehicle{Vehicle: v, PowertrainSummary: summary})
			vehicleIDs = append(vehicleIDs, v.ID)
		}

		// Populate engine options for result vehicles
		if len(vehicleIDs) > 0 {
			placeholders := make([]string, 0, len(vehicleIDs))
			args2 := make([]interface{}, 0, len(vehicleIDs))
			for _, vid := range vehicleIDs {
				placeholders = append(placeholders, "?")
				args2 = append(args2, vid)
			}
			q2 := "SELECT vehicle_id, name FROM trims WHERE vehicle_id IN (" + strings.Join(placeholders, ",") + ")"
			tr, err := database.Query(q2, args2...)
			if err == nil {
				defer tr.Close()
				byVehicle := make(map[int64][]string)
				for tr.Next() {
					var vid int64
					var name string
					if err := tr.Scan(&vid, &name); err != nil {
						continue
					}
					byVehicle[vid] = append(byVehicle[vid], name)
				}
				for i := range results {
					trims := byVehicle[results[i].ID]
					results[i].EngineOptions = extractEngineOptions(trims)
				}
			}
		}

		facetFuels := make([]string, 0, len(seenFuel))
		for f := range seenFuel {
			facetFuels = append(facetFuels, f)
		}
		sort.Strings(facetFuels)

		facetTrans := make([]string, 0, len(seenTrans))
		for t := range seenTrans {
			facetTrans = append(facetTrans, t)
		}
		sort.Strings(facetTrans)

		resp := SearchResponse{
			Results: results,
			Facets: Facets{
				FuelTypes:     facetFuels,
				Transmissions: facetTrans,
				MinHP:         facetsMinHP,
				MaxHP:         facetsMaxHP,
				MinYear:       facetsMinYear,
				MaxYear:       facetsMaxYear,
			},
		}
		json.WriteJSON(w, http.StatusOK, resp, nil)
	})

	// GET /api/vehicles?brand=Fiat
	mux.HandleFunc("GET /api/vehicles", func(w http.ResponseWriter, r *http.Request) {
		brand := r.URL.Query().Get("brand")

		// Query against new schema (Models + Makes)
		// Subquery to get Class from the first trim found (demo logic)
		query := `SELECT m.id, mk.name, m.name, 
        (SELECT s.value FROM specs s JOIN trims t ON s.trim_id = t.id WHERE t.model_id = m.id AND s.name = 'Class' LIMIT 1) 
        FROM models m JOIN makes mk ON m.make_id = mk.id`
		var rows *sql.Rows
		var err error

		if brand != "" {
			query += " WHERE mk.name = ?"
			rows, err = database.Query(query, brand)
		} else {
			rows, err = database.Query(query)
		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var vehicles []models.Vehicle
		var modelIDs []uint
		for rows.Next() {
			var v models.Vehicle
			var mid uint
			// Map Model to Vehicle struct for frontend compatibility
			var class sql.NullString
			if err := rows.Scan(&mid, &v.Brand, &v.Model, &class); err != nil {
				log.Printf("Scan error: %v", err)
				continue
			}
			v.ID = int64(mid)
			// Defaults for missing fields
			v.Generation = ""
			v.ImageURL = ""
			if class.Valid {
				v.Class = class.String
			}

			vehicles = append(vehicles, v)
			modelIDs = append(modelIDs, mid)
		}

		// Populate engine options (derived from trim names)
		if len(modelIDs) > 0 {
			placeholders := make([]string, 0, len(modelIDs))
			args := make([]interface{}, 0, len(modelIDs))
			for _, mid := range modelIDs {
				placeholders = append(placeholders, "?")
				args = append(args, mid)
			}

			// Use model_id instead of vehicle_id
			q := "SELECT model_id, name FROM trims WHERE model_id IN (" + strings.Join(placeholders, ",") + ")"
			tr, err := database.Query(q, args...)
			if err == nil {
				defer tr.Close()
				byVehicle := make(map[int64][]string)
				for tr.Next() {
					var mid int64
					var name string
					if err := tr.Scan(&mid, &name); err != nil {
						continue
					}
					byVehicle[mid] = append(byVehicle[mid], name)
				}

				for i := range vehicles {
					trims := byVehicle[vehicles[i].ID]
					vehicles[i].EngineOptions = extractEngineOptions(trims)
				}
			}
		}

		// Sort by Brand then Model
		sort.Slice(vehicles, func(i, j int) bool {
			if vehicles[i].Brand == vehicles[j].Brand {
				return vehicles[i].Model < vehicles[j].Model
			}
			return vehicles[i].Brand < vehicles[j].Brand
		})

		json.WriteJSON(w, http.StatusOK, vehicles, nil)
	})

	// GET /api/vehicles/{id} (Get Vehicle Details + Trims)
	mux.HandleFunc("GET /api/vehicles/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

		// Query against new schema (Models + Makes)
		query := "SELECT m.id, mk.name, m.name FROM models m JOIN makes mk ON m.make_id = mk.id WHERE m.id = ?"

		var v models.Vehicle
		var mid uint
		err := database.QueryRow(query, id).Scan(&mid, &v.Brand, &v.Model)
		if err != nil {
			http.Error(w, "Vehicle not found", http.StatusNotFound)
			return
		}
		v.ID = int64(mid)
		v.Generation = "" // Default
		v.ImageURL = ""

		// Fetch Trims for this Model
		rows, err := database.Query("SELECT id, name, year FROM trims WHERE model_id = ?", id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var trims []models.Trim
		// In this new schema, we return Trim struct which matches what frontend expects mostly?
		// Frontend expects Trim with specs?
		// Let's stick to returning models.Trim with populated Specs if possible.
		// Wait, legacy handled 'TrimPowertrainMeta'. New schema has 'Specs'.
		// Frontend might rely on 'powertrain_meta'.
		// I should map 'Specs' to 'PowertrainMeta' if possible, or just send Specs.
		// The previous handler enriched Trim with 'PowertrainMeta'.
		// I will fetch Specs and populate 'Specs' field.
		// If frontend needs 'PowertrainMeta', I would need to Map 'Specs' -> 'PowertrainMeta' manually.
		// Let's check 'models.Trim' struct. It has 'Specs []Spec'.
		// It also has 'PowertrainMeta *TrimPowertrainMeta'.

		for rows.Next() {
			var t models.Trim
			var year int
			if err := rows.Scan(&t.ID, &t.Name, &year); err != nil {
				continue
			}
			t.Year = year // Assuming Trim struct has Year field (Step 727 confirmed)

			// Fetch specs for this trim
			specRows, err := database.Query("SELECT category, name, value FROM specs WHERE trim_id = ?", t.ID)
			if err == nil {
				var specs []models.Spec
				for specRows.Next() {
					var s models.Spec
					if err := specRows.Scan(&s.Category, &s.Name, &s.Value); err != nil {
						continue
					}
					specs = append(specs, s)
				}
				specRows.Close()
				t.Specs = specs

				// OPTIONAL: Map Specs to PowertrainMeta for backward compatibility if needed?
				// Frontend 'Home.tsx' doesn't seem to use deep details. 'CarCard' uses specs.
				// Let's assume sending 'specs' list is enough or cleaner.
			}
			trims = append(trims, t)
		}

		// Return combined response
		response := map[string]interface{}{
			"vehicle": v,
			"trims":   trims,
		}
		json.WriteJSON(w, http.StatusOK, response, nil)
	})

	// GET /api/chronicles?model_id=123
	mux.HandleFunc("GET /api/chronicles", func(w http.ResponseWriter, r *http.Request) {
		modelID := r.URL.Query().Get("model_id")
		if modelID == "" {
			http.Error(w, "model_id required", http.StatusBadRequest)
			return
		}

		var chronicles []models.Chronicle
		if err := gormApp.Where("model_id = ?", modelID).Order("upvotes desc, created_at desc").Find(&chronicles).Error; err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.WriteJSON(w, http.StatusOK, chronicles, nil)
	})

	// POST /api/chronicles
	mux.HandleFunc("POST /api/chronicles", func(w http.ResponseWriter, r *http.Request) {
		var c models.Chronicle
		if err := stdjson.NewDecoder(r.Body).Decode(&c); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Simple validation
		if c.ModelID == 0 || c.Content == "" {
			http.Error(w, "model_id and content required", http.StatusBadRequest)
			return
		}
		c.CreatedAt = time.Now()
		c.Upvotes = 0

		if err := gormApp.Create(&c).Error; err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.WriteJSON(w, http.StatusCreated, c, nil)
	})

	// Simple CORS middleware
	corsHandler := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      corsHandler(mux),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("Server starting on %s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func extractEngineOptions(trimNames []string) []string {
	seen := make(map[string]bool)
	var out []string
	for _, tn := range trimNames {
		opt := extractEngineOption(tn)
		if opt == "" {
			continue
		}
		key := strings.ToLower(opt)
		if seen[key] {
			continue
		}
		seen[key] = true
		out = append(out, opt)
	}
	sort.Strings(out)
	if len(out) > 20 {
		out = out[:20]
	}
	return out
}

func extractEngineOption(trimName string) string {
	s := strings.TrimSpace(trimName)
	if s == "" {
		return ""
	}
	parts := strings.Fields(s)
	// Find first token that looks like displacement (e.g. 1.4) or contains digits+dot
	for i := 0; i < len(parts); i++ {
		p := strings.Trim(parts[i], "()[]{}:;,")
		if looksLikeDisplacement(p) {
			// Prefer patterns like "1.4 TSI" / "1.5 eTSI" / "1.6 TDI"
			if i+1 < len(parts) {
				next := strings.Trim(parts[i+1], "()[]{}:;,")
				if looksLikeEngineFamily(next) {
					return p + " " + next
				}
			}
			return p
		}
	}
	return ""
}

func looksLikeDisplacement(token string) bool {
	if token == "" {
		return false
	}
	// crude check for things like 1.0, 1.2, 2.0
	if len(token) < 3 || len(token) > 5 {
		return false
	}
	if token[1] != '.' {
		return false
	}
	if token[0] < '0' || token[0] > '9' {
		return false
	}
	if token[2] < '0' || token[2] > '9' {
		return false
	}
	return true
}

func looksLikeEngineFamily(token string) bool {
	if token == "" {
		return false
	}
	upper := strings.ToUpper(token)
	switch upper {
	case "TSI", "TDI", "TFSI", "TGI", "ETSI", "EDRIVE", "TCE", "DCI", "CDI", "ECONETIC":
		return true
	}
	// Allow cases like eTSI or BlueHDi etc.
	if strings.Contains(upper, "TSI") || strings.Contains(upper, "TDI") || strings.Contains(upper, "ETSI") || strings.Contains(upper, "HDi") {
		return true
	}
	return false
}
