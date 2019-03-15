package main

type Node struct {
	x, y, w, h       int
	used             bool
	right, down, fit *Node
}

type GrowingPacker struct {
	root *Node
}

func (gp *GrowingPacker) Fit(blocks []*Node) {
	var (
		w, h, n int
		l       int = len(blocks)
		block   *Node
	)

	if l > 0 {
		w = blocks[0].w
	}
	if l > 0 {
		h = blocks[0].h
	}

	gp.root = &Node{w: w, h: h}

	for n < l {
		block = blocks[n]
		node := gp.findNode(gp.root, block.w, block.h)
		if node != nil {
			block.fit = gp.splitNode(node, block.w, block.h)
		} else {
			block.fit = gp.growNode(block.w, block.h)
		}
		n++
	}
}

func (gp *GrowingPacker) findNode(root *Node, w, h int) *Node {
	var node *Node
	if root.used {
		node = gp.findNode(root.right, w, h)
		if node == nil {
			node = gp.findNode(root.down, w, h)
		}
	} else if w <= root.w && h <= root.h {
		node = root
	}

	return node
}

func (gp *GrowingPacker) splitNode(node *Node, w, h int) *Node {
	node.used = true
	node.down = &Node{x: node.x, y: node.y + h, w: node.w, h: node.h - h}
	node.right = &Node{x: node.x + w, y: node.y, w: node.w - w, h: h}

	return node
}

func (gp *GrowingPacker) growNode(w, h int) *Node {
	canGrowDown := w <= gp.root.w
	canGrowRight := h <= gp.root.h

	shouldGrowRight := canGrowRight && gp.root.h >= gp.root.w+w // attempt to keep square-ish by growing right when height is much greater than width
	shouldGrowDown := canGrowDown && gp.root.w >= gp.root.h+h   // attempt to keep square-ish by growing down  when width  is much greater than height

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

func (gp GrowingPacker) growRight(w, h int) *Node {
	gp.root = &Node{
		used:  true,
		x:     0,
		y:     0,
		w:     gp.root.w + w,
		h:     gp.root.h,
		down:  gp.root,
		right: &Node{x: gp.root.w, y: 0, w: w, h: gp.root.h},
	}
	node := gp.findNode(gp.root, w, h)
	if node != nil {
		return gp.splitNode(node, w, h)
	}
	return nil
}

func (gp GrowingPacker) growDown(w, h int) *Node {
	gp.root = &Node{
		used:  true,
		x:     0,
		y:     0,
		w:     gp.root.w,
		h:     gp.root.h + h,
		down:  &Node{x: 0, y: gp.root.h, w: gp.root.w, h: h},
		right: gp.root,
	}
	node := gp.findNode(gp.root, w, h)
	if node != nil {
		return gp.splitNode(node, w, h)
	}
	return nil
}
