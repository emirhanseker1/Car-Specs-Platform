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

	fmt.Println("=== ADDING MISSING TRIMS ===\n")

	// Insert VW Golf trims
	fmt.Println("VW Golf MK8 (ID 5):")
	insertGolfMK8Trims(db, 5)

	fmt.Println("\nVW Golf MK7.5 (ID 6):")
	insertGolfMK75Trims(db, 6)

	fmt.Println("\nVW Golf MK7 (ID 7):")
	insertGolfMK7Trims(db, 7)

	fmt.Println("\nVW Golf MK6 (ID 8):")
	insertGolfMK6Trims(db, 8)

	fmt.Println("\nVW Golf MK5 (ID 9):")
	insertGolfMK5Trims(db, 9)

	// VW Passat
	fmt.Println("\nVW Passat B8.5 (ID 10):")
	insertPassatB85Trims(db, 10)

	fmt.Println("\nVW Passat B8 (ID 11):")
	insertPassatB8Trims(db, 11)

	// Verify
	fmt.Println("\n=== VERIFICATION ===")
	rows, _ := db.Query(`
		SELECT g.code, COUNT(t.id) 
		FROM generations g 
		LEFT JOIN trims t ON t.generation_id = g.id 
		GROUP BY g.id 
		ORDER BY g.id
	`)
	for rows.Next() {
		var code string
		var count int
		rows.Scan(&code, &count)
		fmt.Printf("  %s: %d trims\n", code, count)
	}
	rows.Close()
}

