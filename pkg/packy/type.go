package packy

type Node struct {
	X, Y, W, H       int
	used             bool
	right, down, Fit *Node
}

func NewBlock(w, h int) *Node {
	return &Node{W: w, H: h}
}
