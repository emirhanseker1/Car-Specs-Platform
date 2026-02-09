package main

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"runtime"

	_ "modernc.org/sqlite"
)

func main() {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	dbPath := filepath.Join(basepath, "vehicles.db")

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, name, start_year, end_year, code FROM generations WHERE model_id = 1 ORDER BY start_year")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	fmt.Println("--- A3 GENERATIONS (Simple) ---")
	for rows.Next() {
		var id int
		var start, end sql.NullInt64
		var name, code string
		rows.Scan(&id, &name, &start, &end, &code)

		s := "NULL"
		if start.Valid {
			s = fmt.Sprintf("%d", start.Int64)
		}
		e := "NULL"
		if end.Valid {
			e = fmt.Sprintf("%d", end.Int64)
		}

		fmt.Printf("ID: %d | Code: %s | Years: %s-%s | Name: %s\n", id, code, s, e, name)
	}
}
