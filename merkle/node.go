package merkle

type node struct {
	parent *node
	left   *node
	right  *node

	hashInfo string
	data     []byte
}

func (n *node) isLeaf() bool {
	return n.left == nil && n.right == nil
}

func (n *node) equals(right *node) bool {
	return n.hashInfo == right.hashInfo
}
