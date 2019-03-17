package main

import (
	"errors"
	"fmt"
	"math"

	svg "github.com/ajstarks/svgo"
	"github.com/innermond/cobai/packy/pkg/packy"
)

func aproximateHeightText(numchar int, w int) string {
	wchar := float64(w) / float64(numchar+2)
	hchar := math.Floor(1.5*wchar*100.0) / 100
	return fmt.Sprintf("%.2f", hchar)
}

func outsvg(canvas *svg.SVG, blocks []*packy.Node) error {
	canvas.Group("id=\"blocks\"", "inkscape:label=\"blocks\"", "inkscape:groupmode=\"layer\"")
	for _, blk := range blocks {
		if blk.Fit != nil {
			canvas.Rect(blk.Fit.X, blk.Fit.Y, blk.W, blk.H, "fill:none;stroke-width:0.2;stroke-opacity:1;stroke:#000")
		} else {
			return errors.New("unexpected unfit block")
		}
	}
	canvas.Gend()

	canvas.Group("id=\"dimensions\"", "inkscape:label=\"dimensions\"", "inkscape:groupmode=\"layer\"")
	for _, blk := range blocks {
		if blk.Fit != nil {
			x := fmt.Sprintf("%dx%d", blk.W, blk.H)
			canvas.Text(blk.Fit.X+blk.W/2, blk.Fit.Y+blk.H/2,
				x, "text-anchor:middle;font-size:"+aproximateHeightText(len(x), blk.W)+";fill:#000")
		} else {
			return errors.New("unexpected unfit block")
		}
	}
	canvas.Gend()

	return nil
}
