package importer

import (
	"bytes"
	"database/sql"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"net/url"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type importRow struct {
	Brand            string
	Model            string
	Generation       string
	StartYear        *int64
	EndYear          *int64
	IsFacelift       *bool
	Market           string
	TrimName         string
	EngineCode       string
	FuelType         string
	DisplacementCC   *int64
	PowerHP          *int64
	TorqueNM         *int64
	TransmissionType string
	Gears            *int64
	Drive            string
	MarketScope      string
	SpecCategory     string
	SpecName         string
	SpecValue        string
	SourceURL        string
	SourceTitle      string
	SourceType       string
	SourceMarket     string
	SourceRetrieved  string
	SourcePage       string
	SourceNote       string
}

func ImportCSV(db *sql.DB, csvPath string) error {
	f, err := os.Open(csvPath)
	if err != nil {
		return err
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.TrimLeadingSpace = true

	headers, err := r.Read()
	if err != nil {
		return err
	}

	headerIndex := map[string]int{}
	for i, h := range headers {
		headerIndex[strings.ToLower(strings.TrimSpace(h))] = i
	}

	// Auto-detect dataset formats
	if _, ok := headerIndex["model_year"]; ok {
		return importVWGolfEpeyCSV(db, r, headerIndex)
	}

	get := func(cols []string, key string) string {
		i, ok := headerIndex[key]
		if !ok || i < 0 || i >= len(cols) {
			return ""
		}
		return strings.TrimSpace(cols[i])
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	for {
		cols, readErr := r.Read()
		if readErr != nil {
			if errors.Is(readErr, io.EOF) {
				break
			}
			err = readErr
			return err
		}

		row, parseErr := parseRow(cols, get)
		if parseErr != nil {
			err = parseErr
			return err
		}
		if row.Brand == "" || row.Model == "" || row.Generation == "" {
			continue
		}

		vehicleID, vErr := upsertVehicle(tx, row.Brand, row.Model, row.Generation)
		if vErr != nil {
			err = vErr
			return err
		}

		if row.StartYear != nil || row.EndYear != nil || row.IsFacelift != nil || row.Market != "" {
			if gmErr := upsertGenerationMeta(tx, vehicleID, row); gmErr != nil {
				err = gmErr
				return err
			}
		}

		if row.TrimName != "" {
			trimID, tErr := upsertTrim(tx, vehicleID, row.TrimName)
			if tErr != nil {
				err = tErr
				return err
			}

			if row.EngineCode != "" || row.FuelType != "" || row.DisplacementCC != nil || row.PowerHP != nil || row.TorqueNM != nil || row.TransmissionType != "" || row.Gears != nil || row.Drive != "" || row.MarketScope != "" {
				if pmErr := upsertPowertrainMeta(tx, trimID, row); pmErr != nil {
					err = pmErr
					return err
				}
			}

			if row.SpecCategory != "" && row.SpecName != "" && row.SpecValue != "" {
				specID, sErr := insertSpec(tx, trimID, row.SpecCategory, row.SpecName, row.SpecValue)
				if sErr != nil {
					err = sErr
					return err
				}

				if row.SourceURL != "" {
					sourceID, sdErr := upsertSourceDocument(tx, row)
					if sdErr != nil {
						err = sdErr
						return err
					}
					if ssErr := upsertSpecSource(tx, specID, sourceID, row.SourcePage, row.SourceNote); ssErr != nil {
						err = ssErr
						return err
					}
				}
			}
		}
	}

	if cErr := tx.Commit(); cErr != nil {
		return cErr
	}
	return nil
}

func importVWGolfEpeyCSV(db *sql.DB, r *csv.Reader, headerIndex map[string]int) error {
	get := func(cols []string, key string) string {
		i, ok := headerIndex[key]
		if !ok || i < 0 || i >= len(cols) {
			return ""
		}
		return strings.TrimSpace(cols[i])
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	brand := "Volkswagen"
	model := "Golf"
	defaultSourceType := "trusted"
	defaultMarketScope := "TR"
	imageCache := make(map[string]string)

	for {
		cols, readErr := r.Read()
		if readErr != nil {
			if errors.Is(readErr, io.EOF) {
				break
			}
			err = readErr
			return err
		}

		yearStr := get(cols, "model_year")
		if yearStr == "" {
			continue
		}
		year, pErr := strconv.ParseInt(yearStr, 10, 64)
		if pErr != nil {
			continue
		}

		generation := yearStr
		vehicleID, vErr := upsertVehicle(tx, brand, model, generation)
		if vErr != nil {
			err = vErr
			return err
		}

		// Set generation meta as single-year range for now.
		startYear := year
		endYear := year
		gm := importRow{StartYear: &startYear, EndYear: &endYear, Market: defaultMarketScope}
		if gmErr := upsertGenerationMeta(tx, vehicleID, gm); gmErr != nil {
			err = gmErr
			return err
		}

		fullModelName := get(cols, "model_name")
		imageURL := ""
		modelNameCitation := get(cols, "model_name_citation")
		if modelNameCitation != "" {
			if cached, ok := imageCache[modelNameCitation]; ok {
				imageURL = cached
			} else {
				img, _ := fetchOpenGraphImage(modelNameCitation)
				imageCache[modelNameCitation] = img
				imageURL = img
			}
		}
		if imageURL != "" {
			if uErr := updateVehicleImage(tx, vehicleID, imageURL); uErr != nil {
				err = uErr
				return err
			}
		}
		trimName := normalizeVWTrimName(yearStr, fullModelName)
		if trimName == "" {
			trimName = fullModelName
		}
		if trimName == "" {
			trimName = "Variant " + yearStr
		}

		trimID, tErr := upsertTrim(tx, vehicleID, trimName)
		if tErr != nil {
			err = tErr
			return err
		}

		// Powertrain meta (optional)
		disp, _ := parseInt64Ptr(get(cols, "engine_capacity_cc"))
		power, _ := parseInt64Ptr(get(cols, "engine_power_hp"))
		fuelType := get(cols, "fuel_type")
		transmission := get(cols, "epey_com_transmission_type")

		pm := importRow{
			FuelType:         fuelType,
			DisplacementCC:   disp,
			PowerHP:          power,
			TransmissionType: transmission,
			MarketScope:      defaultMarketScope,
		}
		if pmErr := upsertPowertrainMeta(tx, trimID, pm); pmErr != nil {
			err = pmErr
			return err
		}

		// Helper: insert spec + link citation as a source
		addSpec := func(category, name, value, citationURL string) error {
			if value == "" {
				return nil
			}
			specID, sErr := insertSpec(tx, trimID, category, name, value)
			if sErr != nil {
				return sErr
			}
			if citationURL == "" {
				return nil
			}

			sourceRow := importRow{
				SourceURL:       citationURL,
				SourceTitle:     fullModelName,
				SourceType:      defaultSourceType,
				SourceMarket:    defaultMarketScope,
				SourceRetrieved: time.Now().Format(time.RFC3339),
			}
			sourceID, sdErr := upsertSourceDocument(tx, sourceRow)
			if sdErr != nil {
				return sdErr
			}
			if ssErr := upsertSpecSource(tx, specID, sourceID, "", name); ssErr != nil {
				return ssErr
			}
			return nil
		}

		if err := addSpec("General", "Model year", yearStr, get(cols, "model_year_citation")); err != nil {
			return err
		}
		if err := addSpec("General", "Body type", get(cols, "epey_com_body_type"), get(cols, "epey_com_body_type_citation")); err != nil {
			return err
		}
		if err := addSpec("Transmission", "Transmission", transmission, get(cols, "epey_com_transmission_type_citation")); err != nil {
			return err
		}

		if disp != nil {
			if err := addSpec("Engine", "Engine capacity", fmt.Sprintf("%d cc", *disp), get(cols, "engine_capacity_cc_citation")); err != nil {
				return err
			}
		}
		if power != nil {
			if err := addSpec("Engine", "Power", fmt.Sprintf("%d hp", *power), get(cols, "engine_power_hp_citation")); err != nil {
				return err
			}
		}
		if err := addSpec("Engine", "Fuel type", fuelType, get(cols, "fuel_type_citation")); err != nil {
			return err
		}

		if err := addSpec("Consumption", "Fuel consumption (city)", appendUnit(get(cols, "epey_com_fuel_consumption_l_100km_city"), "L/100km"), get(cols, "epey_com_fuel_consumption_l_100km_city_citation")); err != nil {
			return err
		}
		if err := addSpec("Consumption", "Fuel consumption (highway)", appendUnit(get(cols, "epey_com_fuel_consumption_l_100km_highway"), "L/100km"), get(cols, "epey_com_fuel_consumption_l_100km_highway_citation")); err != nil {
			return err
		}

		if err := addSpec("Performance", "0-100 km/h", appendUnit(get(cols, "epey_com_acceleration_0_100kmh_seconds"), "s"), get(cols, "epey_com_acceleration_0_100kmh_seconds_citation")); err != nil {
			return err
		}
		if err := addSpec("Performance", "Top speed", appendUnit(get(cols, "epey_com_top_speed_kmh"), "km/h"), get(cols, "epey_com_top_speed_kmh_citation")); err != nil {
			return err
		}
	}

	if cErr := tx.Commit(); cErr != nil {
		return cErr
	}
	return nil
}

func normalizeVWTrimName(yearStr, fullName string) string {
	name := strings.TrimSpace(fullName)
	if name == "" {
		return ""
	}
	// Remove leading year if present
	if strings.HasPrefix(name, yearStr) {
		name = strings.TrimSpace(strings.TrimPrefix(name, yearStr))
	}
	// Remove leading brand/model if present
	for _, p := range []string{"Volkswagen", "VW"} {
		if strings.HasPrefix(strings.ToLower(name), strings.ToLower(p)) {
			name = strings.TrimSpace(name[len(p):])
		}
	}
	if strings.HasPrefix(strings.ToLower(name), "golf") {
		name = strings.TrimSpace(name[len("golf"):])
	}
	return strings.TrimSpace(name)
}

func appendUnit(v string, unit string) string {
	v = strings.TrimSpace(v)
	if v == "" {
		return ""
	}
	if unit == "" {
		return v
	}
	return v + " " + unit
}

func updateVehicleImage(tx *sql.Tx, vehicleID int64, imageURL string) error {
	imageURL = strings.TrimSpace(imageURL)
	if imageURL == "" {
		return nil
	}
	_, err := tx.Exec(
		"UPDATE vehicles SET image_url = CASE WHEN image_url IS NULL OR image_url = '' THEN ? ELSE image_url END WHERE id = ?",
		imageURL,
		vehicleID,
	)
	return err
}

func fetchOpenGraphImage(pageURL string) (string, error) {
	pageURL = strings.TrimSpace(pageURL)
	if pageURL == "" {
		return "", nil
	}

	client := &http.Client{Timeout: 20 * time.Second}
	req, err := http.NewRequest("GET", pageURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("unexpected status: %s", resp.Status)
	}

	// Limit read to avoid huge pages
	const maxBytes = 1024 * 1024
	buf, err := io.ReadAll(io.LimitReader(resp.Body, maxBytes))
	if err != nil {
		return "", err
	}

	img := extractOGImage(buf)
	if img == "" {
		return "", nil
	}

	// Normalize to absolute URL when possible
	parsedPage, err := url.Parse(pageURL)
	if err == nil {
		parsedImg, err := url.Parse(img)
		if err == nil {
			img = parsedPage.ResolveReference(parsedImg).String()
		}
	}

	return img, nil
}

func extractOGImage(html []byte) string {
	// Look for: <meta property="og:image" content="...">
	lower := bytes.ToLower(html)
	needle := []byte("og:image")
	idx := bytes.Index(lower, needle)
	if idx < 0 {
		return ""
	}

	// Search forward from where og:image appears for content=
	windowStart := idx
	windowEnd := idx + 2000
	if windowEnd > len(lower) {
		windowEnd = len(lower)
	}
	window := lower[windowStart:windowEnd]

	contentIdx := bytes.Index(window, []byte("content="))
	if contentIdx < 0 {
		return ""
	}

	// Find quote char after content=
	qPos := contentIdx + len("content=")
	if qPos >= len(window) {
		return ""
	}
	quote := window[qPos]
	if quote != '"' && quote != '\'' {
		return ""
	}

	start := qPos + 1
	endRel := bytes.IndexByte(window[start:], quote)
	if endRel < 0 {
		return ""
	}
	end := start + endRel

	// Return from original html slice to preserve casing
	origWindow := html[windowStart:windowEnd]
	return strings.TrimSpace(string(origWindow[start:end]))
}

func parseRow(cols []string, get func([]string, string) string) (importRow, error) {
	row := importRow{
		Brand:      get(cols, "brand"),
		Model:      get(cols, "model"),
		Generation: get(cols, "generation"),
		Market:     get(cols, "market"),
		TrimName:   get(cols, "trim_name"),

		EngineCode:       get(cols, "engine_code"),
		FuelType:         get(cols, "fuel_type"),
		TransmissionType: get(cols, "transmission_type"),
		Drive:            get(cols, "drive"),
		MarketScope:      get(cols, "market_scope"),

		SpecCategory: get(cols, "spec_category"),
		SpecName:     get(cols, "spec_name"),
		SpecValue:    get(cols, "spec_value"),

		SourceURL:       get(cols, "source_url"),
		SourceTitle:     get(cols, "source_title"),
		SourceType:      get(cols, "source_type"),
		SourceMarket:    get(cols, "source_market_scope"),
		SourceRetrieved: get(cols, "source_retrieved_at"),
		SourcePage:      get(cols, "source_page"),
		SourceNote:      get(cols, "source_note"),
	}

	startYear, err := parseInt64Ptr(get(cols, "start_year"))
	if err != nil {
		return importRow{}, fmt.Errorf("start_year: %w", err)
	}
	row.StartYear = startYear

	endYear, err := parseInt64Ptr(get(cols, "end_year"))
	if err != nil {
		return importRow{}, fmt.Errorf("end_year: %w", err)
	}
	row.EndYear = endYear

	isFacelift, err := parseBoolPtr(get(cols, "is_facelift"))
	if err != nil {
		return importRow{}, fmt.Errorf("is_facelift: %w", err)
	}
	row.IsFacelift = isFacelift

	displacement, err := parseInt64Ptr(get(cols, "displacement_cc"))
	if err != nil {
		return importRow{}, fmt.Errorf("displacement_cc: %w", err)
	}
	row.DisplacementCC = displacement

	powerHP, err := parseInt64Ptr(get(cols, "power_hp"))
	if err != nil {
		return importRow{}, fmt.Errorf("power_hp: %w", err)
	}
	row.PowerHP = powerHP

	torqueNM, err := parseInt64Ptr(get(cols, "torque_nm"))
	if err != nil {
		return importRow{}, fmt.Errorf("torque_nm: %w", err)
	}
	row.TorqueNM = torqueNM

	gears, err := parseInt64Ptr(get(cols, "gears"))
	if err != nil {
		return importRow{}, fmt.Errorf("gears: %w", err)
	}
	row.Gears = gears

	if row.SourceRetrieved == "" {
		row.SourceRetrieved = time.Now().Format(time.RFC3339)
	}

	return row, nil
}

func parseInt64Ptr(s string) (*int64, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, nil
	}
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return nil, err
	}
	return &n, nil
}

func parseBoolPtr(s string) (*bool, error) {
	s = strings.TrimSpace(strings.ToLower(s))
	if s == "" {
		return nil, nil
	}
	if s == "1" || s == "true" || s == "yes" {
		b := true
		return &b, nil
	}
	if s == "0" || s == "false" || s == "no" {
		b := false
		return &b, nil
	}
	return nil, fmt.Errorf("invalid bool: %s", s)
}

func vehicleLink(brand, model, generation string) string {
	slug := strings.ToLower(strings.TrimSpace(brand)) + "/" + strings.ToLower(strings.TrimSpace(model)) + "/" + strings.ToLower(strings.TrimSpace(generation))
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "_", "-")
	return "manual://" + slug
}

func trimLink(brand, model, generation, trimName string) string {
	slug := strings.ToLower(strings.TrimSpace(trimName))
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "_", "-")
	return vehicleLink(brand, model, generation) + "/trim/" + slug
}

