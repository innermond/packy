package main

import (
	"strconv"
	"strings"

	"github.com/innermond/cobai/packy/pkg/packy"
)

type dim struct {
	w, h float64
	n    int
}

func blocksfrom(dims []dim) (blocks []*packy.Node) {
	for inx := 0; inx < len(dims); inx++ {
		for i := 0; i < dims[inx].n; i++ {
			blocks = append(blocks, &packy.Node{W: int(dims[inx].w), H: int(dims[inx].h)})
		}
	}
	return
}

func blocksArranged(dims []dim) []*packy.Node {
	bb := blocksfrom(dims)
	packy.Arrange(bb)
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
