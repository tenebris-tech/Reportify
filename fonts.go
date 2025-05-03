// Copyright (c) 2025 Tenebris Technologies Inc.
// Released under the MIT License. See LICENSE file for details.

package reportify

import "strings"

// SetFontDefaults uses either the configured default font or, in lieu of that, the default font
// specified in reportify.go constants. In the process, it ensures that the default font is set.
// This function is called by New() but is exported so that it can be called again to reset the fonts
// or easily switch to another font.
func (r *Reportify) SetFontDefaults() {
	baseSize := DefaultBaseFontSize
	baseFont := DefaultBaseFont

	// DefaultFonts must be set
	if r.defaultFont.Family != "" {
		baseFont = r.defaultFont.Family
	} else {
		r.defaultFont.Family = baseFont
	}

	if r.defaultFont.Size != 0 {
		baseSize = r.defaultFont.Size
	} else {
		r.defaultFont.Size = baseSize
	}

	r.SetFont("h1", baseFont, "B", baseSize+3, "L")
	r.SetFont("h2", baseFont, "B", baseSize+2, "L")
	r.SetFont("h3", baseFont, "B", baseSize+1, "L")
	r.SetFont("p", baseFont, "", baseSize, "L")
	r.SetFont("header", baseFont, "", baseSize, "L")
	r.SetFont("footer", baseFont, "", baseSize, "C")
	r.SetFont("table_left", baseFont, "B", baseSize, "C")
	r.SetFont("table_right", baseFont, "", baseSize, "C")
	r.SetFont("end_marker", baseFont, "B", baseSize, "C")
}

// SetFont sets the font for the specified "kind" (which is used because type is a reserved word)
func (r *Reportify) SetFont(kind string, family string, style string, size float64, alignment string) {
	r.fonts[kind] = Font{Family: family, Style: style, Size: size, Align: alignment}
}

// SetDefaultFont sets the default font for the report and is exported so that a calling package
// can set the default. This is helpful before calling SetFontDefaults() to reset the fonts.
func (r *Reportify) SetDefaultFont(family string, style string, size float64, alignment string) {
	r.defaultFont = Font{Family: family, Style: style, Size: size, Align: alignment}
}

// GetFont returns the font for the specified "kind" in a tuple that can be consumed by fpdf
// If an unknown "kind" is specified, the default font is returned
// Alignment, if used, is accessed via getAlign()
func (r *Reportify) GetFont(kind string) (string, string, float64) {
	var f Font
	k := strings.ToLower(kind)
	if _, ok := r.fonts[k]; !ok {
		f = r.defaultFont
	} else {
		f = r.fonts[k]
	}
	return f.Family, f.Style, f.Size
}
