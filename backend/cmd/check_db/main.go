package main

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

func main() {
	databases := []string{
		"../../vehicles.db",
		"vehicles.db",
	}

	for _, dbPath := range databases {
		fmt.Printf("\n=== Checking: %s ===\n", dbPath)

		db, err := sql.Open("sqlite", dbPath)
		if err != nil {
			fmt.Printf("  ❌ Failed to open: %v\n", err)
			continue
		}
		defer db.Close()

		var trimCount, modelCount, brandCount int

		db.QueryRow("SELECT COUNT(*) FROM trims").Scan(&trimCount)
		db.QueryRow("SELECT COUNT(*) FROM models").Scan(&modelCount)
		db.QueryRow("SELECT COUNT(*) FROM brands").Scan(&brandCount)

		fmt.Printf("  ✓ Brands: %d\n", brandCount)
		fmt.Printf("  ✓ Models: %d\n", modelCount)
		fmt.Printf("  ✓ Trims:  %d\n", trimCount)

		// List all brands
		if brandCount > 0 {
			rows, _ := db.Query("SELECT name FROM brands ORDER BY name")
			fmt.Printf("  Brands: ")
			var brands []string
			for rows.Next() {
				var name string
				rows.Scan(&name)
				brands = append(brands, name)
			}
			rows.Close()
			fmt.Printf("%v\n", brands)
		}
	}
}