func upsertVehicle(tx *sql.Tx, brand, model, generation string) (int64, error) {
	var id int64
	err := tx.QueryRow("SELECT id FROM vehicles WHERE brand = ? AND model = ? AND generation = ?", brand, model, generation).Scan(&id)
	if err == nil {
		return id, nil
	}
	if err != sql.ErrNoRows {
		return 0, err
	}

	link := vehicleLink(brand, model, generation)
	res, err := tx.Exec("INSERT INTO vehicles (brand, model, generation, image_url, link) VALUES (?, ?, ?, ?, ?)", brand, model, generation, "", link)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func upsertGenerationMeta(tx *sql.Tx, vehicleID int64, row importRow) error {
	var start interface{} = nil
	var end interface{} = nil
	var isFacelift interface{} = nil
	if row.StartYear != nil {
		start = *row.StartYear
	}
	if row.EndYear != nil {
		end = *row.EndYear
	}
	if row.IsFacelift != nil {
		if *row.IsFacelift {
			isFacelift = 1
		} else {
			isFacelift = 0
		}
	}

	_, err := tx.Exec(
		"INSERT INTO vehicle_generation_meta (vehicle_id, start_year, end_year, is_facelift, market) VALUES (?, ?, ?, COALESCE(?, 0), ?) ON CONFLICT(vehicle_id) DO UPDATE SET start_year=excluded.start_year, end_year=excluded.end_year, is_facelift=excluded.is_facelift, market=excluded.market",
		vehicleID,
		start,
		end,
		isFacelift,
		row.Market,
	)
	return err
}

func upsertTrim(tx *sql.Tx, vehicleID int64, trimName string) (int64, error) {
	var id int64
	err := tx.QueryRow("SELECT id FROM trims WHERE vehicle_id = ? AND name = ?", vehicleID, trimName).Scan(&id)
	if err == nil {
		return id, nil
	}
	if err != sql.ErrNoRows {
		return 0, err
	}

	var brand, model, generation string
	if err := tx.QueryRow("SELECT brand, model, generation FROM vehicles WHERE id = ?", vehicleID).Scan(&brand, &model, &generation); err != nil {
		return 0, err
	}

	link := trimLink(brand, model, generation, trimName)
	res, err := tx.Exec("INSERT INTO trims (vehicle_id, name, link) VALUES (?, ?, ?)", vehicleID, trimName, link)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func upsertPowertrainMeta(tx *sql.Tx, trimID int64, row importRow) error {
	var displacement interface{} = nil
	var power interface{} = nil
	var torque interface{} = nil
	var gears interface{} = nil
	if row.DisplacementCC != nil {
		displacement = *row.DisplacementCC
	}
	if row.PowerHP != nil {
		power = *row.PowerHP
	}
	if row.TorqueNM != nil {
		torque = *row.TorqueNM
	}
	if row.Gears != nil {
		gears = *row.Gears
	}

	_, err := tx.Exec(
		"INSERT INTO trim_powertrain_meta (trim_id, engine_code, fuel_type, displacement_cc, power_hp, torque_nm, transmission_type, gears, drive, market_scope) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?) ON CONFLICT(trim_id) DO UPDATE SET engine_code=excluded.engine_code, fuel_type=excluded.fuel_type, displacement_cc=excluded.displacement_cc, power_hp=excluded.power_hp, torque_nm=excluded.torque_nm, transmission_type=excluded.transmission_type, gears=excluded.gears, drive=excluded.drive, market_scope=excluded.market_scope",
		trimID,
		row.EngineCode,
		row.FuelType,
		displacement,
		power,
		torque,
		row.TransmissionType,
		gears,
		row.Drive,
		row.MarketScope,
	)
	return err
}

func insertSpec(tx *sql.Tx, trimID int64, category, name, value string) (int64, error) {
	var id int64
	err := tx.QueryRow("SELECT id FROM specs WHERE trim_id = ? AND category = ? AND name = ? AND value = ?", trimID, category, name, value).Scan(&id)
	if err == nil {
		return id, nil
	}
	if err != sql.ErrNoRows {
		return 0, err
	}

	res, err := tx.Exec("INSERT INTO specs (trim_id, category, name, value) VALUES (?, ?, ?, ?)", trimID, category, name, value)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func upsertSourceDocument(tx *sql.Tx, row importRow) (int64, error) {
	_, err := tx.Exec("INSERT OR IGNORE INTO source_documents (url, title, source_type, market_scope, retrieved_at) VALUES (?, ?, ?, ?, ?)", row.SourceURL, row.SourceTitle, row.SourceType, row.SourceMarket, row.SourceRetrieved)
	if err != nil {
		return 0, err
	}

	var id int64
	if err := tx.QueryRow("SELECT id FROM source_documents WHERE url = ?", row.SourceURL).Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func upsertSpecSource(tx *sql.Tx, specID, sourceID int64, page, note string) error {
	_, err := tx.Exec("INSERT OR IGNORE INTO spec_sources (spec_id, source_document_id, page, note) VALUES (?, ?, ?, ?)", specID, sourceID, page, note)
	return err
}
