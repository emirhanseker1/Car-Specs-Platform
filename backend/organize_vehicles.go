package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	_ "modernc.org/sqlite"
)

func main() {
	publicImagesDir := "../frontend/public/images"
	vehiclesDir := filepath.Join(publicImagesDir, "vehicles")

	// 1. Move Files
	fmt.Println("📦 Organizing Vehicles into Brand Folders...")
	files, err := os.ReadDir(vehiclesDir)
	if err != nil {
		log.Fatalf("Failed to read vehicles dir: %v", err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		name := strings.ToLower(file.Name())
		var targetBrand string

		if strings.HasPrefix(name, "audi") {
			targetBrand = "audi"
		} else if strings.HasPrefix(name, "bmw") {
			targetBrand = "bmw"
		} else if strings.HasPrefix(name, "vw") || strings.HasPrefix(name, "volkswagen") {
			targetBrand = "volkswagen"
		}

		if targetBrand != "" {
			targetDir := filepath.Join(vehiclesDir, targetBrand)
			if _, err := os.Stat(targetDir); os.IsNotExist(err) {
				os.MkdirAll(targetDir, 0755)
			}

			oldPath := filepath.Join(vehiclesDir, file.Name())
			newPath := filepath.Join(targetDir, file.Name())

			err := os.Rename(oldPath, newPath)
			if err != nil {
				log.Printf("Failed to move %s: %v", file.Name(), err)
			} else {
				fmt.Printf("Moved %s -> %s\n", oldPath, newPath)
			}
		}
	}

	// 2. Update Database
	fmt.Println("💾 Updating Database Records...")
	db, err := sql.Open("sqlite", "./vehicles.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Update Audi
	_, err = db.Exec(`UPDATE generations SET image_url = REPLACE(image_url, '/vehicles/audi-', '/vehicles/audi/audi-') WHERE image_url LIKE '%/vehicles/audi-%'`)
	_, err = db.Exec(`UPDATE trims SET image_url = REPLACE(image_url, '/vehicles/audi-', '/vehicles/audi/audi-') WHERE image_url LIKE '%/vehicles/audi-%'`)

	// Update BMW
	_, err = db.Exec(`UPDATE generations SET image_url = REPLACE(image_url, '/vehicles/bmw-', '/vehicles/bmw/bmw-') WHERE image_url LIKE '%/vehicles/bmw-%'`)
	_, err = db.Exec(`UPDATE trims SET image_url = REPLACE(image_url, '/vehicles/bmw-', '/vehicles/bmw/bmw-') WHERE image_url LIKE '%/vehicles/bmw-%'`)

	// Update VW
	_, err = db.Exec(`UPDATE generations SET image_url = REPLACE(image_url, '/vehicles/vw-', '/vehicles/volkswagen/vw-') WHERE image_url LIKE '%/vehicles/vw-%'`)
	_, err = db.Exec(`UPDATE trims SET image_url = REPLACE(image_url, '/vehicles/vw-', '/vehicles/volkswagen/vw-') WHERE image_url LIKE '%/vehicles/vw-%'`)

	fmt.Println("✅ Vehicle Organization Complete!")
}