func insertTrim(db *sql.DB, genID int, name string, powerHP, torqueNM int, transmission, fuel, drivetrain string, acceleration float64) {
	// Get model_id from generation
	var modelID int
	db.QueryRow("SELECT model_id FROM generations WHERE id = ?", genID).Scan(&modelID)

	_, err := db.Exec(`
		INSERT INTO trims (generation_id, model_id, name, power_hp, torque_nm, transmission_type, fuel_type, drivetrain, acceleration_0_100)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, genID, modelID, name, powerHP, torqueNM, transmission, fuel, drivetrain, acceleration)
	if err != nil {
		fmt.Printf("  ✗ %s: %v\n", name, err)
	} else {
		fmt.Printf("  ✓ %s\n", name)
	}
}

func insertGolfMK8Trims(db *sql.DB, genID int) {
	insertTrim(db, genID, "1.0 eTSI Life", 110, 200, "Otomatik (DSG)", "Mild Hybrid", "Önden Çekiş", 10.2)
	insertTrim(db, genID, "1.5 eTSI Style", 150, 250, "Otomatik (DSG)", "Mild Hybrid", "Önden Çekiş", 8.5)
	insertTrim(db, genID, "2.0 TDI Life", 115, 300, "Otomatik (DSG)", "Dizel", "Önden Çekiş", 10.0)
	insertTrim(db, genID, "GTI 2.0 TSI", 245, 370, "Otomatik (DSG)", "Benzin", "Önden Çekiş", 6.3)
	insertTrim(db, genID, "R 2.0 TSI 4Motion", 320, 420, "Otomatik (DSG)", "Benzin", "4WD", 4.7)
}

func insertGolfMK75Trims(db *sql.DB, genID int) {
	insertTrim(db, genID, "1.0 TSI Comfortline", 110, 200, "Otomatik (DSG)", "Benzin", "Önden Çekiş", 9.9)
	insertTrim(db, genID, "1.4 TSI Comfortline", 125, 200, "Otomatik (DSG)", "Benzin", "Önden Çekiş", 9.1)
	insertTrim(db, genID, "1.5 TSI EVO Highline", 150, 250, "Otomatik (DSG)", "Benzin", "Önden Çekiş", 8.3)
	insertTrim(db, genID, "1.6 TDI Comfortline", 115, 250, "Otomatik (DSG)", "Dizel", "Önden Çekiş", 10.2)
	insertTrim(db, genID, "GTI 2.0 TSI", 230, 350, "Otomatik (DSG)", "Benzin", "Önden Çekiş", 6.4)
	insertTrim(db, genID, "R 2.0 TSI 4Motion", 310, 400, "Otomatik (DSG)", "Benzin", "4WD", 4.6)
}

func insertGolfMK7Trims(db *sql.DB, genID int) {
	insertTrim(db, genID, "1.2 TSI Midline Plus", 105, 175, "Otomatik (DSG)", "Benzin", "Önden Çekiş", 10.2)
	insertTrim(db, genID, "1.4 TSI Comfortline", 122, 200, "Otomatik (DSG)", "Benzin", "Önden Çekiş", 9.3)
	insertTrim(db, genID, "1.4 TSI ACT Highline", 140, 250, "Otomatik (DSG)", "Benzin", "Önden Çekiş", 8.4)
	insertTrim(db, genID, "1.6 TDI Comfortline", 105, 250, "Otomatik (DSG)", "Dizel", "Önden Çekiş", 10.7)
	insertTrim(db, genID, "GTD 2.0 TDI", 184, 380, "Otomatik (DSG)", "Dizel", "Önden Çekiş", 7.5)
	insertTrim(db, genID, "GTI 2.0 TSI", 220, 350, "Otomatik (DSG)", "Benzin", "Önden Çekiş", 6.5)
	insertTrim(db, genID, "R 2.0 TSI 4Motion", 300, 380, "Otomatik (DSG)", "Benzin", "4WD", 4.9)
}

func insertGolfMK6Trims(db *sql.DB, genID int) {
	insertTrim(db, genID, "1.6 Primeline", 102, 148, "Otomatik (DSG)", "Benzin", "Önden Çekiş", 11.3)
	insertTrim(db, genID, "1.4 TSI Comfortline", 122, 200, "Otomatik (DSG)", "Benzin", "Önden Çekiş", 9.5)
	insertTrim(db, genID, "1.4 TSI Highline (Twincharger)", 160, 240, "Otomatik (DSG)", "Benzin", "Önden Çekiş", 8.0)
	insertTrim(db, genID, "1.6 TDI Comfortline", 105, 250, "Otomatik (DSG)", "Dizel", "Önden Çekiş", 11.3)
	insertTrim(db, genID, "GTI 2.0 TSI", 210, 280, "Otomatik (DSG)", "Benzin", "Önden Çekiş", 6.9)
	insertTrim(db, genID, "R 2.0 TSI 4Motion", 270, 350, "Otomatik (DSG)", "Benzin", "4WD", 5.5)
}

func insertGolfMK5Trims(db *sql.DB, genID int) {
	insertTrim(db, genID, "1.6 Primeline", 102, 148, "Otomatik (Tiptronic)", "Benzin", "Önden Çekiş", 12.5)
	insertTrim(db, genID, "1.6 FSI Comfortline", 115, 155, "Otomatik (Tiptronic)", "Benzin", "Önden Çekiş", 11.5)
	insertTrim(db, genID, "1.4 TSI GT", 170, 240, "Otomatik (DSG)", "Benzin", "Önden Çekiş", 7.9)
	insertTrim(db, genID, "1.9 TDI Comfortline", 105, 250, "Otomatik (DSG)", "Dizel", "Önden Çekiş", 11.3)
	insertTrim(db, genID, "GTI 2.0 TFSI", 200, 280, "Otomatik (DSG)", "Benzin", "Önden Çekiş", 6.9)
	insertTrim(db, genID, "R32 3.2 V6 4Motion", 250, 320, "Otomatik (DSG)", "Benzin", "4WD", 6.2)
}

func insertPassatB85Trims(db *sql.DB, genID int) {
	insertTrim(db, genID, "1.5 TSI EVO Elegance", 150, 250, "Otomatik (DSG)", "Benzin", "Önden Çekiş", 8.7)
	insertTrim(db, genID, "2.0 TSI 4Motion", 190, 320, "Otomatik (DSG)", "Benzin", "4WD", 7.3)
	insertTrim(db, genID, "2.0 TDI Business", 150, 360, "Otomatik (DSG)", "Dizel", "Önden Çekiş", 8.8)
	insertTrim(db, genID, "2.0 TDI 4Motion", 190, 400, "Otomatik (DSG)", "Dizel", "4WD", 7.6)
}

func insertPassatB8Trims(db *sql.DB, genID int) {
	insertTrim(db, genID, "1.4 TSI ACT Comfortline", 150, 250, "Otomatik (DSG)", "Benzin", "Önden Çekiş", 8.4)
	insertTrim(db, genID, "1.8 TSI Highline", 180, 280, "Otomatik (DSG)", "Benzin", "Önden Çekiş", 7.9)
	insertTrim(db, genID, "2.0 TSI R-Line", 220, 350, "Otomatik (DSG)", "Benzin", "Önden Çekiş", 6.7)
	insertTrim(db, genID, "1.6 TDI Comfortline", 120, 250, "Otomatik (DSG)", "Dizel", "Önden Çekiş", 10.5)
	insertTrim(db, genID, "2.0 TDI Highline", 150, 340, "Otomatik (DSG)", "Dizel", "Önden Çekiş", 8.7)
	insertTrim(db, genID, "2.0 BiTDI 4Motion", 240, 500, "Otomatik (DSG)", "Dizel", "4WD", 6.1)
}
