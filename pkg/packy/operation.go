package packy

import (
	"math"
)

// Operation holds context for boxes matching operation
type Operation struct {
	Width         float64
	Height        float64
	Cutwidth      float64
	TopLeftMargin float64
}

// NewOperation create an operation
func NewOperation(w, h, cw, tl float64) *Operation {
	return &Operation{w, h, cw, tl}
}

// Expand temporary adjust page to accomodate by cut width or top left offset
func (o *Operation) Expand() {
	expandpage := o.Cutwidth / 2
	if o.TopLeftMargin > 0.0 {
		expandpage = 0.0
	}
	o.Width = o.Width + expandpage - 2*o.TopLeftMargin
	o.Height = o.Height + expandpage - 2*o.TopLeftMargin
}

// Unexpand reverts page to original dimensions
func (o *Operation) Unexpand() {
	expandpage := o.Cutwidth / 2
	if o.TopLeftMargin > 0.0 {
		expandpage = 0.0
	}
	o.Width = o.Width - expandpage + 2*o.TopLeftMargin
	o.Height = o.Height - expandpage + 2*o.TopLeftMargin
}

// Pack puts the boxes in places
func (o *Operation) Pack(unfit []*Node) (report []*Report) {

	unfitlen := len(unfit)
	expand := o.Cutwidth

	o.Expand()
	Arrange(unfit)

	var fit []*Node

	for unfitlen > 0 {
		// Presumably unfit are already expanded
		fit, unfit = PackExpand(o.Width, o.Height, unfit, expand, o.TopLeftMargin)

		if len(fit) == 0 || unfitlen == len(unfit) {
			break
		}

		// calculate calculate de maximum height and width that fit blocks have
		xh := 0.0
		xw := 0.0
		cwx, cwy := 0.0, 0.0
		prevx, prevy := o.TopLeftMargin, o.TopLeftMargin
		for _, blk := range fit {
			if blk.Fit.Y != prevy {
				cwy++
			}
			if blk.Fit.X != prevx {
				cwx++
			}
			if xh < blk.Fit.Y+blk.H {
				xh = blk.Fit.Y + blk.H
			}
			if xw < blk.Fit.X+blk.W {
				xw = blk.Fit.X + blk.W
			}
			prevx = blk.Fit.X
			prevy = blk.Fit.Y
		}

		aria, perim := 0.0, 0.0
		for _, blk := range fit {
			aria += blk.W * blk.H
			perim += 2 * (blk.W + blk.H)
		}

		report = append(report,
			&Report{
				WidthUsed:  xw,
				HeightUsed: xh,
				Aria:       aria,
				Perimeter:  perim,
				Fit:        fit,
			},
		)

		unfitlen = len(unfit)
	}

	o.Unexpand()

	return
}

// ReportOne prepare data for one big box used for fitting
// k is a unit to unit conversion - 1000 mm -> m
func (o *Operation) ReportOne(pkd *Report, k float64, modeReportAria string) []float64 {
	width, height := o.Width/k, o.Height/k
	switch modeReportAria {
	case "tight":
		height = pkd.HeightUsed / k
	case "supertight":
		height = pkd.HeightUsed / k
		width = pkd.WidthUsed / k
	}
	aria, perim := pkd.Aria/(k*k), pkd.Perimeter/k

	lost := width*height - aria
	percent := math.Round(100 * aria / (width * height))

	return []float64{aria, lost, percent, perim}
}

// Report has data needed for stats
type Report struct {
	WidthUsed  float64
	HeightUsed float64

	Aria      float64
	Perimeter float64

	Fit []*Node
}
