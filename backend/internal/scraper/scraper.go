package scraper

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const baseURL = "https://www.cars-data.com"

// Map brands to their URL partials
var brands = map[string]string{
	"Fiat": "fiat",
	"BMW":  "bmw",
	"Ford": "ford",
}

// Helper to make requests with User-Agent
func makeRequest(url string) (*http.Response, error) {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	// Use a common browser User-Agent
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	return client.Do(req)
}

func ScrapeAll(db *sql.DB) error {
	for brand, slug := range brands {
		log.Printf("Scraping brand: %s...", brand)
		if err := scrapeBrand(db, brand, slug); err != nil {
			log.Printf("Error scraping %s: %v", brand, err)
		}
	}
	log.Println("Scraping finished.")
	return nil
}

func scrapeBrand(db *sql.DB, brandName, slug string) error {
	url := fmt.Sprintf("%s/en/%s", baseURL, slug)
	resp, err := makeRequest(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return err
	}

	// Find Models
	// cars-data lists models in section.models .col-4 or similar.
	// We'll look for links inside the main content area that look like /en/brand/model
	var modelLinks []string
	seenLinks := make(map[string]bool)

	// Selector based on observation: links often contain /en/{brand}/
	pathPrefix := fmt.Sprintf("/en/%s/", slug)

	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists && strings.Contains(href, pathPrefix) && !seenLinks[href] {
			// Filter out junk
			if !strings.Contains(href, "#") &&
				!strings.Contains(href, "fastest") &&
				!strings.Contains(href, "sport") &&
				!strings.Contains(href, "consumption") {
				modelLinks = append(modelLinks, href)
				seenLinks[href] = true
			}
		}
	})

	// Limit to top 5 for demo speed
	if len(modelLinks) > 5 {
		modelLinks = modelLinks[:5]
	}

	for _, link := range modelLinks {
		// Ensure absolute URL
		if !strings.HasPrefix(link, "http") {
			link = baseURL + link
		}
		if err := scrapeModel(db, brandName, link); err != nil {
			log.Printf("Error scraping model %s: %v", link, err)
		}
		time.Sleep(1 * time.Second)
	}
	return nil
}

func scrapeModel(db *sql.DB, brandName, modelURL string) error {
	resp, err := makeRequest(modelURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return err
	}

	// Model Name from h1
	rawModelName := doc.Find("h1").Text()
	modelName := strings.ReplaceAll(rawModelName, " models", "")
	modelName = strings.TrimSpace(modelName)

	// Find Image
	imageURL := ""
	doc.Find("img").Each(func(i int, s *goquery.Selection) {
		src, exists := s.Attr("src")
		if exists && imageURL == "" {
			// Filter out junk images
			if strings.Contains(src, "/design/") || strings.Contains(src, "logo") || strings.Contains(src, "icon") || strings.Contains(src, "search") {
				return
			}
			// Must contain cars-data
			if strings.Contains(src, "cars-data") || strings.HasPrefix(src, "http") {
				imageURL = src
			} else if strings.HasPrefix(src, "/") {
				// Relative URL
				imageURL = baseURL + src
			}
		}
	})

	// Find Generations/Years
	// Links look like /en/{brand}-{model}-{year}/{id}
	// We will look for these links in the content
	var genLinks []string
	seenGens := make(map[string]bool)

	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		// heuristic: /en/{slug}-
		// e.g. /en/fiat-500e-2020/
		if exists && strings.Contains(href, "/en/"+strings.ToLower(brandName)+"-") && !seenGens[href] {
			genLinks = append(genLinks, href)
			seenGens[href] = true
		}
	})

	// Limit 1 generation
	if len(genLinks) > 1 {
		genLinks = genLinks[:1]
	}

	for _, link := range genLinks {
		if !strings.HasPrefix(link, "http") {
			link = baseURL + link
		}

		// Insert Vehicle (using Generation as a holder for now)
		vehicleID, err := insertVehicle(db, brandName, modelName, "Generation "+link, imageURL, link)
		if err != nil {
			return err
		}

		if err := scrapeGeneration(db, vehicleID, link); err != nil {
			log.Printf("Error scraping generation %s: %v", link, err)
		}
		time.Sleep(1 * time.Second)
	}
	return nil
}

func scrapeGeneration(db *sql.DB, vehicleID int64, genURL string) error {
	resp, err := makeRequest(genURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return err
	}

	// Find Trims
	// Links to "...-specs/ID"
	var trimLinks []string
	seenTrims := make(map[string]bool)

	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists && strings.Contains(href, "-specs/") && !seenTrims[href] {
			trimLinks = append(trimLinks, href)
			seenTrims[href] = true
		}
	})

	// Limit 3 trims
	if len(trimLinks) > 3 {
		trimLinks = trimLinks[:3]
	}

	for _, link := range trimLinks {
		if !strings.HasPrefix(link, "http") {
			link = baseURL + link
		}

		// Use the link text as trim name if available?
		// We'll scrape the trim page for the real name
		trimName := "Trim " + link
		trimID, err := insertTrim(db, vehicleID, trimName, link)
		if err != nil {
			return err
		}

		// Scrape Tech Specs
		techURL := link + "/tech"
		if err := scrapeTrimSpecs(db, trimID, techURL); err != nil {
			log.Printf("Error scraping trim specs %s: %v", techURL, err)
		}
		time.Sleep(1 * time.Second)
	}
	return nil
}

func scrapeTrimSpecs(db *sql.DB, trimID int64, techURL string) error {
	resp, err := makeRequest(techURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return err
	}

	// Update trim name
	h1 := doc.Find("h1").Text()
	h1 = strings.ReplaceAll(h1, " technical specs", "")
	_, err = db.Exec("UPDATE trims SET name = ? WHERE id = ?", h1, trimID)
	if err != nil {
		log.Printf("Failed to update trim name: %v", err)
	}

	// Parse Specs
	// cars-data often uses tables or DLs. Let's try finding dl dt dd
	// Structure seems to be sections with H2, then DL or Table.

	// Try finding all tr
	doc.Find("tr").Each(func(i int, s *goquery.Selection) {
		// key is first td, val is second td? or th/td
		var key, val string
		s.Find("td").Each(func(j int, td *goquery.Selection) {
			if j == 0 {
				key = strings.TrimSpace(td.Text())
			} else if j == 1 {
				val = strings.TrimSpace(td.Text())
			}
		})

		if key != "" && val != "" {
			// Category? Can't easy infer without efficient traversal
			insertSpec(db, trimID, "Tech", key, val)
		}
	})

	return nil
}

// -- Database Helpers --

func insertVehicle(db *sql.DB, brand, model, generation, imageURL, link string) (int64, error) {
	res, err := db.Exec("INSERT INTO vehicles (brand, model, generation, image_url, link) VALUES (?, ?, ?, ?, ?)", brand, model, generation, imageURL, link)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func insertTrim(db *sql.DB, vehicleID int64, name, link string) (int64, error) {
	res, err := db.Exec("INSERT INTO trims (vehicle_id, name, link) VALUES (?, ?, ?)", vehicleID, name, link)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func insertSpec(db *sql.DB, trimID int64, category, name, value string) {
	_, err := db.Exec("INSERT INTO specs (trim_id, category, name, value) VALUES (?, ?, ?, ?)", trimID, category, name, value)
	if err != nil {
		log.Printf("Error inserting spec: %v", err)
	}
}
