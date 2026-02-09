package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	_ "modernc.org/sqlite"
)

func main() {
	db, err := sql.Open("sqlite", "./vehicles.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// List of migration files to apply
	migrations := []string{
		"migrations/015_split_audi_a3_8v.sql",      // 8V1/8V2 split
		"migrations/017_replace_vw_golf_data.sql",  // VW Golf
		"migrations/018_populate_vw_passat.sql",    // VW Passat
		"migrations/019_populate_bmw_1_series.sql", // BMW 1 Series
		"migrations/020_populate_bmw_3_series.sql", // BMW 3 Series
		"migrations/021_populate_bmw_5_series.sql", // BMW 5 Series
	}

	for _, migFile := range migrations {
		fmt.Printf("Applying: %s\n", filepath.Base(migFile))

		content, err := ioutil.ReadFile(migFile)
		if err != nil {
			log.Printf("  ✗ Error reading file: %v", err)
			continue
		}

		// Execute migration
		_, err = db.Exec(string(content))
		if err != nil {
			// Try to extract just the error message
			errStr := err.Error()
			if strings.Contains(errStr, "UNIQUE constraint") || strings.Contains(errStr, "duplicate") {
				fmt.Printf("  ⚠ Skipped (data already exists)\n")
			} else {
				log.Printf("  ✗ Error: %v", err)
			}
			continue
		}
		fmt.Printf("  ✓ Applied successfully\n")
	}

	// Now fix the image paths to use actual file paths from /images/vehicles/
	fmt.Println("\nFixing image paths...")

	// Fix VW Golf generation images
	golfImageFixes := []struct {
		code     string
		imageURL string
	}{
		{"MK8", "/images/vehicles/volkswagen/mk8/vw-golf-8.png"},
		{"MK7.5", "/images/vehicles/volkswagen/mk7-5/vw-golf-7-5.png"},
		{"MK7", "/images/vehicles/volkswagen/mk7/vw-golf-7.png"},
		{"MK6", "/images/vehicles/volkswagen/mk6/vw-golf-6.png"},
		{"MK5", "/images/vehicles/volkswagen/mk5/vw-golf-5.png"},
	}

	for _, fix := range golfImageFixes {
		_, err := db.Exec(`UPDATE generations SET image_url = ? WHERE code = ?`, fix.imageURL, fix.code)
		if err == nil {
			fmt.Printf("  ✓ Golf %s: %s\n", fix.code, fix.imageURL)
		}
	}

	// Fix VW Passat generation images
	passatImageFixes := []struct {
		code     string
		imageURL string
	}{
		{"B8.5", "/images/vehicles/volkswagen/b8-5/vw-passat-b8-5.png"},
		{"B8", "/images/vehicles/volkswagen/b8/vw-passat-b8.png"},
		{"B7", "/images/vehicles/volkswagen/b7/vw-passat-b7.png"},
		{"B6", "/images/vehicles/volkswagen/b6/vw-passat-b6.png"},
		{"B5.5", "/images/vehicles/volkswagen/b5-5/vw-passat-b5-5.png"},
	}

	for _, fix := range passatImageFixes {
		_, err := db.Exec(`UPDATE generations SET image_url = ? WHERE code = ?`, fix.imageURL, fix.code)
		if err == nil {
			fmt.Printf("  ✓ Passat %s: %s\n", fix.code, fix.imageURL)
		}
	}

	// Fix BMW generation images - use actual folder structure
	bmwImageFixes := []struct {
		code     string
		imageURL string
	}{
		// 1 Series
		{"F40", "/images/vehicles/bmw/f40/bmw-1-f40.png"},
		{"F20 LCI", "/images/vehicles/bmw/f20-lci/bmw-1-f20-lci.png"},
		{"F20", "/images/vehicles/bmw/f20/bmw-1-f20.png"},
		{"E87 LCI", "/images/vehicles/bmw/e87-lci/bmw-1-e87-lci.png"},
		{"E87", "/images/vehicles/bmw/e87/bmw-1-e87.png"},
		// 3 Series
		{"G20 LCI", "/images/vehicles/bmw/g20-lci/bmw-3-g20-lci.png"},
		{"G20", "/images/vehicles/bmw/g20/bmw-3-g20.png"},
		{"F30 LCI", "/images/vehicles/bmw/f30-lci/bmw-3-f30-lci.png"},
		{"F30", "/images/vehicles/bmw/f30/bmw-3-f30.png"},
		{"E90 LCI", "/images/vehicles/bmw/e90-lci/bmw-3-e90-lci.png"},
		{"E90", "/images/vehicles/bmw/e90/bmw-3-e90.png"},
		// 5 Series
		{"G60", "/images/vehicles/bmw/g60/bmw-5-g60.png"},
		{"G30 LCI", "/images/vehicles/bmw/g30-lci/bmw-5-g30-lci.png"},
		{"G30", "/images/vehicles/bmw/g30/bmw-5-g30.png"},
		{"F10 LCI", "/images/vehicles/bmw/f10-lci/bmw-5-f10-lci.png"},
		{"F10", "/images/vehicles/bmw/f10/bmw-5-f10.png"},
		{"E60 LCI", "/images/vehicles/bmw/e60-lci/bmw-5-e60-lci.png"},
		{"E60", "/images/vehicles/bmw/e60/bmw-5-e60.png"},
	}

	for _, fix := range bmwImageFixes {
		_, err := db.Exec(`UPDATE generations SET image_url = ? WHERE code = ?`, fix.imageURL, fix.code)
		if err == nil {
			fmt.Printf("  ✓ BMW %s: %s\n", fix.code, fix.imageURL)
		}
	}

	// Fix Audi 8V1/8V2 images
	audiImageFixes := []struct {
		code     string
		imageURL string
	}{
		{"8V1", "/images/vehicles/audi/8v1/audi-a3-8v-2013-2016.png"},
		{"8V2", "/images/vehicles/audi/8v2/audi-a3-8v-2016-2020.png"},
	}

	for _, fix := range audiImageFixes {
		_, err := db.Exec(`UPDATE generations SET image_url = ? WHERE code = ?`, fix.imageURL, fix.code)
		if err == nil {
			fmt.Printf("  ✓ Audi A3 %s: %s\n", fix.code, fix.imageURL)
		}
	}

	// Verify data
	fmt.Println("\n=== Verification ===")

	// Count generations by brand
	rows, err := db.Query(`
		SELECT b.name, COUNT(g.id) as gen_count
		FROM brands b
		LEFT JOIN models m ON m.brand_id = b.id
		LEFT JOIN generations g ON g.model_id = m.id
		GROUP BY b.id
		ORDER BY b.name
	`)
	if err == nil {
		defer rows.Close()
		fmt.Println("Generations by brand:")
		for rows.Next() {
			var brand string
			var count int
			rows.Scan(&brand, &count)
			fmt.Printf("  %s: %d generations\n", brand, count)
		}
	}

	// List all generations with images
	fmt.Println("\nGenerations with images:")
	rows2, err := db.Query(`
		SELECT m.name, g.code, g.image_url 
		FROM generations g 
		JOIN models m ON g.model_id = m.id
		ORDER BY m.name, g.start_year DESC
	`)
	if err == nil {
		defer rows2.Close()
		for rows2.Next() {
			var model, code string
			var imageURL sql.NullString
			rows2.Scan(&model, &code, &imageURL)
			if imageURL.Valid && imageURL.String != "" {
				fmt.Printf("  %s %s: %s\n", model, code, imageURL.String)
			} else {
				fmt.Printf("  %s %s: (no image)\n", model, code)
			}
		}
	}
}
