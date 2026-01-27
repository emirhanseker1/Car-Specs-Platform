package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	dbPath := "./vehicles.db"
	migrationsDir := "./migrations"

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Enable foreign keys
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		log.Fatalf("Failed to enable foreign keys: %v", err)
	}

	// Get all migration files
	files, err := filepath.Glob(filepath.Join(migrationsDir, "*.sql"))
	if err != nil {
		log.Fatalf("Failed to read migrations directory: %v", err)
	}

	sort.Strings(files)

	fmt.Println("Running migrations...")
	for _, file := range files {
		fmt.Printf("Applying: %s\n", filepath.Base(file))

		content, err := os.ReadFile(file)
		if err != nil {
			log.Fatalf("Failed to read migration file %s: %v", file, err)
		}

		if _, err := db.Exec(string(content)); err != nil {
			log.Fatalf("Failed to execute migration %s: %v", file, err)
		}
	}

	fmt.Println("✓ All migrations applied successfully")
}
