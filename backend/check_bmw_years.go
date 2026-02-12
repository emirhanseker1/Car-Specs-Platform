package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./vehicles.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	query := `
		SELECT g.name, g.code, g.start_year, g.end_year, m.name as model_name
		FROM generations g
		JOIN models m ON g.model_id = m.id
		JOIN brands b ON m.brand_id = b.id
		WHERE b.name = 'BMW'
		ORDER BY m.name, g.start_year DESC
	`

	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	fmt.Println("BMW Generations:")
	fmt.Println("================")
	for rows.Next() {
		var name, code, modelName string
		var startYear, endYear int
		err := rows.Scan(&name, &code, &startYear, &endYear, &modelName)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s (%s): %d - %d (Model: %s)\n", name, code, startYear, endYear, modelName)
	}
}
