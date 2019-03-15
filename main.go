package main

import (
	"os"
	"strconv"

	svg "github.com/ajstarks/svgo"
)

func main() {
	dimstr := "100x100x2 20x60x7 50x100 70x130x2 50x50x3 50x100x1 50x50"
	dims := dimString(dimstr)
	blocks := blocksArranged(dims)

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
