package importer

import (
	"log"
	"strconv"
	"time"

	"github.com/emirh/car-specs/backend/internal/service"
	"github.com/emirh/car-specs/backend/pkg/carquery"
	"gorm.io/gorm"
)

// SyncCarQueryData fetches data from CarQuery API for the given year range and saves it to the DB.
func SyncCarQueryData(db *gorm.DB, client *carquery.Client, startYear, endYear int) error {
	log.Printf("Starting CarQuery Sync for years %d-%d...", startYear, endYear)

	for year := startYear; year <= endYear; year++ {
		log.Printf("Fetching makes for %d...", year)
		makes, err := client.GetMakes(year)
		if err != nil {
			log.Printf("Error fetching makes for year %d: %v. Skipping year.", year, err)
			continue
		}

		for _, mk := range makes {
			// Optimization: Skip fetching models if we want to save API calls,
			// but we need models to get trims effectively or iterating makes is safest.
			// CarQuery GetTrims can filter by Make+Year directly?
			// Let's check existing client usage.
			// Client has GetTrims(filter).
			// If we just request GetTrims(Make: mk, Year: year), we get all trims for that make in that year.
			// This might be more efficient than getting models first.
			// But GetTrims might have a limit on results? CarQuery usually returns JSON.

			// Let's try fetching Trims by Make and Year directly to batch it.
			log.Printf("Fetching trims for %s (%d)...", mk.MakeDisplay, year)

			trims, err := client.GetTrims(carquery.TrimFilter{
				Make: mk.MakeID, // use ID for filtering
				Year: year,
			})
			if err != nil {
				log.Printf("Error fetching trims for %s %d: %v", mk.MakeDisplay, year, err)
				continue
			}

			if len(trims) == 0 {
				continue
			}

			// Map to CarJSON DTO
			var batch []service.CarJSON
			for _, t := range trims {
				// Parse year safely
				y, _ := strconv.Atoi(t.ModelYear)
				if y == 0 {
					y = year
				}

				// Construct Trim Name (e.g. "Series 3 320i")
				// ModelTrim often contains the specific version
				trimName := t.ModelTrim
				if trimName == "" {
					trimName = t.ModelName // Fallback
				}

				batch = append(batch, service.CarJSON{
					Make:  t.ModelMakeDisplay,
					Model: t.ModelName,
					Trim:  trimName,
					Year:  y,
				})
			}

			// Save to DB via Service
			if err := service.ImportCarData(db, batch); err != nil {
				log.Printf("Error importing batch for %s %d: %v", mk.MakeDisplay, year, err)
			}

			// Respect API rate limits
			time.Sleep(200 * time.Millisecond)
		}
	}
	log.Println("CarQuery Sync Completed.")
	return nil
}
