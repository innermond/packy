package main

import (
	"errors"
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

func pack(width int, height int, blocks []*Node) (fit []*Node, unfit []*Node) {
	packer := NewPacker(width, height)
	packer.Fit(blocks)

	for _, blk := range blocks {
		if blk.fit != nil {
			fit = append(fit, blk)
		} else {
			unfit = append(unfit, blk)
		}
	}
	arrange(unfit)

	return
}

func outsvg(canvas *svg.SVG, blocks []*Node) error {
	canvas.Group("id=\"blocks\"", "inkscape:label=\"blocks\"", "inkscape:groupmode=\"layer\"")
	for _, blk := range blocks {
		if blk.fit != nil {
			canvas.Rect(blk.fit.x, blk.fit.y, blk.w, blk.h, "fill:none;stroke-width:0.1;stroke-opacity:1;stroke:#000")
		} else {
			return errors.New("unexpected unfit block")
		}
	}
	canvas.Gend()

	return nil
}

func arrange(bb []*Node) {
	if len(bb) == 0 {
		return
	}

	sort.Slice(bb, func(i, j int) bool {
		a := bb[i]
		b := bb[j]
		aw := float64(a.w)
		ah := float64(a.h)
		bw := float64(b.w)
		bh := float64(b.h)

		if math.Max(bw, bh) <= math.Max(aw, ah) {
			return true
		}

		if math.Min(bw, bh) <= math.Min(aw, ah) {
			return true
		}

		if bh <= ah {
			return true
		}

		if bw <= aw {
			return true
		}

		return false
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

func dimString(dimarr []string) (dims []dim) {
	for _, dd := range dimarr {
		d := strings.Split(dd, "x")
		if len(d) == 2 {
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
