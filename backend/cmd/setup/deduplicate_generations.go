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

	// Enable foreign keys
	db.Exec("PRAGMA foreign_keys = ON")

	fmt.Println("🚀 Starting Generation Deduplication...")

	// 1. Find models with duplicate generation codes
	rows, err := db.Query(`
		SELECT model_id, code, COUNT(*) as c 
		FROM generations 
		GROUP BY model_id, code 
		HAVING c > 1
	`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	type DuplicateGroup struct {
		ModelID int
		Code    string
	}
	var groups []DuplicateGroup

	for rows.Next() {
		var g DuplicateGroup
		var count int
		rows.Scan(&g.ModelID, &g.Code, &count)
		groups = append(groups, g)
	}
	rows.Close()

	if len(groups) == 0 {
		fmt.Println("✅ No duplicates found.")
		return
	}

	for _, g := range groups {
		fmt.Printf("Processing duplicates for Model %d, Code %s...\n", g.ModelID, g.Code)

		// A. Fetch all candidates
		genRows, err := db.Query("SELECT id, name, image_url FROM generations WHERE model_id = ? AND code = ? ORDER BY id ASC", g.ModelID, g.Code)
		if err != nil {
			log.Printf("Error fetching gens: %v", err)
			continue
		}

		var ids []int
		var bestID int
		bestScore := -1

		for genRows.Next() {
			var id int
			var name string
			var img sql.NullString
			genRows.Scan(&id, &name, &img)
			ids = append(ids, id)

			// Scoring strategy:
			// 1. Has Image? (+10)
			// 2. Name contains year range? (+5)
			// 3. Lower ID? (+1)
			score := 0
			if img.Valid && img.String != "" {
				score += 10
			}
			if len(name) > 10 { // simplistic check for "Tip 8L (1996-2003)" vs "Tip 8L"
				score += 5
			}

			if score > bestScore {
				bestScore = score
				bestID = id
			}
		}
		genRows.Close()

		fmt.Printf("  -> Candidates: %v. Winner: %d (Score: %d)\n", ids, bestID, bestScore)

		// B. Merge Trims to Winner
		for _, id := range ids {
			if id == bestID {
				continue
			}
			// Update trims
			res, err := db.Exec("UPDATE trims SET generation_id = ? WHERE generation_id = ?", bestID, id)
			if err != nil {
				log.Printf("    ❌ Failed to re-link trims from %d: %v", id, err)
			} else {
				n, _ := res.RowsAffected()
				if n > 0 {
					fmt.Printf("    Moved %d trims from %d to %d.\n", n, id, bestID)
				}
			}

			// C. Delete Loser
			_, err = db.Exec("DELETE FROM generations WHERE id = ?", id)
			if err != nil {
				log.Printf("    ❌ Failed to delete dup gen %d: %v", id, err)
			} else {
				fmt.Printf("    Deleted duplicate generation %d.\n", id)
			}
		}
	}

	fmt.Println("✅ Deduplication Complete!")
}
