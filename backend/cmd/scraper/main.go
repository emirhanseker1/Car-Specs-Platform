package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/emirh/car-specs/backend/internal/db"
	"github.com/gocolly/colly/v2"
)

// Global DB
var database *sql.DB

func main() {
	// 1. Initialize DB
	var err error
	database, err = db.InitDB("vehicles.db")
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	// 2. Load manual overrides
	if err := LoadOverrides("../../data/manual_overrides.json"); err != nil {
		log.Printf("Warning: Failed to load overrides: %v", err)
	}

	// 3. Ensure Audi and A3 exist
	brandID := ensureBrand("Audi")
	modelID := ensureModel(brandID, "A3")

	// 4. Setup Collector
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"),
	)
	c.SetRequestTimeout(120 * time.Second)

	// Rate Limiting
	// c.Limit(&colly.LimitRule{
	// 	DomainGlob:  "*ultimatespecs.com*",
	// 	Delay:       1 * time.Second,
	// 	RandomDelay: 1 * time.Second,
	// })

	c.OnError(func(r *colly.Response, err error) {
		log.Printf("Request URL: %s failed with error: %v", r.Request.URL, err)
	})

	// --- HANDLERS ---

	// A. Audi A3 Main Page -> Find Generations
	c.OnHTML("div.home_models_line", func(e *colly.HTMLElement) {
		genTitle := strings.TrimSpace(e.ChildText("h2"))
		if genTitle == "" {
			return
		}

		fmt.Printf("Found Generation Block: %s\n", genTitle)

		// Parse Generation Code (e.g., "Type 8Y" -> "8Y")
		genCode := parseGenerationCode(genTitle)

		// Create Generation in DB
		genID := ensureGeneration(modelID, genCode, genTitle)

		// Iterate over Body Styles
		e.ForEach("a", func(_ int, el *colly.HTMLElement) {
			link := el.Attr("href")
			bodyStyleName := strings.TrimSpace(el.ChildText("h3"))

			if link != "" && bodyStyleName != "" {
				if strings.HasPrefix(link, "/") {
					link = "https://www.ultimatespecs.com" + link
				}

				fmt.Printf("  -> Body Style: %s (%s)\n", bodyStyleName, link)

				fmt.Printf("  -> Body Style: %s (%s)\n", bodyStyleName, link)

				// URL Hack: Pass context via URL query params
				u, _ := url.Parse(link)
				q := u.Query()
				q.Set("gen_id", strconv.FormatInt(genID, 10))
				u.RawQuery = q.Encode()

				c.Visit(u.String())
			}
		})
	})

	// B. Body Style Page -> Find Trims (Engines)
	c.OnHTML("div.content-wrapper", func(e *colly.HTMLElement) {
		e.ForEach("a", func(_ int, el *colly.HTMLElement) {
			link := el.Attr("href")
			text := strings.TrimSpace(el.Text)

			// Filter irrelevant links
			if strings.Contains(link, "car-specs/Audi") && !strings.Contains(text, "Versions") && !strings.Contains(text, "View more") {
				if strings.HasPrefix(link, "/") {
					link = "https://www.ultimatespecs.com" + link
				}

				// Verify it's a trim page (usually ends in .html)
				if strings.HasSuffix(link, ".html") {
					// Check params from URL (Handler A passes it)
					qReq := e.Request.URL.Query()
					genIDStr := qReq.Get("gen_id")

					if genIDStr == "" {
						// Fallback to context
						genIDStr = e.Request.Ctx.Get("gen_id")
					}

					if genIDStr == "" {
						return
					}

					fmt.Printf("    -> Found Trim: %s\n", text)

					// URL Hack: Pass context via URL query params to ensure persistence
					u, _ := url.Parse(link)
					q := u.Query()
					q.Set("gen_id", genIDStr)
					q.Set("trim_name", text)
					q.Set("fuel_type", parseFuelType(text))
					u.RawQuery = q.Encode()

					// Use Visit instead of Request to ensure full Colly pipeline with URL params
					c.Visit(u.String())
				}
			}
		})
	})

	// C. Trim Page -> Parse Specs
	c.OnHTML("body", func(e *colly.HTMLElement) {
		// Parse from URL
		q := e.Request.URL.Query()
		trimName := q.Get("trim_name")

		if trimName == "" {
			// Try context as fallback
			trimName = e.Request.Ctx.Get("trim_name")
		}

		if trimName == "" {
			return
		}

		// Create context manually from URL params
		ctxMock := colly.NewContext()
		ctxMock.Put("gen_id", q.Get("gen_id"))
		ctxMock.Put("trim_name", trimName)
		ctxMock.Put("fuel_type", q.Get("fuel_type"))

		fmt.Printf("      -> Scraping Specs for: %s\n", trimName)

		// Parse table data
		specs := make(map[string]string)
		e.ForEach("table tr", func(_ int, el *colly.HTMLElement) {
			label := strings.TrimSpace(el.ChildText("td:nth-child(1)"))
			value := strings.TrimSpace(el.ChildText("td:nth-child(2)"))
			if label != "" {
				label = strings.TrimSuffix(label, ":")
				label = strings.TrimSpace(label)
				specs[label] = value
			}
		})

		saveTrim(ctxMock, specs)
	})

	fmt.Println("Starting Scraper for Audi A3...")
	c.Visit("https://www.ultimatespecs.com/car-specs/Audi-models/Audi-A3")

	// Save validation report
	if err := SaveReport("validation_report.json"); err != nil {
		log.Printf("Warning: Failed to save validation report: %v", err)
	}
}

