package theap

import "math"

type Node struct {
	Value int

	left  *Node
	right *Node

	// doublely linked, otherwise down() is O(logn^2)
	parent *Node
}

// find will find the indexed node given a root and number of nodes contained in its subtree.
// it runs in log(n)
func find(root *Node, n int) *Node {
	n = n / 2
	bits := int(math.Log2(float64(n))) + 1

	cur := root
	for bits != 0 {
		if n>>bits == 0 {
			cur = cur.left
		} else {
			cur = cur.right
		}
		n = n / 2
		bits--
	}
	return cur
}

// swap two nodes
func swap(a, b *Node) {
	a.Value, b.Value = b.Value, a.Value
}
