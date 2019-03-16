package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

	svg "github.com/ajstarks/svgo"
)

var (
	outname, unit, dimensions, bigbox string
	report, output, tight             bool
)

func initFlag() {
	flag.StringVar(&outname, "o", "fit", "name of the maching project")
	flag.StringVar(&unit, "u", "mm", "unit of measurements")
	flag.StringVar(&bigbox, "b", "0x0", "dimensions as \"wxh\" in units")
	flag.BoolVar(&report, "r", true, "match report")
	flag.BoolVar(&output, "f", false, "outputing files representing matching")
	flag.BoolVar(&tight, "tight", false, "when true only aria used is taken into account")

	flag.Parse()
}

func main() {
	initFlag()
	dimensions := flag.Args()

	dims := dimString(dimensions)

	unfit := blocksArranged(dims)
	fit := []*Node{}

	wh := strings.Split(bigbox, "x")
	width, err := strconv.Atoi(wh[0])
	if err != nil {
		panic("can't get width")
	}
	height, err := strconv.Atoi(wh[1])
	if err != nil {
		panic("can't get height")
	}

	initialheight := height

	var canvas *svg.SVG
	unfitlen := len(unfit)
	inx := 0
	stats := ""
	mpused := 0.0
	mplost := 0.0
	for unfitlen > 0 {
		inx++
		fit, unfit = pack(width, initialheight, unfit)

		if tight {
			h := 0
			for _, blk := range fit {
				if h < blk.fit.y+blk.h {
					h = blk.fit.y + blk.h
				}
			}
			height = h
		}

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
			k := 1.0
			switch unit {
			case "mm":
				k = 1000 * 1000
			case "cm":
				k = 100 * 100
			}
			used := float64(aria) / k
			mpused += used
			mplost += float64(width*height)/k - used
		}

		if unfitlen == len(unfit) {
			break
		}
		unfitlen = len(unfit)
	}

	if report {
		stats += fmt.Sprintf("used %.2f lost %.2f\n", mpused, mplost)
		fmt.Print(stats)
	}
}