// --- DB HELPERS ---

func ensureBrand(name string) int64 {
	var id int64
	err := database.QueryRow("SELECT id FROM brands WHERE name = ?", name).Scan(&id)
	if err == sql.ErrNoRows {
		res, err := database.Exec("INSERT INTO brands (name) VALUES (?)", name)
		if err != nil {
			log.Fatal(err)
		}
		id, _ = res.LastInsertId()
		fmt.Printf("Created Brand: %s (ID: %d)\n", name, id)
	}
	return id
}

func ensureModel(brandID int64, name string) int64 {
	var id int64
	err := database.QueryRow("SELECT id FROM models WHERE brand_id = ? AND name = ?", brandID, name).Scan(&id)
	if err == sql.ErrNoRows {
		res, err := database.Exec("INSERT INTO models (brand_id, name) VALUES (?, ?)", brandID, name)
		if err != nil {
			log.Fatal(err)
		}
		id, _ = res.LastInsertId()
		fmt.Printf("Created Model: %s (ID: %d)\n", name, id)
	}
	return id
}

func ensureGeneration(modelID int64, code, name string) int64 {
	var id int64
	err := database.QueryRow("SELECT id FROM generations WHERE model_id = ? AND code = ?", modelID, code).Scan(&id)

	startYear, endYear := parseYears(name)

	if err == sql.ErrNoRows {
		res, err := database.Exec(`
			INSERT INTO generations (model_id, code, name, start_year, end_year) 
			VALUES (?, ?, ?, ?, ?)`,
			modelID, code, name, startYear, endYear)
		if err != nil {
			log.Printf("Error creating generation %s: %v", name, err)
			return 0
		}
		id, _ = res.LastInsertId()
		fmt.Printf("Created Generation: %s (ID: %d)\n", code, id)
	}
	return id
}

