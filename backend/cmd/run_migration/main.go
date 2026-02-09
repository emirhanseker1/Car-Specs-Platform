package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

func main() {
	filePtr := flag.String("file", "", "Path to SQL file to execute")
	dbPtr := flag.String("db", "vehicles.db", "Path to SQLite DB")
	flag.Parse()

	if *filePtr == "" {
		log.Fatal("Please provide -file argument")
	}

	// 1. Read SQL file
	content, err := os.ReadFile(*filePtr)
	if err != nil {
		log.Fatalf("Failed to read SQL file: %v", err)
	}
	sqlContent := string(content)

	// 2. Connect to DB
	// Check if db path is relative or absolute
	dbPath := *dbPtr
	if !filepath.IsAbs(dbPath) {
		wd, _ := os.Getwd()
		dbPath = filepath.Join(wd, dbPath)
	}

	fmt.Printf("📂 Database: %s\n", dbPath)
	fmt.Printf("📜 Migration: %s\n", *filePtr)

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatalf("Failed to open DB: %v", err)
	}
	defer db.Close()

	// 3. Execute
	// Split by semi-colon to execute statement by statement if needed,
	// but mostly Exec handles multi-statement if supported by driver.
	// modernc.org/sqlite supports Exec with multiple statements.

	_, err = db.Exec(sqlContent)
	if err != nil {
		log.Fatalf("❌ Execution Failed: %v", err)
	}

	fmt.Println("✅ Migration executed successfully!")
}
