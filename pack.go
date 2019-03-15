package main

import (
	"io"
	"math"
	"sort"
	"strconv"
	"strings"

	svg "github.com/ajstarks/svgo"
)

type dim struct {
	w, h float64
	n    int
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

	arrange(unfit)

	return
}

func arrange(bb []*Node) {
	sort.Slice(bb, func(i, j int) bool {
		a := bb[i]
		b := bb[j]
		aw := float64(a.w)
		ah := float64(a.h)
		bw := float64(b.w)
		bh := float64(b.h)

		if math.Max(bw, bh) < math.Max(aw, ah) {
			return true
		}

		if math.Min(bw, bh) < math.Min(aw, ah) {
			return true
		}

		if bh < ah {
			return true
		}

		return bw < aw
	})
}

func blocksfrom(dims []dim) (blocks []*Node) {
	for inx := 0; inx < len(dims); inx++ {
		for i := 0; i < dims[inx].n; i++ {
			blocks = append(blocks, &Node{w: int(dims[inx].w), h: int(dims[inx].h)})
		}
	}
	return
}

func blocksArranged(dims []dim) []*Node {
	bb := blocksfrom(dims)
	arrange(bb)
	return bb
}

func dimString(dimstr string) (dims []dim) {
	dimarr := strings.Fields(dimstr)
	for _, dd := range dimarr {
		d := strings.Split(dd, "x")
		if len(d) < 3 {
			d = append(d, "1")
		}
		w, err := strconv.ParseFloat(d[0], 64)
		if err != nil {
			panic(err)
		}
		h, err := strconv.ParseFloat(d[1], 64)
		if err != nil {
			panic(err)
		}
		n, err := strconv.Atoi(d[2])
		if err != nil {
			panic(err)
		}

		dims = append(dims, dim{w: w, h: h, n: n})
	}
	return
}