func saveTrim(ctx *colly.Context, specs map[string]string) {
	genID, _ := strconv.ParseInt(ctx.Get("gen_id"), 10, 64)
	name := ctx.Get("trim_name")
	fuelType := ctx.Get("fuel_type")

	// Parse Specs
	// Normalized keys: "Horsepower", "Maximum torque"
	hpStr := specs["Horsepower"]
	if hpStr == "" {
		hpStr = specs["Power"]
	}

	hp := parseHP(hpStr)

	// Fetch Model ID from Generation
	var modelID int64
	err := database.QueryRow("SELECT model_id FROM generations WHERE id = ?", genID).Scan(&modelID)
	if err != nil {
		if genID == 0 {
			modelID = 0
		}
	}

	// Parse production years
	prodYears := specs["Production years"]
	if prodYears == "" {
		prodYears = specs["Years"]
	}
	startYear, endYear := parseYears(prodYears)

	// Parse Year from trim name or generation or use startYear
	nameYear := parseTrimYear(name)
	if startYear == 0 {
		startYear = nameYear
	}

	year := nameYear
	if year == 0 {
		year = startYear
	}

	// Fallback: If still no years, try to get from generation overrides
	if startYear == 0 {
		// Get generation code from database
		var genCode string
		err := database.QueryRow("SELECT code FROM generations WHERE id = ?", genID).Scan(&genCode)
		if err == nil && genCode != "" {
			// Try to get years from overrides
			overrideStart, overrideEnd, found := GetGenerationYears("Audi", "A3", genCode)
			if found {
				startYear = overrideStart
				if overrideEnd != nil {
					endYear.Valid = true
					endYear.Int64 = int64(*overrideEnd)
				}
				fmt.Printf("      ✓ Applied generation override for %s: %d-%v\n", genCode, startYear, overrideEnd)
			}
		}
	}

	// Check duplicate
	var exists int
	database.QueryRow("SELECT 1 FROM trims WHERE generation_id = ? AND name = ?", genID, name).Scan(&exists)
	if exists == 1 {
		fmt.Printf("Skipping existing trim: %s\n", name)
		return
	}

	// Validate data quality
	validationReport := ValidateTrimData(name, year, startYear, endYear, hp, specs)
	AddToReport(validationReport)

	_, err = database.Exec(`
		INSERT INTO trims (generation_id, model_id, name, year, start_year, end_year, fuel_type, power_hp, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`, genID, modelID, name, year, sql.NullInt64{Int64: int64(startYear), Valid: startYear != 0}, endYear, fuelType, hp)

	if err != nil {
		log.Printf("Error insert trim %s: %v", name, err)
	} else {
		fmt.Printf("Saved Trim: %s (HP: %d, Years: %d-%v)\n", name, hp, startYear, endYear.Int64)
	}
}

// --- PARSERS ---

func parseTrimYear(name string) int {
	re := regexp.MustCompile(`\b(19|20)\d{2}\b`)
	matches := re.FindStringSubmatch(name)
	if len(matches) > 0 {
		y, _ := strconv.Atoi(matches[0])
		return y
	}
	return 0
}

func parseGenerationCode(s string) string {
	re := regexp.MustCompile(`Type\s+(\w+)`)
	matches := re.FindStringSubmatch(s)
	if len(matches) > 1 {
		return matches[1]
	}
	return "Unknown"
}

func parseYears(s string) (int, sql.NullInt64) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, sql.NullInt64{}
	}

	// Pattern 1: (2013 - 2016) or (2013 - Present) or 2013 - 2016
	re := regexp.MustCompile(`(\d{4})\s*-\s*(\d{4}|Present)`)
	matches := re.FindStringSubmatch(s)
	if len(matches) > 1 {
		start, _ := strconv.Atoi(matches[1])
		var end sql.NullInt64
		if matches[2] != "Present" {
			e, _ := strconv.Atoi(matches[2])
			end = sql.NullInt64{Int64: int64(e), Valid: true}
		}
		return start, end
	}

	// Pattern 2: Single year like "2024" or "(2024)"
	reSingle := regexp.MustCompile(`(\d{4})`)
	match := reSingle.FindStringSubmatch(s)
	if len(match) > 1 {
		start, _ := strconv.Atoi(match[1])
		return start, sql.NullInt64{}
	}

	return 0, sql.NullInt64{}
}

func parseBodyStyle(raw string) string {
	if strings.Contains(raw, "Sportback") {
		return "Sportback"
	}
	if strings.Contains(raw, "Sedan") || strings.Contains(raw, "Limousine") {
		return "Sedan"
	}
	if strings.Contains(raw, "Cabrio") {
		return "Cabriolet"
	}
	return "Hatchback"
}

func parseFuelType(name string) string {
	if strings.Contains(name, "TFSI") {
		return "Petrol"
	}
	if strings.Contains(name, "TDI") {
		return "Diesel"
	}
	if strings.Contains(name, "e-tron") || strings.Contains(name, "TFSIe") {
		return "Hybrid"
	}
	return "Petrol"
}

func parseHP(val string) int {
	re := regexp.MustCompile(`(\d+)\s*PS`)
	m := re.FindStringSubmatch(val)
	if len(m) > 1 {
		v, _ := strconv.Atoi(m[1])
		return v
	}
	return 0
}
