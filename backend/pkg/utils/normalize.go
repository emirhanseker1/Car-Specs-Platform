package utils

import (
	"strings"
	"unicode"
)

// ToTitleCase converts a string to Title Case with special handling for automotive terms
func ToTitleCase(s string) string {
	if s == "" {
		return s
	}

	words := strings.Fields(s)
	for i, word := range words {
		// Handle special cases - automotive abbreviations
		upper := strings.ToUpper(word)
		switch upper {
		case "BMW", "GTI", "TSI", "TDI", "FSI", "TFSI", "SDI", "RS", "GT", "SE", "SL", "AMG", "SUV", "AWD", "FWD", "RWD", "HP", "PS", "KW":
			words[i] = upper
		case "E-CLASS", "C-CLASS", "S-CLASS", "A-CLASS", "B-CLASS", "G-CLASS", "CLS-CLASS", "GLE-CLASS", "GLC-CLASS":
			// Mercedes class names
			words[i] = strings.ToUpper(word[:1]) + "-" + strings.ToUpper(string(word[2])) + strings.ToLower(word[3:])
		default:
			// Standard title case
			if len(word) > 0 {
				runes := []rune(word)
				runes[0] = unicode.ToUpper(runes[0])
				for j := 1; j < len(runes); j++ {
					runes[j] = unicode.ToLower(runes[j])
				}
				words[i] = string(runes)
			}
		}
	}
	return strings.Join(words, " ")
}

// NormalizeVehicleName normalizes brand and model names to Title Case
func NormalizeVehicleName(brand, model string) (string, string) {
	return ToTitleCase(brand), ToTitleCase(model)
}
