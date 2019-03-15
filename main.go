package main

import (
	"io"
	"math"
	"os"
	"sort"
	"strconv"

	svg "github.com/ajstarks/svgo"
)

type dim struct {
	w, h float64
	n    int
}

func main() {
	dims := []dim{
		{100, 100, 2},
		{20, 60, 7},
		{50, 100, 1},
		{70, 130, 2},
		{50, 50, 3},
		{50, 100, 1},
		{50, 50, 1},
	}

	sort.Slice(dims, func(i, j int) bool {
		a := dims[i]
		b := dims[j]

		if math.Max(b.w, b.h) < math.Max(a.w, a.h) {
			return true
		}

		if math.Min(b.w, b.h) < math.Min(a.w, a.h) {
			return true
		}

		if b.h < a.h {
			return true
		}

		return b.w < a.w

	})

	var blocks []*Node
	for inx := 0; inx < len(dims); inx++ {
		for i := 0; i < dims[inx].n; i++ {
			blocks = append(blocks, &Node{w: int(dims[inx].w), h: int(dims[inx].h)})
		}
	}
	width := 200
	height := 200
	unfit := blocks[:]
	var canvas *svg.SVG
	unfitlen := len(unfit)
	inx := 0
	for unfitlen > 0 {
		f, err := os.Create("fit." + strconv.Itoa(inx) + ".svg")
		if err != nil {
			panic("cannot create file")
		}
		inx++
		unfit, canvas = pack(width, height, unfit, f)
		canvas.End()
		if unfitlen == len(unfit) {
			break
		}
		unfitlen = len(unfit)
	}
}

func pack(width int, height int, blocks []*Node, out io.Writer) (unfit []*Node, canvas *svg.SVG) {
	packer := NewPacker(width, height)
	packer.Fit(blocks)

	canvas = svg.New(out)
	canvas.Start(width, height)
	canvas.Group("id=\"blocks\"", "inkscape:label=\"blocks\"", "inkscape:groupmode=\"layer\"")
	for _, blk := range blocks {
		if blk.fit != nil {
			canvas.Rect(blk.fit.x, blk.fit.y, blk.w, blk.h, "fill:none;stroke-width:0.1;stroke-opacity:1;stroke:#000")
		} else {
			unfit = append(unfit, blk)
		}
	}
	canvas.Gend()

	return
}
