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
	// Paths
	publicImagesDir := "../frontend/public/images"
	brandsDir := filepath.Join(publicImagesDir, "brands")
	generationsDir := filepath.Join(publicImagesDir, "generations")
	vehiclesDir := filepath.Join(publicImagesDir, "vehicles")

	// 1. Move Brand Logos: images/brands/*.png -> images/*.png (Idempotent)
	fmt.Println("📦 Checking Brand Logos...")
	if _, err := os.Stat(brandsDir); err == nil {
		files, err := os.ReadDir(brandsDir)
		if err == nil {
			for _, file := range files {
				if !file.IsDir() && strings.HasSuffix(file.Name(), ".png") {
					oldPath := filepath.Join(brandsDir, file.Name())
					newPath := filepath.Join(publicImagesDir, file.Name())
					err := os.Rename(oldPath, newPath)
					if err != nil {
						log.Printf("Failed to move %s: %v", file.Name(), err)
					} else {
						fmt.Printf("Moved %s -> %s\n", oldPath, newPath)
					}
				}
			}
			os.Remove(brandsDir)
		}
	} else {
		fmt.Println("Brand logos directory not found or already moved.")
	}

	// 2. Move Generations Files: images/generations/* -> images/vehicles/*
	fmt.Println("📦 Checking Generations Folder...")

	// Ensure target directory exists
	if _, err := os.Stat(vehiclesDir); os.IsNotExist(err) {
		os.MkdirAll(vehiclesDir, 0755)
	}

	if _, err := os.Stat(generationsDir); err == nil {
		files, err := os.ReadDir(generationsDir)
		if err == nil {
			for _, file := range files {
				oldPath := filepath.Join(generationsDir, file.Name())
				newPath := filepath.Join(vehiclesDir, file.Name())

				err := os.Rename(oldPath, newPath)
				if err != nil {
					log.Printf("Failed to move %s: %v", file.Name(), err)
				} else {
					fmt.Printf("Moved %s -> %s\n", oldPath, newPath)
				}
			}
			// Try removing the now-empty generations directory
			os.Remove(generationsDir)
		} else {
			log.Printf("Could not list generations: %v", err)
		}
	} else {
		fmt.Println("Generations directory not found (already merged?)")
	}

	// 3. Update Database
	fmt.Println("💾 Updating Database Records...")
	db, err := sql.Open("sqlite", "./vehicles.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Update Brands (Idempotent)
	_, err = db.Exec(`UPDATE brands SET logo_url = REPLACE(logo_url, '/images/brands/', '/images/')`)
	if err != nil {
		log.Printf("Error updating brands: %v", err)
	}

	// Update Generations (Idempotent)
	_, err = db.Exec(`UPDATE generations SET image_url = REPLACE(image_url, '/images/generations/', '/images/vehicles/')`)
	if err != nil {
		log.Printf("Error updating generations: %v", err)
	}

	// Update Trims
	_, err = db.Exec(`UPDATE trims SET image_url = REPLACE(image_url, '/images/generations/', '/images/vehicles/')`)
	if err != nil {
		log.Printf("Error updating trims: %v", err)
	}

	fmt.Println("✅ Restructuring Complete!")
}
