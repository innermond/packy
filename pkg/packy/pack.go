package packy

import (
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
	var fitzero []*Node
	fit, unfit = Pack(width, height, blocks)
	if len(fit) == 0 {
		return fitzero, blocks
	}

	if topleftmargin > 0.0 {
		// margin physically must have room for a half cut width
		if topleftmargin <= expand/2 {
			return fitzero, blocks
		}
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
	blk.W -= expand / 2
	blk.H -= expand / 2

	for _, blk := range fit[1:] {
		if blk.Fit != nil {
			// blocks on the top edge must be shortened on height by a expand = half cutwidth
			if blk.Fit.Y == 0.0 {
				blk.Fit.X -= expand
				blk.H -= expand / 2
				continue
			}
			// blocks on the left edge must be shortened on width by a expand = half cutwidth
			if blk.Fit.X == 0.0 {
				blk.Fit.Y -= expand
				blk.W -= expand / 2
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
		// aw := a.W
		ah := a.H
		// bw := b.W
		bh := b.H

		if ah <= bh {
			return false
		}

		// sort boxes about ratio
		// rr := (aw / ah) / (bw / bh)
		// if rr < 1.1 && rr > 0.9 {
		// 	return false
		// }

		// then about aria
		// if aw*ah <= bw*bh {
		// 	return false
		// }

		// if math.Min(aw, ah) <= math.Min(bw, bh) {
		// 	return false
		// }

		// if math.Max(aw, ah) <= math.Max(bw, bh) {
		// 	return false
		// }

		// if aw <= bw {
		// 	return false
		// }

		return true
	})
}
