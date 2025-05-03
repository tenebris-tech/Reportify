// Copyright (c) 2025 Tenebris Technologies Inc.
// Released under the MIT License. See LICENSE file for details.

package reportify

// SetMarginDefaults sets the default margins for the report
func (r *Reportify) SetMarginDefaults() {
	r.SetMargins(DefaultMargins, DefaultMargins, DefaultMargins)
}

// SetMargins sets the margins for the report. For consistency, they use the
// same order as the fpdf library: left, top, right.
func (r *Reportify) SetMargins(left, top, right float64) {
	if left < 0 {
		left = 0
	}
	if right < 0 {
		right = 0
	}
	if top < 0 {
		top = 0
	}

	r.margins.Left = left
	r.margins.Right = right
	r.margins.Top = top
}

func (r *Reportify) GetTopMargin() float64 {
	return r.margins.Top
}

func (r *Reportify) GetLeftMargin() float64 {
	return r.margins.Left
}

func (r *Reportify) GetRightMargin() float64 {
	return r.margins.Right
}
