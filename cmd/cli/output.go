package main

import (
	"errors"
	"fmt"

	svg "github.com/ajstarks/svgo"
	"github.com/innermond/cobai/packy/pkg/packy"
)

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
			canvas.Text(blk.Fit.X+blk.W/2, blk.Fit.Y+blk.H/2,
				fmt.Sprintf("%dx%d", blk.W, blk.H), "text-anchor:middle;font-size:72pt;fill:#000")
		} else {
			return errors.New("unexpected unfit block")
		}
	}
	canvas.Gend()

	return nil
}
