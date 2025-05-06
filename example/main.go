// Copyright (c) 2025 Tenebris Technologies Inc.
// Released under the MIT License. See LICENSE file for details.

package main

import (
	"fmt"
	"time"

	"github.com/tenebris-tech/Reportify"
)

func main() {

	// For testing purposes, out a file in the current directory
	outputFile := "example-report"

	// Create a header table for the report
	meta := [][]string{
		{"Report", "Sample report version 1.0"},
		{"Generated", time.Now().Format("2006-01-02 15:04:05 EDT")},
		{"Subject", "This report was generated to demonstrate using the reportify package."},
	}

	// Start the report
	report := reportify.New(
		reportify.WithOutputFile(outputFile),
		reportify.WithHeader("Sample Report"),
		reportify.WithHeadingTable(meta),
		reportify.WithEndMarker(true),
		reportify.WithPageOrientation("P"), // P for portrait, L for landscape
	) // Create a sample report

	// Add some information to the report
	report.AddH1("This is H1 centered", "C")
	report.AddH2("This is H2 left", "L")
	report.AddPara("This is a paragraph. It is left aligned by default.")
	report.AddPara("This is another paragraph.")
	report.AddH2("Part Two")
	report.AddPara("This is more information.\nLine breaks are fully supported and long lines will wrap. Spaces will be added between block of text. This is just more text to make it wrap.")
	report.AddPara("This is another paragraph.")

	// Generate the report
	err := report.Generate()
	if err != nil {
		fmt.Println("Error generating report:", err)
	}
}
