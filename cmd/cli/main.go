package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

	svg "github.com/ajstarks/svgo"
	"github.com/innermond/cobai/packy/pkg/packy"
)

var (
	outname, unit, dimensions, bigbox       string
	report, output, tight, expandtocutwidth bool
	mu, ml, pp, pd, cutwidth                float64
)

func param() {
	flag.StringVar(&outname, "o", "fit", "name of the maching project")
	flag.StringVar(&unit, "u", "mm", "unit of measurements")
	flag.StringVar(&bigbox, "b", "0x0", "dimensions as \"wxh\" in units for bigest box / mother surface")
	flag.BoolVar(&report, "r", true, "match report")
	flag.BoolVar(&output, "f", false, "outputing files representing matching")
	flag.BoolVar(&tight, "tight", false, "when true only aria used is taken into account")
	flag.Float64Var(&mu, "mu", 15.0, "used material price per 1 square meter")
	flag.Float64Var(&ml, "ml", 5.0, "lost material price per 1 square meter")
	flag.Float64Var(&pp, "pp", 0.25, "perimeter price per 1 linear meter; used for evaluating cuts price")
	flag.Float64Var(&pd, "pd", 10, "travel price to location")
	flag.Float64Var(&cutwidth, "cutwidth", 0.0, "the with of material that is lost due to a cut")

	flag.Parse()
}

func main() {
	param()

	dimensions := flag.Args()
	// if the cut can eat half of its width along cutline
	// we compensate expanding boxes with an entire cut width
	dims := dimString(dimensions, cutwidth)

	unfit := blocksArranged(dims)
	fit := []*packy.Node{}

	wh := strings.Split(bigbox, "x")
	width, err := strconv.ParseFloat(wh[0], 64)
	if err != nil {
		panic("can't get width")
	}
	height, err := strconv.ParseFloat(wh[1], 64)
	if err != nil {
		panic("can't get height")
	}

	expand := 0.0
	if expandtocutwidth {
		expand = cutwidth / 2
	}
	// first row and first column will not have to be expanded with an entire cutwidth
	// but with an expand as they will not need a cut - being already cut on external side
	// we have to expand big box but later we will shrink big box with an expand
	width += expand
	height += expand

	initialheight := height

	canvas := &svg.SVG{}
	inx := 0
	stats := ""
	mpused := 0.0
	mplost := 0.0
	mperim := 0.0
	unfitlen := len(unfit)

	for unfitlen > 0 {
		inx++

		fit, unfit = packy.Pack(width, initialheight, unfit)

		// How much the boxes are spreading on page. Are they dangerously approaching to edges
		xh := 0
		xw := 0
		cwx, cwy := 0.0, 0.0
		prevx, prevy := 0, 0
		for _, blk := range fit {
			if blk.Fit.Y != prevy {
				cwx++
			}
			if blk.Fit.X != prevx {
				cwy++
			}
			if xh < blk.Fit.Y+blk.H {
				xh = blk.Fit.Y + blk.H
			}
			if xw < blk.Fit.X+blk.W {
				xw = blk.Fit.X + blk.W
			}
			prevx = blk.Fit.X
			prevy = blk.Fit.Y
			if cwx*cutwidth+float64(xw) >= float64(width) {
				panic(fmt.Sprintf("for cut width %.2f big box width %d is not enough", cutwidth, width))
			}
			if cwy*cutwidth+float64(xh) >= float64(height) {
				panic(fmt.Sprintf("for cut width %.2f big box height %d is not enough", cutwidth, height))
			}
		}

		if tight {
			height = float64(xh)
		}

		if output {
			f, err := os.Create(fmt.Sprintf("%s.%d.svg", outname, inx))
			if err != nil {
				panic("cannot create file")
			}

			// revert big box to the original size
			width -= expand
			height -= expand

			canvas = svg.New(f)
			canvas.Startunit(width, height, unit, fmt.Sprintf("viewBox=\"0 0 %d %d\"", width, height))
			outsvg(canvas, fit, expand)
			canvas.End()
		}

		if report {
			aria, perim := 0, 0.0
			for _, blk := range fit {
				aria += blk.W * blk.H
				perim += 2 * float64(blk.W+blk.H)
			}

			percent := math.Round(100 * float64(aria) / float64(width*height))

			k := 1.0
			kp := 1.0
			switch unit {
			case "mm":
				k = 1000 * 1000
				kp = 1000
			case "cm":
				k = 100 * 100
				kp = 100
			}

			used := float64(aria) / k
			lost := float64(width*height)/k - used
			perim = perim / kp
			mplost += lost
			mpused += used
			mperim += perim

			stats += fmt.Sprintf(
				"%d %s %d%sx%d%s fit %d used %.2f lost %.2f percent %.2f perim %.2f\n",
				inx,
				outname,
				width, unit,
				height, unit,
				len(fit),
				used,
				lost,
				percent,
				perim,
			)
		}

		if unfitlen == len(unfit) {
			break
		}
		unfitlen = len(unfit)
	}

	if report {
		price := mpused*mu + mplost*ml + mperim*pp + pd
		stats += fmt.Sprintf("used %.2f lost %.2f total %.2f perim %.2f price %.2f\n", mpused, mplost, mpused+mplost, mperim, price)
		fmt.Print(stats)
	}
}
