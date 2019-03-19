package main

import (
	"errors"
	"fmt"
	"math"

	"github.com/innermond/cobai/packy/pkg/packy"
)

func aproximateHeightText(numchar int, w float64) string {
	wchar := w / float64(numchar+2)
	hchar := math.Floor(1.5*wchar*100.0) / 100
	return fmt.Sprintf("%.2f", hchar)
}

func outsvg(blocks []*packy.Node, expand float64) (string, error) {
	gb := svgGroupStart("id=\"blocks\"", "inkscape:label=\"blocks\"", "inkscape:groupmode=\"layer\"")
	for _, blk := range blocks {
		if blk.Fit != nil {
			// first row and first column must be shrink by expand
			prevx, prevy := 0.0, 0.0
			if blk.Fit.Y != prevy || blk.Fit.X != prevx {
				gb += svgRect(blk.Fit.X-expand,
					blk.Fit.Y,
					blk.W-expand,
					blk.H-expand,
					"fill:none;stroke-width:0.2;stroke-opacity:1;stroke:green")
			} else {
				gb += svgRect(blk.Fit.X,
					blk.Fit.Y,
					blk.W-expand,
					blk.H-expand,
					"fill:none;stroke-width:0.2;stroke-opacity:1;stroke:red")

			}
			prevx, prevy = blk.Fit.X, blk.Fit.Y
		} else {
			return "", errors.New("unexpected unfit block")
		}
	}
	gb = svgGroupEnd(gb)

	gt := svgGroupStart("id=\"dimensions\"", "inkscape:label=\"dimensions\"", "inkscape:groupmode=\"layer\"")
	for _, blk := range blocks {
		if blk.Fit != nil {
			x := fmt.Sprintf("%.2fx%.2f", blk.W, blk.H)
			gt += svgText(blk.Fit.X+blk.W/2, blk.Fit.Y+blk.H/2,
				x, "text-anchor:middle;font-size:"+aproximateHeightText(len(x), blk.W)+";fill:#000")
		} else {
			return "", errors.New("unexpected unfit block")
		}
	}
	gt = svgGroupEnd(gt)

	return gb + gt, nil
}
