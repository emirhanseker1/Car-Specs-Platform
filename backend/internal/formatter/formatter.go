package formatter

import (
	"fmt"
	"strings"

	"github.com/emirh/car-specs/backend/internal/models"
)

// FormatTransmission converts raw transmission codes to Turkish
func FormatTransmission(raw string) string {
	if raw == "" {
		return ""
	}

	lower := strings.ToLower(strings.TrimSpace(raw))

	switch {
	case strings.Contains(lower, "auto") || lower == "a":
		return "Otomatik"
	case strings.Contains(lower, "manual") || lower == "m":
		return "Manuel"
	case strings.Contains(lower, "cvt") || lower == "cv":
		return "CVT"
	case strings.Contains(lower, "dct") || strings.Contains(lower, "dual"):
		return "Çift Kavramalı"
	default:
		return TitleCase(raw)
	}
}

// FormatFuelType converts raw fuel type codes to Turkish
func FormatFuelType(raw string) string {
	if raw == "" {
		return ""
	}

	lower := strings.ToLower(strings.TrimSpace(raw))

	switch {
	case lower == "gas" || lower == "gasoline" || lower == "petrol":
		return "Benzin"
	case lower == "diesel" || lower == "dizel":
		return "Dizel"
	case lower == "electric" || lower == "electricity" || lower == "elektrik":
		return "Elektrik"
	case lower == "hybrid":
		return "Hibrit"
	case lower == "plug-in hybrid" || lower == "phev":
		return "Plug-in Hibrit"
	default:
		return TitleCase(raw)
	}
}

// FormatDrivetrain converts raw drivetrain codes to Turkish
func FormatDrivetrain(raw string) string {
	if raw == "" {
		return ""
	}

	lower := strings.ToLower(strings.TrimSpace(raw))

	switch lower {
	case "fwd":
		return "Önden Çekiş"
	case "rwd":
		return "Arkadan İtiş"
	case "awd", "4wd":
		return "Dört Tekerlekten Çekiş"
	case "4x4":
		return "4x4"
	default:
		return strings.ToUpper(raw)
	}
}

// TitleCase converts a string to title case
func TitleCase(s string) string {
	if s == "" {
		return ""
	}

	// Handle special cases
	upper := strings.ToUpper(s)
	switch upper {
	case "BMW", "AUDI", "VW", "SUV", "MPV", "AWD", "FWD", "RWD", "CVT", "DCT", "TFSI", "TDI", "S TRONIC", "S-TRONIC":
		return upper
	}

	words := strings.Fields(s)
	for i, word := range words {
		if len(word) > 0 {
			upperWord := strings.ToUpper(word)
			// Keep acronyms uppercase
			if (len(word) <= 3 && upperWord == word) || upperWord == "TFSI" || upperWord == "TDI" || upperWord == "DSG" {
				words[i] = upperWord
			} else {
				words[i] = strings.ToUpper(string(word[0])) + strings.ToLower(word[1:])
			}
		}
	}
	return strings.Join(words, " ")
}

// FormatBrandName formats brand names with proper capitalization
func FormatBrandName(raw string) string {
	if raw == "" {
		return ""
	}

	lower := strings.ToLower(strings.TrimSpace(raw))

	// Special brand name formatting
	switch lower {
	case "bmw":
		return "BMW"
	case "audi":
		return "Audi"
	case "volkswagen", "vw":
		return "Volkswagen"
	case "mercedes-benz", "mercedes":
		return "Mercedes-Benz"
	case "fiat":
		return "Fiat"
	case "toyota":
		return "Toyota"
	case "ford":
		return "Ford"
	case "renault":
		return "Renault"
	default:
		return TitleCase(raw)
	}
}

// FormatModelName formats model names with proper capitalization
func FormatModelName(raw string) string {
	if raw == "" {
		return ""
	}

	// Keep certain patterns uppercase
	raw = strings.ReplaceAll(raw, "xdrive", "xDrive")
	raw = strings.ReplaceAll(raw, "quattro", "Quattro")

	return TitleCase(raw)
}

// FormatPowerWithUnit formats power value with HP unit
func FormatPowerWithUnit(hp *int) string {
	if hp == nil || *hp == 0 {
		return ""
	}
	return fmt.Sprintf("%d HP", *hp)
}

// FormatTorqueWithUnit formats torque value with Nm unit
func FormatTorqueWithUnit(nm *int) string {
	if nm == nil || *nm == 0 {
		return ""
	}
	return fmt.Sprintf("%d Nm", *nm)
}

// FormatAccelerationWithUnit formats 0-100 km/h time with seconds unit
func FormatAccelerationWithUnit(sec *float64) string {
	if sec == nil || *sec == 0 {
		return ""
	}
	return fmt.Sprintf("%.1f s", *sec)
}

// FormatSpeedWithUnit formats top speed with km/h unit
func FormatSpeedWithUnit(kmh *int) string {
	if kmh == nil || *kmh == 0 {
		return ""
	}
	return fmt.Sprintf("%d km/h", *kmh)
}

// FormatTrim formats a complete trim object with all human-readable values
func FormatTrim(trim *models.Trim) {
	if trim == nil {
		return
	}

	// Format transmission
	if trim.TransmissionType != nil {
		formatted := FormatTransmission(*trim.TransmissionType)
		trim.TransmissionType = &formatted
	}

	// Format fuel type
	if trim.FuelType != nil {
		formatted := FormatFuelType(*trim.FuelType)
		trim.FuelType = &formatted
	}

	// Format drivetrain
	if trim.Drivetrain != nil {
		formatted := FormatDrivetrain(*trim.Drivetrain)
		trim.Drivetrain = &formatted
	}

	// Format trim name
	trim.Name = TitleCase(trim.Name)

	// Format related model and brand if present
	if trim.Model != nil {
		trim.Model.Name = FormatModelName(trim.Model.Name)

		if trim.Model.Brand != nil {
			trim.Model.Brand.Name = FormatBrandName(trim.Model.Brand.Name)
		}
	}
}

// FormatTrims formats a slice of trims
func FormatTrims(trims []*models.Trim) {
	for _, trim := range trims {
		FormatTrim(trim)
	}
}
