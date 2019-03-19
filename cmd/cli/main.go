package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/innermond/cobai/packy/pkg/packy"
)

var (
	outname, unit, dimensions, bigbox                   string
	report, output, tight, supertight, expandtocutwidth bool
	mu, ml, pp, pd, cutwidth                            float64
)

func param() {
	flag.StringVar(&outname, "o", "fit", "name of the maching project")
	flag.StringVar(&unit, "u", "mm", "unit of measurements")
	flag.StringVar(&bigbox, "bb", "0x0", "dimensions as \"wxh\" in units for bigest box / mother surface")
	flag.BoolVar(&report, "r", true, "match report")
	flag.BoolVar(&output, "f", false, "outputing files representing matching")
	flag.BoolVar(&tight, "tight", false, "when true only aria used tighten by height is taken into account")
	flag.BoolVar(&supertight, "supertight", false, "when true only aria used tighten bu height and width is taken into account")
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
	lendims := 0
	for _, dim := range dims {
		lendims += dim.n
	}

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
	if cutwidth > 0.0 {
		expand = cutwidth / 2
	}
	// first row and first column will not have to be expanded with an entire cutwidth
	// but with an expand as they will not need a cut - being already cut on external side
	// we have to expand big box but later we will shrink big box with an expand
	width += expand
	height += expand

	initialheight := height

	inx := 0
	stats := ""
	mpused := 0.0
	mplost := 0.0
	mperim := 0.0
	unfitlen := len(unfit)

	for unfitlen > 0 {
		inx++
		// Presumably unfit are already expanded
		fit, unfit = packy.PackExpand(width, initialheight, unfit, expand)

		if len(fit) == 0 {
			break
		}

		// calculate calculate de maximum height and width that fit blocks have
		xh := 0.0
		xw := 0.0
		cwx, cwy := 0.0, 0.0
		prevx, prevy := 0.0, 0.0
		for _, blk := range fit {
			if blk.Fit.Y != prevy {
				cwy++
			}
			if blk.Fit.X != prevx {
				cwx++
			}
			if xh < blk.Fit.Y+blk.H {
				xh = blk.Fit.Y + blk.H
			}
			if xw < blk.Fit.X+blk.W {
				xw = blk.Fit.X + blk.W
			}
			prevx = blk.Fit.X
			prevy = blk.Fit.Y
		}

		if supertight {
			height = xh
			width = xw
		} else if tight {
			height = xh
		}

		// output only when we have fit blocks
		if output {
			fn := fmt.Sprintf("%s.%d.svg", outname, inx)

			f, err := os.Create(fn)
			if err != nil {
				panic("cannot create file")
			}

			// revert big box to the original size
			width -= expand
			height -= expand

			s := svgStart(width, height, unit)
			si, err := outsvg(fit)
			if err != nil {
				f.Close()
				os.Remove(fn)
			} else {
				s += svgEnd(si)

				_, err = f.WriteString(s)
				if err != nil {
					panic(err)
				}
				f.Close()
			}
		}

		if report {
			aria, perim := 0.0, 0.0
			for _, blk := range fit {
				aria += blk.W * blk.H
				perim += 2 * (blk.W + blk.H)
			}

			percent := math.Round(100 * aria / (width * height))

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

			used := aria / k
			lost := (width*height)/k - used
			perim = perim / kp
			mplost += lost
			mpused += used
			mperim += perim

			stats += fmt.Sprintf(
				"%d %s %.2f%sx%.2f%s fit %d of %d unfit %d used %.2f lost %.2f percent %.2f perim %.2f\n",
				inx,
				outname,
				width, unit,
				height, unit,
				len(fit),
				unfitlen,
				unfitlen-len(fit),
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
