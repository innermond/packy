package packy

// Node is a rectangle like repesenting a block
type Node struct {
	X, Y, W, H       float64
	used             bool
	right, down, Fit *Node
}

// Dim represents size
type Dim struct {
	W, H float64
	N    int // how many
}

// DimNode transforms a bunch of Dims into echivalent Nodes
func DimNode(dims []Dim) (blocks []*Node) {
	for inx := 0; inx < len(dims); inx++ {
		for i := 0; i < dims[inx].N; i++ {
			blocks = append(blocks, &Node{W: dims[inx].W, H: dims[inx].H})
		}
	}
	return
}

// DimFlat see how many Dim (dim.N) and builds them as we have a flat list of dims
func DimFlat(dims []Dim) (blocks []Dim) {
	for inx := 0; inx < len(dims); inx++ {
		for i := 0; i < dims[inx].N; i++ {
			blocks = append(blocks, Dim{W: dims[inx].W, H: dims[inx].H})
		}
	}
	return
}
