package packy

import (
	"math"
	"sort"
)

// Pack fits given blocks
func Pack(width float64, height float64, blocks []*Node) (fit []*Node, unfit []*Node) {
	packer := NewPacker(width, height)
	packer.Fit(blocks)

	for _, blk := range blocks {
		if blk.Fit != nil {
			fit = append(fit, blk)
		} else {
			unfit = append(unfit, blk)
		}
	}
	Arrange(unfit)

	return
}

// PackExpand adjust blocks dimensions in porder to compensate for cut width
func PackExpand(width float64, height float64, blocks []*Node, expand float64, topleftmargin float64) (fit []*Node, unfit []*Node) {
	fit, unfit = Pack(width, height, blocks)
	if len(fit) == 0 {
		return fit, unfit
	}

	if topleftmargin > 0.0 {
		for _, blk := range fit {
			if blk.Fit != nil {
				blk.Fit.X += topleftmargin
				blk.Fit.Y += topleftmargin
			}
		}
		return fit, unfit
	}

	// first block
	blk := fit[0]
	blk.W -= expand
	blk.H -= expand

	for _, blk := range fit[1:] {
		if blk.Fit != nil {
			// blocks on the top edge must be shortened on height by a expand = half cutwidth
			if blk.Fit.Y == 0.0 {
				blk.Fit.X -= expand
				blk.H -= expand
				continue
			}
			// blocks on the left edge must be shortened on width by a expand = half cutwidth
			if blk.Fit.X == 0.0 {
				blk.Fit.Y -= expand
				blk.W -= expand
				continue
			}
			// blocks that do not touch any big box edges keeps their expanded dimensions
			blk.Fit.X -= expand
			blk.Fit.Y -= expand
		}
	}
	return fit, unfit
}

// Arrange sorts blocks bigest to smallest
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
