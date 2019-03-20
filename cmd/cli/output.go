package main

import (
	"errors"
	"fmt"
	"math"

	"github.com/innermond/cobai/packy/pkg/packy"
)

func aproximateHeightText(numchar int, w float64) float64 {
	wchar := w / float64(numchar+2)
	return math.Floor(1.5*wchar*100.0) / 100
}

func outsvg(blocks []*packy.Node, topleftmargin float64, plain bool) (string, error) {
	if len(blocks) == 0 {
		return "", errors.New("no blocks")
	}

	gb := svgGroupStart("id=\"blocks\"")
	if !plain {
		gb = svgGroupStart("id=\"blocks\"", "inkscape:label=\"blocks\"", "inkscape:groupmode=\"layer\"")
	}
	// first block
	blk := blocks[0]
	gb += svgRect(blk.Fit.X,
		blk.Fit.Y,
		blk.W,
		blk.H,
		"fill:magenta;stroke:none",
	)

	for _, blk := range blocks[1:] {
		if blk.Fit != nil {
			// blocks on the top edge must be shortened on height by a expand = half cutwidth
			if blk.Fit.Y == topleftmargin {
				gb += svgRect(blk.Fit.X,
					blk.Fit.Y,
					blk.W,
					blk.H,
					"fill:red;stroke:none",
				)
				continue
			}
			// blocks on the left edge must be shortened on width by a expand = half cutwidth
			if blk.Fit.X == topleftmargin {
				gb += svgRect(blk.Fit.X,
					blk.Fit.Y,
					blk.W,
					blk.H,
					"fill:green;stroke:none",
				)
				continue
			}
			// blocks that do not touch any big box edges keeps their expanded dimensions
			gb += svgRect(blk.Fit.X,
				blk.Fit.Y,
				blk.W,
				blk.H,
				"fill:#eee;stroke:none",
			)
		} else {
			return "", errors.New("unexpected unfit block")
		}
	}
	gb = svgGroupEnd(gb)

	gt := svgGroupStart("id=\"dimensions\"")
	if !plain {
		gt = svgGroupStart("id=\"dimensions\"", "inkscape:label=\"dimensions\"", "inkscape:groupmode=\"layer\"")
	}
	for _, blk := range blocks {
		if blk.Fit != nil {
			x := fmt.Sprintf("%.2fx%.2f", blk.W, blk.H)
			y := aproximateHeightText(len(x), blk.W)
			gt += svgText(blk.Fit.X+blk.W/2, blk.Fit.Y+blk.H/2+y/3, // y/3 is totally empirical
				x, "text-anchor:middle;font-size:"+fmt.Sprintf("%.2f", y)+";fill:#000")
		} else {
			return "", errors.New("unexpected unfit block")
		}
	}
	gt = svgGroupEnd(gt)

	return gb + gt, nil
}
