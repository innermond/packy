package packy

type Node struct {
	X, Y, W, H       int
	used             bool
	right, down, Fit *Node
}
