package main

import (
	"strconv"
	"strings"

	"github.com/innermond/cobai/packy/pkg/packy"
)

func dimString(dimarr []string, extra float64) (dims []packy.Dim) {
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

		dims = append(dims, packy.Dim{W: w + extra, H: h + extra, N: n})
	}
	return
}
