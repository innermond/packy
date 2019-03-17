package packy

type Packer struct {
	root *Node
}

func NewPacker(w, h int) *Packer {
	return &Packer{root: &Node{W: w, H: h}}
}

func (pk *Packer) findNode(root *Node, w, h int) *Node {
	var node *Node
	if root.used {
		node = pk.findNode(root.right, w, h)
		if node == nil {
			node = pk.findNode(root.down, w, h)
		}
	} else if w <= root.w && h <= root.h {
		node = root
	}

	return node
}

func (pk *Packer) splitNode(node *Node, w, h int) *Node {
	node.used = true
	node.down = &Node.X: node.X, y: node.y + h, w: node.w, h: node.h - h}
	node.right = &Node.X: node.X + w, y: node.y, w: node.w - w, h: h}

	return node
}

func (pk *Packer) Fit(blocks []*Node) {
	var (
		n     int
		block *Node
	)

	for n < len(blocks) {
		block = blocks[n]
		node := pk.findNode(pk.root, block.w, block.h)
		if node != nil {
			block.fit = pk.splitNode(node, block.w, block.h)
		}
		n++
	}
}