package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

// DataQuality represents the quality assessment of a single data field
type DataQuality struct {
	Field      string      `json:"field"`
	Value      interface{} `json:"value"`
	Confidence string      `json:"confidence"` // "high", "medium", "low", "missing"
	Source     string      `json:"source"`     // "ultimatespecs", "carquery", "manual"
	Notes      string      `json:"notes"`
}

// ValidationReport represents issues found for a single trim
type ValidationReport struct {
	TrimName       string        `json:"trim_name"`
	Issues         []DataQuality `json:"issues"`
	RequiresReview bool          `json:"requires_review"`
}

// ScraperReport represents the complete validation report
type ScraperReport struct {
	Timestamp            string             `json:"timestamp"`
	TotalTrims           int                `json:"total_trims"`
	IssuesFound          int                `json:"issues_found"`
	TrimsRequiringReview []ValidationReport `json:"trims_requiring_review"`
}

var globalReport = &ScraperReport{
	Timestamp:            time.Now().Format(time.RFC3339),
	TrimsRequiringReview: []ValidationReport{},
}

// ValidateTrimData validates all fields of a trim and returns issues
func ValidateTrimData(name string, year, startYear int, endYear sql.NullInt64, hp int, specs map[string]string) ValidationReport {
	report := ValidationReport{
		TrimName: name,
		Issues:   []DataQuality{},
	}

	// Validate start_year
	if startYear == 0 {
		report.Issues = append(report.Issues, DataQuality{
			Field:      "start_year",
			Value:      0,
			Confidence: "missing",
			Source:     "ultimatespecs",
			Notes:      "Could not parse production start year from specs or name",
		})
		report.RequiresReview = true
	}

	// Validate end_year (only if production has ended)
	prodYears := specs["Production years"]
	if prodYears == "" {
		prodYears = specs["Years"]
	}
	if !endYear.Valid && prodYears != "" && !contains(prodYears, "Present") {
		report.Issues = append(report.Issues, DataQuality{
			Field:      "end_year",
			Value:      nil,
			Confidence: "low",
			Source:     "ultimatespecs",
			Notes:      fmt.Sprintf("Production years field exists (%s) but end year could not be parsed", prodYears),
		})
		report.RequiresReview = true
	}

	// Validate year range logic
	if startYear > 0 && endYear.Valid && startYear > int(endYear.Int64) {
		report.Issues = append(report.Issues, DataQuality{
			Field:      "year_range",
			Value:      fmt.Sprintf("%d-%d", startYear, endYear.Int64),
			Confidence: "low",
			Source:     "ultimatespecs",
			Notes:      "Start year is greater than end year - illogical range",
		})
		report.RequiresReview = true
	}

	// Validate future years
	currentYear := time.Now().Year()
	if startYear > currentYear+2 {
		report.Issues = append(report.Issues, DataQuality{
			Field:      "start_year",
			Value:      startYear,
			Confidence: "low",
			Source:     "ultimatespecs",
			Notes:      fmt.Sprintf("Start year is too far in the future (current: %d)", currentYear),
		})
		report.RequiresReview = true
	}

	// Validate power_hp
	if hp == 0 {
		report.Issues = append(report.Issues, DataQuality{
			Field:      "power_hp",
			Value:      0,
			Confidence: "missing",
			Source:     "ultimatespecs",
			Notes:      "Could not parse horsepower from specs",
		})
		report.RequiresReview = true
	}

	return report
}

// AddToReport adds a validation report to the global report
func AddToReport(vr ValidationReport) {
	if len(vr.Issues) > 0 {
		globalReport.TrimsRequiringReview = append(globalReport.TrimsRequiringReview, vr)
		globalReport.IssuesFound += len(vr.Issues)
	}
	globalReport.TotalTrims++
}

// SaveReport saves the validation report to a JSON file
func SaveReport(filename string) error {
	data, err := json.MarshalIndent(globalReport, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal report: %w", err)
	}

	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write report file: %w", err)
	}

	fmt.Printf("\n✅ Validation report saved to %s\n", filename)
	fmt.Printf("   Total trims: %d\n", globalReport.TotalTrims)
	fmt.Printf("   Issues found: %d\n", globalReport.IssuesFound)
	fmt.Printf("   Trims requiring review: %d\n", len(globalReport.TrimsRequiringReview))

	return nil
}

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
