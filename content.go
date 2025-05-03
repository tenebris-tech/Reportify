// Copyright (c) 2025 Tenebris Technologies Inc.
// Released under the MIT License. See LICENSE file for details.

package reportify

import (
	"strings"
)

// AddPara adds a paragraph to the report with optional alignment
func (r *Reportify) AddPara(s string, opts ...string) {
	align := ""
	if len(opts) > 0 && opts[0] != "" {
		align = opts[0]
	}
	r.AddBlock(s, "p", align)
}

// AddH1 adds a header 1 to the report with optional alignment
func (r *Reportify) AddH1(s string, opts ...string) {
	align := ""
	if len(opts) > 0 && opts[0] != "" {
		align = opts[0]
	}
	r.AddBlock(s, "h1", align)
}

// AddH2 adds a header 2 to the report with optional alignment
func (r *Reportify) AddH2(s string, opts ...string) {
	align := ""
	if len(opts) > 0 && opts[0] != "" {
		align = opts[0]
	}
	r.AddBlock(s, "h2", align)
}

// AddH3 adds a header 3 to the report with optional alignment
func (r *Reportify) AddH3(s string, opts ...string) {
	align := ""
	if len(opts) > 0 && opts[0] != "" {
		align = opts[0]
	}
	r.AddBlock(s, "h3", align)
}

// AddBlock adds to the report using the specified type and alignment
func (r *Reportify) AddBlock(s string, kind string, align string) {
	r.blocks = append(r.blocks,
		block{text: s, kind: kind, align: align})
}

// getAlign returns the alignment for the block, ensuring that it is one of L, C, or R
// If the block has an alignment specified, it will be used. Otherwise, it will fall back to
// the default alignment for the font and then to the default font's alignment.
func (r *Reportify) getAlign(b *block) string {

	// Start with the block's alignment
	a := b.align

	// If no alignment is specified, use the applicable font's alignment
	// or fall back to the default
	if a == "" {
		if f, ok := r.fonts[b.kind]; ok {
			a = f.Align
		} else {
			a = r.defaultFont.Align
		}
	}
	return constrainAlign(a)
}

// getHeaderAlign returns the alignment for the header
func (r *Reportify) getHeaderAlign() string {
	return r.getFontAlign("header")
}

// getFooterAlign returns the alignment for the footer
func (r *Reportify) getFooterAlign() string {
	return r.getFontAlign("footer")
}

// getFontAlign returns the alignment for the font
func (r *Reportify) getFontAlign(k string) string {
	a := r.defaultFont.Align

	if f, ok := r.fonts[k]; ok {
		if f.Align != "" {
			a = f.Align
		}
	}
	return constrainAlign(a)
}

// constrainAlign constrains the alignment to L, C, or R
func constrainAlign(a string) string {
	switch strings.ToUpper(a) {
	case "L":
		return "L"
	case "C":
		return "C"
	case "R":
		return "R"
	default:
		return "L"
	}
}
