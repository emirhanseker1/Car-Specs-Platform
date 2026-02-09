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

	fmt.Println("=== DATABASE ANALYSIS ===\n")

	// 1. List all tables
	fmt.Println("1. TABLES IN DATABASE:")
	rows, _ := db.Query(`SELECT name FROM sqlite_master WHERE type='table' ORDER BY name`)
	for rows.Next() {
		var name string
		rows.Scan(&name)
		fmt.Printf("   - %s\n", name)
	}
	rows.Close()

	// 2. Count records in each table
	fmt.Println("\n2. RECORD COUNTS:")
	tables := []string{"brands", "models", "generations", "trims", "features"}
	for _, t := range tables {
		var count int
		db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s", t)).Scan(&count)
		fmt.Printf("   %s: %d records\n", t, count)
	}

	// 3. Brands
	fmt.Println("\n3. BRANDS:")
	rows, _ = db.Query(`SELECT id, name FROM brands ORDER BY name`)
	for rows.Next() {
		var id int
		var name string
		rows.Scan(&id, &name)
		fmt.Printf("   [%d] %s\n", id, name)
	}
	rows.Close()

	// 4. Models by brand
	fmt.Println("\n4. MODELS BY BRAND:")
	rows, _ = db.Query(`
		SELECT b.name, m.id, m.name, m.body_style 
		FROM models m 
		JOIN brands b ON m.brand_id = b.id 
		ORDER BY b.name, m.name
	`)
	for rows.Next() {
		var brandName, modelName, bodyStyle string
		var modelID int
		rows.Scan(&brandName, &modelID, &modelName, &bodyStyle)
		fmt.Printf("   %s > [%d] %s (%s)\n", brandName, modelID, modelName, bodyStyle)
	}
	rows.Close()

	// 5. Generations by model
	fmt.Println("\n5. GENERATIONS BY MODEL:")
	rows, _ = db.Query(`
		SELECT m.name, g.id, g.code, g.name, g.start_year, g.end_year
		FROM generations g 
		JOIN models m ON g.model_id = m.id 
		ORDER BY m.name, g.start_year DESC
	`)
	for rows.Next() {
		var modelName, code, genName string
		var genID, startYear int
		var endYear sql.NullInt64
		rows.Scan(&modelName, &genID, &code, &genName, &startYear, &endYear)
		endStr := "now"
		if endYear.Valid {
			endStr = fmt.Sprintf("%d", endYear.Int64)
		}
		fmt.Printf("   %s > [%d] %s: %s (%d-%s)\n", modelName, genID, code, genName, startYear, endStr)
	}
	rows.Close()

	// 6. Trims - sample
	fmt.Println("\n6. TRIMS (SAMPLE - first 20):")
	rows, _ = db.Query(`
		SELECT t.id, g.code, t.name, t.power_hp, t.torque_nm, t.transmission_type
		FROM trims t 
		JOIN generations g ON t.generation_id = g.id 
		ORDER BY t.id
		LIMIT 20
	`)
	for rows.Next() {
		var trimID int
		var genCode, trimName string
		var powerHP, torqueNM sql.NullInt64
		var transmission sql.NullString
		rows.Scan(&trimID, &genCode, &trimName, &powerHP, &torqueNM, &transmission)
		hp := "-"
		if powerHP.Valid {
			hp = fmt.Sprintf("%d HP", powerHP.Int64)
		}
		tx := "-"
		if transmission.Valid {
			tx = transmission.String
		}
		fmt.Printf("   [%d] %s > %s | %s | %s\n", trimID, genCode, trimName, hp, tx)
	}
	rows.Close()

	// 7. Check trims for specific generation
	fmt.Println("\n7. TRIMS FOR VW GOLF MK7.5:")
	var mk75ID int
	err = db.QueryRow(`SELECT id FROM generations WHERE code = 'MK7.5'`).Scan(&mk75ID)
	if err != nil {
		fmt.Printf("   MK7.5 not found: %v\n", err)
	} else {
		fmt.Printf("   Generation ID: %d\n", mk75ID)
		rows, _ = db.Query(`SELECT id, name, power_hp FROM trims WHERE generation_id = ?`, mk75ID)
		count := 0
		for rows.Next() {
			var id int
			var name string
			var hp sql.NullInt64
			rows.Scan(&id, &name, &hp)
			count++
			fmt.Printf("   [%d] %s | %v HP\n", id, name, hp.Int64)
		}
		rows.Close()
		fmt.Printf("   Total trims for MK7.5: %d\n", count)
	}

	// 8. Check schema of trims table
	fmt.Println("\n8. TRIMS TABLE SCHEMA:")
	rows, _ = db.Query(`PRAGMA table_info(trims)`)
	for rows.Next() {
		var cid int
		var name, ctype string
		var notnull, pk int
		var dflt sql.NullString
		rows.Scan(&cid, &name, &ctype, &notnull, &dflt, &pk)
		fmt.Printf("   %s (%s)\n", name, ctype)
	}
	rows.Close()
}
