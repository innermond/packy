package packy

// Node is a rectangle like repesenting a block
type Node struct {
	X, Y, W, H       float64
	used             bool
	right, down, Fit *Node
}
