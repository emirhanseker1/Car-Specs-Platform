package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

func main() {
	log.Println("=== Engine Codes Migration Script ===")
	log.Println("This script will:")
	log.Println("  1. Backup the existing database")
	log.Println("  2. Create the 'engines' table")
	log.Println("  3. Add 'engine_id' column to 'trims' table")
	log.Println("  4. Create necessary indexes")
	log.Println("")

	// Get database path
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		cwd, err := os.Getwd()
		if err != nil {
			log.Fatalf("Failed to get working directory: %v", err)
		}
		// We're in backend/cmd/migrate_engines, go up to backend/
		dbPath = filepath.Join(cwd, "..", "..", "vehicles.db")
		dbPath, err = filepath.Abs(dbPath)
		if err != nil {
			log.Fatalf("Failed to get absolute path: %v", err)
		}
	}

	log.Printf("📁 Using database: %s\n", dbPath)

	// Check if database exists
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		log.Fatalf("❌ Database does not exist: %s", dbPath)
	}

	// Step 1: Backup database
	log.Println("\n🔄 Step 1: Creating backup...")
	backupPath := createBackup(dbPath)
	log.Printf("  ✅ Backup created: %s\n", backupPath)

	// Open database
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatalf("❌ Failed to open database: %v", err)
	}
	defer db.Close()

	// Enable foreign keys
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		log.Fatalf("❌ Failed to enable foreign keys: %v", err)
	}

	// Step 2: Create engines table
	log.Println("\n🏗️  Step 2: Creating 'engines' table...")
	if err := createEnginesTable(db); err != nil {
		log.Fatalf("❌ Failed to create engines table: %v", err)
	}
	log.Println("  ✅ 'engines' table created successfully")

	// Step 3: Add engine_id column to trims
	log.Println("\n🔧 Step 3: Adding 'engine_id' column to 'trims' table...")
	if err := addEngineIDColumn(db); err != nil {
		// This might fail if column already exists, which is OK
		log.Printf("  ⚠️  Note: %v (column might already exist)\n", err)
	} else {
		log.Println("  ✅ 'engine_id' column added successfully")
	}

	// Step 4: Create indexes
	log.Println("\n📊 Step 4: Creating indexes...")
	if err := createIndexes(db); err != nil {
		log.Fatalf("❌ Failed to create indexes: %v", err)
	}
	log.Println("  ✅ Indexes created successfully")

	// Step 5: Populate engines table with unique codes from trims
	log.Println("\n📥 Step 5: Populating 'engines' table with existing engine codes...")
	populated, err := populateEnginesFromTrims(db)
	if err != nil {
		log.Fatalf("❌ Failed to populate engines: %v", err)
	}
	log.Printf("  ✅ Populated %d unique engine codes\n", populated)

	// Step 6: Link trims to engines
	log.Println("\n🔗 Step 6: Linking trims to engines...")
	linked, err := linkTrimsToEngines(db)
	if err != nil {
		log.Fatalf("❌ Failed to link trims: %v", err)
	}
	log.Printf("  ✅ Linked %d trims to engines\n", linked)

	// Verification
	log.Println("\n✅ Migration completed successfully!")
	log.Println("\n📊 Verification:")
	verifyMigration(db)

	log.Println("\n💡 Next steps:")
	log.Println("  1. You can now manually add detailed engine information via the API")
	log.Println("  2. Use POST /api/engines to add engine details")
	log.Println("  3. Frontend will automatically display engine information")
}

func createBackup(dbPath string) string {
	timestamp := time.Now().Format("20060102_150405")
	backupPath := filepath.Join(
		filepath.Dir(dbPath),
		fmt.Sprintf("vehicles_backup_%s.db", timestamp),
	)

	// Read original database
	data, err := os.ReadFile(dbPath)
	if err != nil {
		log.Fatalf("Failed to read database: %v", err)
	}

	// Write backup
	if err := os.WriteFile(backupPath, data, 0644); err != nil {
		log.Fatalf("Failed to create backup: %v", err)
	}

	return backupPath
}

func createEnginesTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS engines (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		code TEXT NOT NULL UNIQUE,
		name TEXT,
		manufacturer TEXT,
		engine_type TEXT,
		fuel_type TEXT,
		displacement_cc INTEGER,
		cylinders INTEGER,
		cylinder_layout TEXT,
		valves_per_cylinder INTEGER,
		aspiration TEXT,
		
		description TEXT,
		technology_features TEXT,
		production_start_year INTEGER,
		production_end_year INTEGER,
		
		common_problems TEXT,
		solutions TEXT,
		maintenance_notes TEXT,
		reliability_rating INTEGER,
		
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`

	_, err := db.Exec(query)
	return err
}

func addEngineIDColumn(db *sql.DB) error {
	query := `ALTER TABLE trims ADD COLUMN engine_id INTEGER REFERENCES engines(id)`
	_, err := db.Exec(query)
	return err
}

func createIndexes(db *sql.DB) error {
	indexes := []string{
		`CREATE INDEX IF NOT EXISTS idx_engines_code ON engines(code)`,
		`CREATE INDEX IF NOT EXISTS idx_engines_fuel_type ON engines(fuel_type)`,
		`CREATE INDEX IF NOT EXISTS idx_engines_manufacturer ON engines(manufacturer)`,
		`CREATE INDEX IF NOT EXISTS idx_trims_engine_id ON trims(engine_id)`,
	}

	for _, idx := range indexes {
		if _, err := db.Exec(idx); err != nil {
			return err
		}
	}

	return nil
}

func populateEnginesFromTrims(db *sql.DB) (int, error) {
	// Get unique engine codes from trims where code is not null and not empty
	query := `
		SELECT DISTINCT 
			UPPER(TRIM(engine_code)) AS code,
			fuel_type,
			displacement_cc,
			cylinders,
			cylinder_layout,
			engine_type
		FROM trims
		WHERE engine_code IS NOT NULL 
		  AND TRIM(engine_code) != ''
		  AND UPPER(TRIM(engine_code)) NOT IN (SELECT code FROM engines)
		ORDER BY code
	`

	rows, err := db.Query(query)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	count := 0
	insertQuery := `
		INSERT INTO engines (
			code, fuel_type, displacement_cc, cylinders, 
			cylinder_layout, engine_type
		) VALUES (?, ?, ?, ?, ?, ?)
	`

	for rows.Next() {
		var code string
		var fuelType, cylinderLayout, engineType sql.NullString
		var displacementCC, cylinders sql.NullInt64

		if err := rows.Scan(&code, &fuelType, &displacementCC, &cylinders, &cylinderLayout, &engineType); err != nil {
			log.Printf("  ⚠️  Skipping row due to error: %v", err)
			continue
		}

		_, err := db.Exec(insertQuery, code,
			nullStringToPtr(fuelType),
			nullInt64ToPtr(displacementCC),
			nullInt64ToPtr(cylinders),
			nullStringToPtr(cylinderLayout),
			nullStringToPtr(engineType))

		if err != nil {
			log.Printf("  ⚠️  Failed to insert engine code '%s': %v", code, err)
			continue
		}

		count++
	}

	return count, nil
}

func linkTrimsToEngines(db *sql.DB) (int, error) {
	query := `
		UPDATE trims
		SET engine_id = (
			SELECT id FROM engines 
			WHERE engines.code = UPPER(TRIM(trims.engine_code))
		)
		WHERE engine_code IS NOT NULL 
		  AND TRIM(engine_code) != ''
		  AND engine_id IS NULL
	`

	result, err := db.Exec(query)
	if err != nil {
		return 0, err
	}

	rowsAffected, _ := result.RowsAffected()
	return int(rowsAffected), nil
}

func verifyMigration(db *sql.DB) {
	// Count engines
	var engineCount int
	db.QueryRow("SELECT COUNT(*) FROM engines").Scan(&engineCount)
	log.Printf("  • Total engines: %d", engineCount)

	// Count trims with engine_id
	var linkedTrims int
	db.QueryRow("SELECT COUNT(*) FROM trims WHERE engine_id IS NOT NULL").Scan(&linkedTrims)
	log.Printf("  • Trims linked to engines: %d", linkedTrims)

	// Count trims without engine_id
	var unlinkedTrims int
	db.QueryRow("SELECT COUNT(*) FROM trims WHERE engine_id IS NULL").Scan(&unlinkedTrims)
	log.Printf("  • Trims without engine link: %d", unlinkedTrims)

	// Check for orphaned references
	var orphaned int
	db.QueryRow(`
		SELECT COUNT(*) 
		FROM trims t 
		LEFT JOIN engines e ON t.engine_id = e.id 
		WHERE t.engine_id IS NOT NULL AND e.id IS NULL
	`).Scan(&orphaned)

	if orphaned > 0 {
		log.Printf("  ⚠️  WARNING: %d orphaned trim references found!", orphaned)
	} else {
		log.Printf("  • No orphaned references ✓")
	}
}

// Helper functions
func nullStringToPtr(ns sql.NullString) *string {
	if ns.Valid {
		return &ns.String
	}
	return nil
}

func nullInt64ToPtr(ni sql.NullInt64) *int {
	if ni.Valid {
		val := int(ni.Int64)
		return &val
	}
	return nil
}
