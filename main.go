package main

import (
	"flag"
	"fmt"
	"math"
	"os"

	svg "github.com/ajstarks/svgo"
)

var (
	outname, unit, dimensions string
	report, output            bool
)

func initFlag() {
	flag.StringVar(&outname, "o", "fit", "name of the maching project")
	flag.StringVar(&unit, "u", "mm", "unit of measurements")
	flag.BoolVar(&report, "r", true, "match report")
	flag.BoolVar(&output, "f", false, "outputing files representing matching")

	flag.Parse()
}

func main() {
	initFlag()
	dimensions := flag.Args()

	dims := dimString(dimensions)
	unfit := blocksArranged(dims)
	fit := []*Node{}

	width := 2030
	height := 3050
	var canvas *svg.SVG
	unfitlen := len(unfit)
	inx := 0
	stats := ""
	for unfitlen > 0 {
		inx++
		fit, unfit = pack(width, height, unfit)

		if output {
			f, err := os.Create(fmt.Sprintf("%s.%d.svg", outname, inx))
			if err != nil {
				panic("cannot create file")
			}
			canvas = svg.New(f)
			canvas.Startunit(width, height, unit, fmt.Sprintf("viewBox=\"0 0 %d %d\"", width, height))
			outsvg(canvas, fit)
			canvas.End()
		}

		if report {
			aria := 0
			for _, blk := range fit {
				aria += blk.w * blk.h
			}
			percent := math.Round(100 * float64(aria) / float64(width*height))
			stats += fmt.Sprintf(
				"%d %s %d%sx%d%s fit %d percent %.2f\n",
				inx,
				outname,
				width, unit,
				height, unit,
				len(fit),
				percent,
			)
		}

		if unfitlen == len(unfit) {
			break
		}
		unfitlen = len(unfit)
	}

	if report {
		fmt.Print(stats)
	}
}
