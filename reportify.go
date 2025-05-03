// Copyright (c) 2025 Tenebris Technologies Inc.
// Released under the MIT License. See LICENSE file for details.

// Package reportify uses the fpdf library to create simple reports in PDF format.
package reportify

import "strings"

const (
	DefaultBaseFontSize = 11.0
	DefaultBaseFont     = "Arial"
	DefaultMargins      = 20.0
)

type Reportify struct {
	outputFile      string
	pageOrientation string
	defaultFont     Font
	fonts           map[string]Font
	lineHeight      float64
	header          string
	headingTable    [][]string
	blocks          []block
	endMarker       bool
	margins         Margins
}

type block struct {
	kind  string // Block type (e.g., "p", "h1", "h2", "h3")
	text  string // Text content
	align string // Optional alignment (e.g., "L", "C", "R")
}

type Font struct {
	Family string  // Font family (e.g., "Arial", "Times")
	Style  string  // Font style (e.g., "B" for bold, "I" for italic)
	Size   float64 // Font size in points
	Align  string  // Optional alignment (e.g., "L", "C", "R")
}

type Margins struct {
	Left  float64 // Left margin in mm
	Right float64 // Right margin in mm
	Top   float64 // Top margin in mm
}

type Option func(*Reportify)

// New creates a new LogReport with the provided options
func New(opts ...Option) *Reportify {

	// Create default LogReport
	r := &Reportify{
		fonts:           make(map[string]Font),
		lineHeight:      6.0,
		pageOrientation: "P",
	}

	// Set default margins
	r.SetMarginDefaults()

	// Apply all options
	for _, opt := range opts {
		opt(r)
	}

	// Set the default fonts
	// This wil respect defaults set via the options
	r.SetFontDefaults()

	return r
}

// WithOutputFile sets the output file
func WithOutputFile(outputFile string) Option {
	return func(r *Reportify) {
		r.outputFile = outputFile
	}
}

// WithHeadingTable creates a simple table at the start of the report
func WithHeadingTable(h [][]string) Option {
	return func(r *Reportify) {
		r.headingTable = h
	}
}

// WithEndMarker adds a marker at the end of the report
func WithEndMarker(endMarker bool) Option {
	return func(r *Reportify) {
		r.endMarker = endMarker
	}
}

// WithHeader adds a classification to the header
func WithHeader(header string) Option {
	return func(r *Reportify) {
		r.header = header
	}
}

// WithDefaultFont sets the default font
func WithDefaultFont(font string, style string, size float64, alignment string) Option {
	return func(r *Reportify) {
		r.defaultFont = Font{
			Family: font,
			Style:  style,
			Align:  alignment,
			Size:   size,
		}
	}
}

// WithPageOrientation sets the page orientation (L or P)
func WithPageOrientation(orientation string) Option {
	return func(r *Reportify) {

		// must be "L" or "P"
		switch strings.ToUpper(orientation) {
		case "L":
			r.pageOrientation = "L"
		case "P":
			r.pageOrientation = "P"
		default:
			r.pageOrientation = "P"
		}
	}
}

// WithDefaultAlignment sets the default alignment
// This will be used where no alignment is specified
func WithDefaultAlignment(align string) Option {
	return func(r *Reportify) {
		r.defaultFont.Align = align
	}
}

// WithMargins sets the margins
func WithMargins(left, top, right float64) Option {
	return func(r *Reportify) {
		r.SetMargins(left, top, right)
	}
}
