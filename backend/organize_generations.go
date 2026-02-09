package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	_ "modernc.org/sqlite"
)

type Generation struct {
	ID        int
	BrandName string
	Code      string
	ImageURL  sql.NullString
}

func main() {
	publicImagesDir := "../frontend/public" // Base public dir to resolve relative paths

	db, err := sql.Open("sqlite", "./vehicles.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 1. Get all generations with logos
	query := `
		SELECT g.id, b.name as brand_name, g.code, g.image_url 
		FROM generations g
		JOIN models m ON g.model_id = m.id
		JOIN brands b ON m.brand_id = b.id
		WHERE g.image_url IS NOT NULL AND g.image_url != ''
	`
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var jobs []Generation
	for rows.Next() {
		var g Generation
		if err := rows.Scan(&g.ID, &g.BrandName, &g.Code, &g.ImageURL); err != nil {
			log.Fatal(err)
		}
		jobs = append(jobs, g)
	}

	// Regex for slugifying code
	reg, _ := regexp.Compile("[^a-zA-Z0-9]+")

	fmt.Printf("📦 Processing %d generations...\n", len(jobs))

	for _, g := range jobs {
		if !g.ImageURL.Valid {
			continue
		}
		currentRelPath := g.ImageURL.String
		// Construct absolute current path
		// currentRelPath usually starts with /images/...
		currentAbsPath := filepath.Join(publicImagesDir, strings.TrimPrefix(currentRelPath, "/"))

		// Determine expected brand folder name (already organized)
		// Usually 'audi', 'bmw', 'volkswagen'
		brandFolder := strings.ToLower(g.BrandName)
		if brandFolder == "vw" {
			brandFolder = "volkswagen"
		}

		// Create Generation Slug
		genSlug := strings.ToLower(reg.ReplaceAllString(g.Code, "-"))
		genSlug = strings.Trim(genSlug, "-")

		// Target Directory: images/vehicles/[brand]/[genSlug]
		targetDirRel := filepath.Join("images", "vehicles", brandFolder, genSlug)
		targetDirAbs := filepath.Join(publicImagesDir, targetDirRel)

		if _, err := os.Stat(targetDirAbs); os.IsNotExist(err) {
			os.MkdirAll(targetDirAbs, 0755)
		}

		// New File Path
		fileName := filepath.Base(currentAbsPath)
		targetFileAbs := filepath.Join(targetDirAbs, fileName)
		targetFileRel := "/" + filepath.ToSlash(filepath.Join(targetDirRel, fileName))

		// Move File
		if currentAbsPath != targetFileAbs {
			// Check if source exists
			if _, err := os.Stat(currentAbsPath); err == nil {
				err := os.Rename(currentAbsPath, targetFileAbs)
				if err != nil {
					log.Printf("❌ Failed to move %s: %v", fileName, err)
					continue
				}
				fmt.Printf("Moved: %s -> %s\n", fileName, targetFileRel)

				// Update DB
				_, err = db.Exec("UPDATE generations SET image_url = ? WHERE id = ?", targetFileRel, g.ID)
				if err != nil {
					log.Printf("Failed to update DB for gen %d: %v", g.ID, err)
				}
			} else {
				// Log but don't error, maybe already moved or path mismatch
				// log.Printf("⚠️  Source file not found: %s", currentAbsPath)
			}
		}
	}

	fmt.Println("✅ Generation Organization Complete!")
}
