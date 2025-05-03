// Copyright (c) 2025 Tenebris Technologies Inc.
// Released under the MIT License. See LICENSE file for details.

package reportify

import (
	"fmt"

	"codeberg.org/go-pdf/fpdf"
)

// Generate creates the report
func (r *Reportify) Generate() error {
	paddingX := 1.0
	paddingY := 4.0
	topMargin := r.GetTopMargin()

	pdf := fpdf.New(r.pageOrientation, "mm", "Letter", "")

	// Expand top margin if required for the header
	if r.header != "" {
		pdf.SetFont(r.GetFont("header"))
		_, lineHeight := pdf.GetFontSize() // in points
		topMargin += lineHeight * 1.5      // 1.5x line height for padding
	}

	// Set margins
	pdf.SetMargins(r.GetLeftMargin(), topMargin, r.GetRightMargin())
	pdf.AliasNbPages("{nb}") // placeholder for total pages

	// Optionally add header to every page
	if r.header != "" {
		pdf.SetHeaderFunc(func() {
			currentY := pdf.GetY()           // Save current Y position
			pdf.SetY(10)                     // Set Y position
			pdf.SetX(r.GetLeftMargin())      // Set X position
			pdf.SetFont(r.GetFont("header")) // Set font
			pdf.CellFormat(0, 10, r.header, "", 0, r.getHeaderAlign(), false, 0, "")
			pdf.SetY(currentY) // Restore Y position
		})
	}

	// Add footer to every page
	pdf.SetFooterFunc(func() {
		pdf.SetY(-15) // position from the bottom of the page
		pdf.SetFont(r.GetFont("footer"))
		pdf.CellFormat(0, 10,
			fmt.Sprintf("Page %d of {nb}", pdf.PageNo()), "", 0, r.getFooterAlign(), false, 0, "")
	})

	pdf.AddPage()
	pdf.SetFont(r.GetFont("default"))

	if len(r.headingTable) > 0 {

		// Find max width of first column
		maxWidth := 0.0
		for _, row := range r.headingTable {
			w := pdf.GetStringWidth(row[0])
			if w > maxWidth {
				maxWidth = w
			}
		}
		maxWidth += 4 // add paddingY

		// Calculate second column width
		pageWidth, _ := pdf.GetPageSize()
		left, _, right, _ := pdf.GetMargins()
		usableWidth := pageWidth - left - right
		secondColWidth := usableWidth - maxWidth

		// Draw table with wrapping second column
		for _, row := range r.headingTable {

			lines := pdf.SplitLines([]byte(row[1]), secondColWidth)
			rowHeight := float64(len(lines)) * r.lineHeight
			rowHeightWithPadding := rowHeight + paddingY

			xLeft := pdf.GetX()
			yBefore := pdf.GetY()

			pdf.SetLineWidth(0.3)
			pdf.SetDrawColor(0, 0, 0)

			// Draw first column box
			pdf.Rect(xLeft, yBefore, maxWidth, rowHeightWithPadding, "D")

			// Draw first column text with horizontal paddingY
			pdf.SetFont(r.GetFont("table_left"))
			pdf.SetXY(xLeft+paddingX, yBefore+paddingY/2)
			pdf.MultiCell(maxWidth-paddingX, r.lineHeight, row[0], "", "L", false)

			// Draw second column box
			pdf.SetFont(r.GetFont("table_right"))
			pdf.Rect(xLeft+maxWidth, yBefore, secondColWidth, rowHeightWithPadding, "D")
			pdf.SetXY(xLeft+maxWidth+paddingX, yBefore+paddingY/2)
			pdf.MultiCell(secondColWidth-paddingX, r.lineHeight, row[1], "", "L", false)

			// Move Y to next row
			pdf.SetY(yBefore + rowHeightWithPadding)
		}

		// Add space after the table
		pdf.SetY(pdf.GetY() + 5)
	}

	// Add the text blocks using MultiCell for text wrapping and automatic page breaks
	for _, b := range r.blocks {
		pdf.SetFont(r.GetFont(b.kind))
		pdf.MultiCell(0, r.lineHeight, b.text, "", r.getAlign(&b), false)
		pdf.Ln(r.lineHeight) // One line between blocks
	}

	if r.endMarker {
		pdf.SetFont(r.GetFont("end_marker"))
		pdf.MultiCell(0, r.lineHeight, "### end of report ###", "", "C", false)
	}

	// Output the PDF
	return pdf.OutputFileAndClose(r.outputFile)
}
