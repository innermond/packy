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
	} else if w <= root.W && h <= root.H {
		node = root
	}

	return node
}

func (pk *Packer) splitNode(node *Node, w, h int) *Node {
	node.used = true
	node.down = &Node{X: node.X, Y: node.Y + h, W: node.W, H: node.H - h}
	node.right = &Node{X: node.X + w, Y: node.Y, W: node.W - w, H: h}

	return node
}

func (pk *Packer) Fit(blocks []*Node) {
	var (
		n     int
		block *Node
	)

	for n < len(blocks) {
		block = blocks[n]
		node := pk.findNode(pk.root, block.W, block.H)
		if node != nil {
			block.Fit = pk.splitNode(node, block.W, block.H)
		}
		n++
	}
}
