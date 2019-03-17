package packy

import (
	"math"
	"sort"
)

func Pack(width int, height int, blocks []*Node) (fit []*Node, unfit []*Node) {
	packer := NewPacker(width, height)
	packer.Fit(blocks)

	for _, blk := range blocks {
		if blk.fit != nil {
			fit = append(fit, blk)
		} else {
			unfit = append(unfit, blk)
		}
	}
	Arrange(unfit)

	return
}

func Arrange(bb []*Node) {
	if len(bb) == 0 {
		return
	}

	sort.Slice(bb, func(i, j int) bool {
		a := bb[i]
		b := bb[j]
		aw := float64(a.W)
		ah := float64(a.H)
		bw := float64(b.W)
		bh := float64(b.H)

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
