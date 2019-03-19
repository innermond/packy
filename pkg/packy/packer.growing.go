package packy

// GrowingPacker is an elastic packer
type GrowingPacker struct {
	root *Node
}

// Fit fits the blocks
func (gp *GrowingPacker) Fit(blocks []*Node) {
	var (
		w, h  float64
		n     int
		l     = len(blocks)
		block *Node
	)

	if l > 0 {
		w = blocks[0].W
	}
	if l > 0 {
		h = blocks[0].H
	}

	gp.root = &Node{W: w, H: h}

	for n < l {
		block = blocks[n]
		node := gp.findNode(gp.root, block.W, block.H)
		if node != nil {
			block.Fit = gp.splitNode(node, block.W, block.H)
		} else {
			block.Fit = gp.growNode(block.W, block.H)
		}
		n++
	}
}

func (gp *GrowingPacker) findNode(root *Node, w, h float64) *Node {
	var node *Node
	if root.used {
		node = gp.findNode(root.right, w, h)
		if node == nil {
			node = gp.findNode(root.down, w, h)
		}
	} else if w <= root.W && h <= root.H {
		node = root
	}

	return node
}

func (gp *GrowingPacker) splitNode(node *Node, w, h float64) *Node {
	node.used = true
	node.down = &Node{X: node.X, Y: node.Y + h, W: node.W, H: node.H - h}
	node.right = &Node{X: node.X + w, Y: node.Y, W: node.W - w, H: h}

	return node
}

func (gp *GrowingPacker) growNode(w, h float64) *Node {
	canGrowDown := w <= gp.root.W
	canGrowRight := h <= gp.root.H

	shouldGrowRight := canGrowRight && gp.root.H >= gp.root.W+w // attempt to keep square-ish by growing right when height is much greater than width
	shouldGrowDown := canGrowDown && gp.root.W >= gp.root.H+h   // attempt to keep square-ish by growing down  when width  is much greater than height

	if shouldGrowRight {
		return gp.growRight(w, h)
	} else if shouldGrowDown {
		return gp.growDown(w, h)
	} else if canGrowRight {
		return gp.growRight(w, h)
	} else if canGrowDown {
		return gp.growDown(w, h)
	}
	return nil // need to ensure sensible root starting size to avoid this happening
}

func (gp GrowingPacker) growRight(w, h float64) *Node {
	gp.root = &Node{
		used:  true,
		X:     0,
		Y:     0,
		W:     gp.root.W + w,
		H:     gp.root.H,
		down:  gp.root,
		right: &Node{X: gp.root.W, Y: 0, W: w, H: gp.root.H},
	}
	node := gp.findNode(gp.root, w, h)
	if node != nil {
		return gp.splitNode(node, w, h)
	}
	return nil
}

func (gp GrowingPacker) growDown(w, h float64) *Node {
	gp.root = &Node{
		used:  true,
		X:     0,
		Y:     0,
		W:     gp.root.W,
		H:     gp.root.H + h,
		down:  &Node{X: 0, Y: gp.root.H, W: gp.root.W, H: h},
		right: gp.root,
	}
	node := gp.findNode(gp.root, w, h)
	if node != nil {
		return gp.splitNode(node, w, h)
	}
	return nil
}
