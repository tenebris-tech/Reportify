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
		maxWidth += 4 // add padding

		// Calculate second column width
		pageWidth, _ := pdf.GetPageSize()
		left, _, right, _ := pdf.GetMargins()
		usableWidth := pageWidth - left - right
		secondColWidth := usableWidth - maxWidth

		// Draw table with wrapping second column
		for i, row := range r.headingTable {
			paddingY := 4.0 // Define padding for this row

			xLeft := pdf.GetX()
			yBefore := pdf.GetY()

			// Get page dimensions for boundary checks
			_, pageHeight := pdf.GetPageSize()
			_, top, _, bottom := pdf.GetMargins()
			pageBottomY := pageHeight - bottom

			// Now actually draw the content
			// Draw first column text with horizontal paddingX
			pdf.SetFont(r.GetFont("table_left"))
			pdf.SetXY(xLeft+paddingX, yBefore+paddingY/2)
			pdf.MultiCell(maxWidth-paddingX, r.lineHeight, row[0], "", "L", false)

			// Draw second column text
			pdf.SetFont(r.GetFont("table_right"))
			pdf.SetXY(xLeft+maxWidth+paddingX, yBefore+paddingY/2)

			// Store current page before drawing content
			currentPage := pdf.PageNo()

			// Draw the second column
			pdf.MultiCell(secondColWidth-paddingX, r.lineHeight, row[1], "", "L", false)

			// Check if a page break occurred
			newPage := pdf.PageNo()
			newY := pdf.GetY()

			// Set drawing properties
			pdf.SetLineWidth(0.3)
			pdf.SetDrawColor(0, 0, 0)

			if newPage > currentPage {
				// Content caused a page break

				// First portion - draw borders to page bottom
				firstPageHeight := pageBottomY - yBefore
				pdf.SetPage(currentPage)
				pdf.Rect(xLeft, yBefore, maxWidth, firstPageHeight, "D")
				pdf.Rect(xLeft+maxWidth, yBefore, secondColWidth, firstPageHeight, "D")

				// There may be pages between the first and last page of the table, in which case
				// we can assume that the entire page is occupied by the table
				if newPage > currentPage+1 {
					for p := currentPage + 1; p < newPage; p++ {
						pdf.SetPage(p)
						pdf.Rect(xLeft, top, maxWidth, pageHeight-top-bottom, "D")
						pdf.Rect(xLeft+maxWidth, top, secondColWidth, pageHeight-top-bottom, "D")
					}
				}
				for p := currentPage + 1; p < newPage; p++ {
					pdf.SetPage(p)
					pdf.Rect(xLeft, top, maxWidth, pageHeight-top-bottom, "D")
					pdf.Rect(xLeft+maxWidth, top, secondColWidth, pageHeight-top-bottom, "D")
				}

				// Last page portion - draw borders from top to current position
				pdf.SetPage(newPage)
				secondPageHeight := newY - top + paddingY/2
				pdf.Rect(xLeft, top, maxWidth, secondPageHeight, "D")
				pdf.Rect(xLeft+maxWidth, top, secondColWidth, secondPageHeight, "D")

				// Ensure we're on the correct page for the next row
				pdf.SetPage(newPage)
				pdf.SetY(newY + paddingY/2)
			} else {
				// Content stayed on same page
				actualHeight := newY - yBefore + paddingY/2

				// Draw borders around the content
				pdf.Rect(xLeft, yBefore, maxWidth, actualHeight, "D")
				pdf.Rect(xLeft+maxWidth, yBefore, secondColWidth, actualHeight, "D")

				// Move to position for next row
				pdf.SetY(newY + paddingY/2)
			}

			// If this is the last row, ensure we're on the correct page
			if i == len(r.headingTable)-1 {
				pdf.SetPage(newPage)
			}
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
	return pdf.OutputFileAndClose(r.outputFile + ".pdf")
}
