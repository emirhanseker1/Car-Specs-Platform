package main

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

func main() {
	db, err := sql.Open("sqlite", "./vehicles.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	query := `
		UPDATE generations
		SET image_url = '/images/vehicles/volkswagen/b6/cover.png'
		WHERE model_id IN (SELECT id FROM models WHERE name = 'Passat') AND code = 'B6'
	`

	res, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}

	affected, _ := res.RowsAffected()
	log.Printf("Updated %d rows for Passat B6 image URL", affected)
}
