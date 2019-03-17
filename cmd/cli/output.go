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
		bf := blk.FitDim()
		bd := blk.Dim()
		if bf != nil {
			canvas.Rect(bf.X, bf.Y, bd.W, bd.H, "fill:none;stroke-width:0.2;stroke-opacity:1;stroke:#000")
		} else {
			return errors.New("unexpected unfit block")
		}
	}
	canvas.Gend()

	canvas.Group("id=\"dimensions\"", "inkscape:label=\"dimensions\"", "inkscape:groupmode=\"layer\"")
	for _, blk := range blocks {
		bf := blk.FitDim()
		bd := blk.Dim()
		if bf != nil {
			canvas.Text(bf.X+bd.W/2, bf.Y+bd.H/2,
				fmt.Sprintf("%dx%d", bd.W, bd.H), "text-anchor:middle;font-size:72pt;fill:#000")
		} else {
			return errors.New("unexpected unfit block")
		}
	}
	canvas.Gend()

	return nil
}
