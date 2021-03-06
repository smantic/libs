package theap

type MinHeap struct {
	NumElems int

	hd *Node
}

// Push puts
func (m *MinHeap) Push(node *Node) {

	n := m.NumElems + 1
	parent := find(m.hd, n>>1)
	node.parent = parent

	if n&1 == 0 {
		parent.left = node
	} else {
		parent.right = node
	}

	m.NumElems = n

	m.up(node)
}

func (m *MinHeap) Pop() *Node {
	hd := m.hd
	last := find(m.hd, m.NumElems)
	parent := last.parent

	if parent.left == last {
		parent.left = nil
	}
	if parent.right == last {
		parent.right = nil
	}

	last.parent = nil
	last.left = hd.left
	last.right = hd.right
	m.hd = last

	m.down(last)

	return hd
}

// Peek will return the top of the heap.
func (m *MinHeap) Peek() *Node {
	return m.hd
}

// up will bouble a node up if it should
func (m *MinHeap) up(n *Node) {
	for n.parent != nil {
		if n.Value < n.parent.Value {
			swap(n, n.parent)
			n = n.parent
		} else {
			return
		}
	}
}

// down will bouble a node down if it should
func (m *MinHeap) down(n *Node) {
	for n != nil {
		if n.left == nil {
			return
		}

		if n.Value > n.left && n.left <= n.right {
			swap(n, n.left)
			n = n.left
		}
		if n.right != nil && n.Value > n.right && n.right < n.left {
			swap(n, n.right)
			n = n.right
		}
	}
}
