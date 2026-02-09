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

	fmt.Println("=== ADDING BMW TRIMS ===\n")

	// BMW 1 Series
	addBMW1SeriesTrims(db)

	// BMW 3 Series
	addBMW3SeriesTrims(db)

	// BMW 5 Series
	addBMW5SeriesTrims(db)

	// Additional VW Passat
	addPassatB7B6Trims(db)

	// Audi 8V1/8V2
	addAudi8VTrims(db)

	// Verify
	fmt.Println("\n=== FINAL VERIFICATION ===")
	rows, _ := db.Query(`
		SELECT m.name, g.code, COUNT(t.id) as trim_count
		FROM generations g 
		JOIN models m ON g.model_id = m.id
		LEFT JOIN trims t ON t.generation_id = g.id 
		GROUP BY g.id 
		HAVING trim_count > 0
		ORDER BY m.name, g.start_year DESC
	`)
	for rows.Next() {
		var model, code string
		var count int
		rows.Scan(&model, &code, &count)
		fmt.Printf("  %s %s: %d trims\n", model, code, count)
	}
	rows.Close()

	// Total
	var total int
	db.QueryRow("SELECT COUNT(*) FROM trims").Scan(&total)
	fmt.Printf("\nTotal trims in database: %d\n", total)
}

func insertTrim(db *sql.DB, genID int, name string, powerHP, torqueNM int, transmission, fuel, drivetrain string, accel float64) {
	var modelID int
	db.QueryRow("SELECT model_id FROM generations WHERE id = ?", genID).Scan(&modelID)
	_, err := db.Exec(`
		INSERT INTO trims (generation_id, model_id, name, power_hp, torque_nm, transmission_type, fuel_type, drivetrain, acceleration_0_100)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, genID, modelID, name, powerHP, torqueNM, transmission, fuel, drivetrain, accel)
	if err != nil {
		fmt.Printf("  ✗ %s: %v\n", name, err)
	} else {
		fmt.Printf("  ✓ %s\n", name)
	}
}

func addBMW1SeriesTrims(db *sql.DB) {
	fmt.Println("BMW 1 Serisi F40 (ID 15):")
	insertTrim(db, 15, "118i M Sport", 140, 220, "Otomatik (DCT)", "Benzin", "Önden Çekiş", 8.5)
	insertTrim(db, 15, "120i M Sport", 178, 280, "Otomatik (DCT)", "Benzin", "Önden Çekiş", 7.5)
	insertTrim(db, 15, "M135i xDrive", 306, 450, "Otomatik (DCT)", "Benzin", "4WD", 4.8)
	insertTrim(db, 15, "116d", 116, 270, "Otomatik (DCT)", "Dizel", "Önden Çekiş", 10.2)
	insertTrim(db, 15, "118d", 150, 350, "Otomatik (DCT)", "Dizel", "Önden Çekiş", 8.5)

	fmt.Println("\nBMW 1 Serisi F20 LCI (ID 16):")
	insertTrim(db, 16, "118i", 136, 220, "Otomatik (ZF)", "Benzin", "Arkadan İtiş", 8.7)
	insertTrim(db, 16, "120i", 177, 250, "Otomatik (ZF)", "Benzin", "Arkadan İtiş", 7.3)
	insertTrim(db, 16, "M140i", 340, 500, "Otomatik (ZF)", "Benzin", "Arkadan İtiş", 4.6)
	insertTrim(db, 16, "116d", 116, 270, "Otomatik (ZF)", "Dizel", "Arkadan İtiş", 10.3)
	insertTrim(db, 16, "118d", 150, 320, "Otomatik (ZF)", "Dizel", "Arkadan İtiş", 8.4)

	fmt.Println("\nBMW 1 Serisi F20 (ID 17):")
	insertTrim(db, 17, "116i", 136, 220, "Otomatik (ZF)", "Benzin", "Arkadan İtiş", 8.7)
	insertTrim(db, 17, "118i", 170, 250, "Otomatik (ZF)", "Benzin", "Arkadan İtiş", 7.4)
	insertTrim(db, 17, "M135i", 320, 450, "Otomatik (ZF)", "Benzin", "Arkadan İtiş", 4.9)
	insertTrim(db, 17, "116d", 116, 260, "Otomatik (ZF)", "Dizel", "Arkadan İtiş", 10.5)

	fmt.Println("\nBMW 1 Serisi E87 LCI (ID 18):")
	insertTrim(db, 18, "116i", 122, 160, "Otomatik", "Benzin", "Arkadan İtiş", 10.8)
	insertTrim(db, 18, "118i", 143, 190, "Otomatik", "Benzin", "Arkadan İtiş", 9.4)
	insertTrim(db, 18, "130i", 265, 315, "Otomatik", "Benzin", "Arkadan İtiş", 6.1)
	insertTrim(db, 18, "118d", 143, 300, "Otomatik", "Dizel", "Arkadan İtiş", 9.0)

	fmt.Println("\nBMW 1 Serisi E87 (ID 19):")
	insertTrim(db, 19, "116i", 115, 150, "Otomatik", "Benzin", "Arkadan İtiş", 10.8)
	insertTrim(db, 19, "118i", 129, 180, "Otomatik", "Benzin", "Arkadan İtiş", 9.6)
	insertTrim(db, 19, "120i", 150, 200, "Otomatik", "Benzin", "Arkadan İtiş", 8.7)
	insertTrim(db, 19, "130i", 265, 315, "Otomatik", "Benzin", "Arkadan İtiş", 6.1)
}

func addBMW3SeriesTrims(db *sql.DB) {
	fmt.Println("\nBMW 3 Serisi G20 LCI (ID 20):")
	insertTrim(db, 20, "318i", 156, 250, "Otomatik (ZF)", "Benzin", "Arkadan İtiş", 8.4)
	insertTrim(db, 20, "320i", 184, 300, "Otomatik (ZF)", "Benzin", "Arkadan İtiş", 7.1)
	insertTrim(db, 20, "330i", 245, 400, "Otomatik (ZF)", "Benzin", "Arkadan İtiş", 5.8)
	insertTrim(db, 20, "M340i xDrive", 374, 500, "Otomatik (ZF)", "Benzin", "4WD", 4.4)
	insertTrim(db, 20, "320d", 190, 400, "Otomatik (ZF)", "Dizel", "Arkadan İtiş", 6.8)

	fmt.Println("\nBMW 3 Serisi G20 (ID 21):")
	insertTrim(db, 21, "318i", 156, 250, "Otomatik (ZF)", "Benzin", "Arkadan İtiş", 8.4)
	insertTrim(db, 21, "320i", 184, 300, "Otomatik (ZF)", "Benzin", "Arkadan İtiş", 7.1)
	insertTrim(db, 21, "330i", 258, 400, "Otomatik (ZF)", "Benzin", "Arkadan İtiş", 5.8)
	insertTrim(db, 21, "M340i xDrive", 374, 500, "Otomatik (ZF)", "Benzin", "4WD", 4.4)
	insertTrim(db, 21, "320d", 190, 400, "Otomatik (ZF)", "Dizel", "Arkadan İtiş", 6.8)

	fmt.Println("\nBMW 3 Serisi F30 LCI (ID 22):")
	insertTrim(db, 22, "318i", 136, 220, "Otomatik (ZF)", "Benzin", "Arkadan İtiş", 9.1)
	insertTrim(db, 22, "320i", 184, 270, "Otomatik (ZF)", "Benzin", "Arkadan İtiş", 7.0)
	insertTrim(db, 22, "330i", 252, 350, "Otomatik (ZF)", "Benzin", "Arkadan İtiş", 5.9)
	insertTrim(db, 22, "340i", 326, 450, "Otomatik (ZF)", "Benzin", "Arkadan İtiş", 5.1)
	insertTrim(db, 22, "320d", 190, 400, "Otomatik (ZF)", "Dizel", "Arkadan İtiş", 7.0)

	fmt.Println("\nBMW 3 Serisi F30 (ID 23):")
	insertTrim(db, 23, "316i", 136, 220, "Otomatik (ZF)", "Benzin", "Arkadan İtiş", 9.0)
	insertTrim(db, 23, "320i", 184, 270, "Otomatik (ZF)", "Benzin", "Arkadan İtiş", 7.0)
	insertTrim(db, 23, "328i", 245, 350, "Otomatik (ZF)", "Benzin", "Arkadan İtiş", 5.9)
	insertTrim(db, 23, "335i", 306, 400, "Otomatik (ZF)", "Benzin", "Arkadan İtiş", 5.3)
	insertTrim(db, 23, "320d", 184, 380, "Otomatik (ZF)", "Dizel", "Arkadan İtiş", 7.2)

	fmt.Println("\nBMW 3 Serisi E90 LCI (ID 24):")
	insertTrim(db, 24, "316i", 122, 160, "Otomatik", "Benzin", "Arkadan İtiş", 10.8)
	insertTrim(db, 24, "320i", 170, 210, "Otomatik", "Benzin", "Arkadan İtiş", 8.2)
	insertTrim(db, 24, "325i", 218, 250, "Otomatik", "Benzin", "Arkadan İtiş", 6.8)
	insertTrim(db, 24, "335i", 306, 400, "Otomatik", "Benzin", "Arkadan İtiş", 5.6)
	insertTrim(db, 24, "320d", 177, 350, "Otomatik", "Dizel", "Arkadan İtiş", 7.5)

	fmt.Println("\nBMW 3 Serisi E90 (ID 25):")
	insertTrim(db, 25, "316i", 122, 160, "Otomatik", "Benzin", "Arkadan İtiş", 10.8)
	insertTrim(db, 25, "320i", 150, 200, "Otomatik", "Benzin", "Arkadan İtiş", 9.2)
	insertTrim(db, 25, "325i", 218, 250, "Otomatik", "Benzin", "Arkadan İtiş", 6.9)
	insertTrim(db, 25, "330i", 258, 300, "Otomatik", "Benzin", "Arkadan İtiş", 6.3)
	insertTrim(db, 25, "320d", 163, 340, "Otomatik", "Dizel", "Arkadan İtiş", 8.0)
}

func addBMW5SeriesTrims(db *sql.DB) {
	fmt.Println("\nBMW 5 Serisi G60 (ID 26):")
	insertTrim(db, 26, "520i", 208, 330, "Otomatik (ZF)", "Benzin", "Arkadan İtiş", 7.3)
	insertTrim(db, 26, "530i", 272, 400, "Otomatik (ZF)", "Benzin", "Arkadan İtiş", 6.1)
	insertTrim(db, 26, "540i xDrive", 374, 540, "Otomatik (ZF)", "Benzin", "4WD", 4.6)
	insertTrim(db, 26, "M550e xDrive", 489, 700, "Otomatik (ZF)", "Plug-in Hybrid", "4WD", 4.2)
	insertTrim(db, 26, "520d", 197, 400, "Otomatik (ZF)", "Dizel", "Arkadan İtiş", 7.4)

	fmt.Println("\nBMW 5 Serisi G30 LCI (ID 27):")
	insertTrim(db, 27, "520i", 184, 290, "Otomatik (ZF)", "Benzin", "Arkadan İtiş", 7.9)
	insertTrim(db, 27, "530i", 252, 350, "Otomatik (ZF)", "Benzin", "Arkadan İtiş", 6.1)
	insertTrim(db, 27, "540i xDrive", 333, 450, "Otomatik (ZF)", "Benzin", "4WD", 4.6)
	insertTrim(db, 27, "M550i xDrive", 530, 750, "Otomatik (ZF)", "Benzin", "4WD", 3.8)
	insertTrim(db, 27, "520d", 190, 400, "Otomatik (ZF)", "Dizel", "Arkadan İtiş", 7.2)

	fmt.Println("\nBMW 5 Serisi G30 (ID 28):")
	insertTrim(db, 28, "520i", 184, 290, "Otomatik (ZF)", "Benzin", "Arkadan İtiş", 7.8)
	insertTrim(db, 28, "530i", 252, 350, "Otomatik (ZF)", "Benzin", "Arkadan İtiş", 6.2)
	insertTrim(db, 28, "540i xDrive", 340, 450, "Otomatik (ZF)", "Benzin", "4WD", 4.8)
	insertTrim(db, 28, "M550i xDrive", 462, 650, "Otomatik (ZF)", "Benzin", "4WD", 4.0)
	insertTrim(db, 28, "520d", 190, 400, "Otomatik (ZF)", "Dizel", "Arkadan İtiş", 7.5)

	fmt.Println("\nBMW 5 Serisi F10 LCI (ID 29):")
	insertTrim(db, 29, "520i", 184, 270, "Otomatik (ZF)", "Benzin", "Arkadan İtiş", 7.9)
	insertTrim(db, 29, "528i", 245, 350, "Otomatik (ZF)", "Benzin", "Arkadan İtiş", 6.2)
	insertTrim(db, 29, "535i", 306, 400, "Otomatik (ZF)", "Benzin", "Arkadan İtiş", 5.7)
	insertTrim(db, 29, "550i xDrive", 450, 650, "Otomatik (ZF)", "Benzin", "4WD", 4.3)
	insertTrim(db, 29, "520d", 190, 400, "Otomatik (ZF)", "Dizel", "Arkadan İtiş", 7.4)

	fmt.Println("\nBMW 5 Serisi F10 (ID 30):")
	insertTrim(db, 30, "520i", 184, 270, "Otomatik (ZF)", "Benzin", "Arkadan İtiş", 8.1)
	insertTrim(db, 30, "528i", 258, 310, "Otomatik (ZF)", "Benzin", "Arkadan İtiş", 6.2)
	insertTrim(db, 30, "535i", 306, 400, "Otomatik (ZF)", "Benzin", "Arkadan İtiş", 5.7)
	insertTrim(db, 30, "550i", 407, 600, "Otomatik (ZF)", "Benzin", "Arkadan İtiş", 5.0)
	insertTrim(db, 30, "520d", 184, 380, "Otomatik (ZF)", "Dizel", "Arkadan İtiş", 7.9)

	fmt.Println("\nBMW 5 Serisi E60 LCI (ID 31):")
	insertTrim(db, 31, "520i", 170, 210, "Otomatik", "Benzin", "Arkadan İtiş", 9.0)
	insertTrim(db, 31, "525i", 218, 250, "Otomatik", "Benzin", "Arkadan İtiş", 7.5)
	insertTrim(db, 31, "530i", 272, 320, "Otomatik", "Benzin", "Arkadan İtiş", 6.4)
	insertTrim(db, 31, "550i", 367, 490, "Otomatik", "Benzin", "Arkadan İtiş", 5.4)
	insertTrim(db, 31, "520d", 177, 350, "Otomatik", "Dizel", "Arkadan İtiş", 8.3)

	fmt.Println("\nBMW 5 Serisi E60 (ID 32):")
	insertTrim(db, 32, "520i", 170, 210, "Otomatik", "Benzin", "Arkadan İtiş", 9.2)
	insertTrim(db, 32, "525i", 192, 245, "Otomatik", "Benzin", "Arkadan İtiş", 7.8)
	insertTrim(db, 32, "530i", 258, 300, "Otomatik", "Benzin", "Arkadan İtiş", 6.6)
	insertTrim(db, 32, "545i", 333, 450, "Otomatik", "Benzin", "Arkadan İtiş", 5.4)
	insertTrim(db, 32, "520d", 163, 340, "Otomatik", "Dizel", "Arkadan İtiş", 8.7)
}

func addPassatB7B6Trims(db *sql.DB) {
	fmt.Println("\nVW Passat B7 (ID 12):")
	insertTrim(db, 12, "1.4 TSI EcoFuel", 150, 220, "Otomatik (DSG)", "Benzin", "Önden Çekiş", 9.0)
	insertTrim(db, 12, "1.8 TSI Comfortline", 160, 250, "Otomatik (DSG)", "Benzin", "Önden Çekiş", 8.4)
	insertTrim(db, 12, "2.0 TSI Highline", 210, 280, "Otomatik (DSG)", "Benzin", "Önden Çekiş", 7.2)
	insertTrim(db, 12, "1.6 TDI BlueMotion", 105, 250, "Otomatik (DSG)", "Dizel", "Önden Çekiş", 11.4)
	insertTrim(db, 12, "2.0 TDI Highline", 170, 350, "Otomatik (DSG)", "Dizel", "Önden Çekiş", 8.4)

	fmt.Println("\nVW Passat B6 (ID 13):")
	insertTrim(db, 13, "1.6 FSI Comfortline", 115, 155, "Otomatik (Tiptronic)", "Benzin", "Önden Çekiş", 11.8)
	insertTrim(db, 13, "2.0 FSI Highline", 150, 200, "Otomatik (Tiptronic)", "Benzin", "Önden Çekiş", 9.6)
	insertTrim(db, 13, "2.0 TSI Highline", 200, 280, "Otomatik (DSG)", "Benzin", "Önden Çekiş", 7.5)
	insertTrim(db, 13, "3.6 V6 4Motion", 280, 350, "Otomatik (DSG)", "Benzin", "4WD", 6.1)
	insertTrim(db, 13, "1.9 TDI Comfortline", 105, 250, "Otomatik (DSG)", "Dizel", "Önden Çekiş", 11.9)
	insertTrim(db, 13, "2.0 TDI Highline", 140, 320, "Otomatik (DSG)", "Dizel", "Önden Çekiş", 9.3)

	fmt.Println("\nVW Passat B5.5 (ID 14):")
	insertTrim(db, 14, "1.6 Comfortline", 102, 148, "Otomatik (Tiptronic)", "Benzin", "Önden Çekiş", 13.0)
	insertTrim(db, 14, "1.8T Highline", 150, 210, "Otomatik (Tiptronic)", "Benzin", "Önden Çekiş", 9.5)
	insertTrim(db, 14, "2.3 V5 Highline", 170, 220, "Otomatik (Tiptronic)", "Benzin", "Önden Çekiş", 9.0)
	insertTrim(db, 14, "2.8 V6 4Motion", 193, 280, "Otomatik (Tiptronic)", "Benzin", "4WD", 8.0)
	insertTrim(db, 14, "1.9 TDI Comfortline", 130, 310, "Otomatik (Tiptronic)", "Dizel", "Önden Çekiş", 10.2)
}

func addAudi8VTrims(db *sql.DB) {
	fmt.Println("\nAudi A3 8V1 (ID 2):")
	insertTrim(db, 2, "1.2 TFSI", 105, 175, "S tronic", "Benzin", "Önden Çekiş", 10.5)
	insertTrim(db, 2, "1.4 TFSI", 122, 200, "S tronic", "Benzin", "Önden Çekiş", 9.3)
	insertTrim(db, 2, "1.8 TFSI", 180, 250, "S tronic", "Benzin", "Önden Çekiş", 7.3)
	insertTrim(db, 2, "2.0 TFSI Quattro", 220, 350, "S tronic", "Benzin", "4WD", 6.2)
	insertTrim(db, 2, "S3 2.0 TFSI Quattro", 300, 380, "S tronic", "Benzin", "4WD", 5.0)
	insertTrim(db, 2, "1.6 TDI", 105, 250, "S tronic", "Dizel", "Önden Çekiş", 10.5)
	insertTrim(db, 2, "2.0 TDI", 150, 320, "S tronic", "Dizel", "Önden Çekiş", 8.2)

	fmt.Println("\nAudi A3 8V2 (ID 33):")
	insertTrim(db, 33, "1.0 TFSI", 115, 200, "S tronic", "Benzin", "Önden Çekiş", 9.9)
	insertTrim(db, 33, "1.4 TFSI CoD", 150, 250, "S tronic", "Benzin", "Önden Çekiş", 8.2)
	insertTrim(db, 33, "1.5 TFSI CoD", 150, 250, "S tronic", "Benzin", "Önden Çekiş", 8.2)
	insertTrim(db, 33, "2.0 TFSI Quattro", 190, 320, "S tronic", "Benzin", "4WD", 6.7)
	insertTrim(db, 33, "S3 2.0 TFSI Quattro", 310, 400, "S tronic", "Benzin", "4WD", 4.8)
	insertTrim(db, 33, "RS3 2.5 TFSI Quattro", 400, 480, "S tronic", "Benzin", "4WD", 4.1)
	insertTrim(db, 33, "1.6 TDI", 116, 250, "S tronic", "Dizel", "Önden Çekiş", 9.7)
	insertTrim(db, 33, "2.0 TDI", 150, 340, "S tronic", "Dizel", "Önden Çekiş", 8.1)
}
