package theap

type MaxHeap struct {
	NumElems int

	hd *Node
}

// Push will insert a new node into the heap
// O(logn)
func (m *MaxHeap) Push(n *Node) {

}

// Pop will pop off the top of the heap
// O(1)
func (m *MaxHeap) Pop() *Node {
	hd := m.hd

	return hd
}

// Peek will return the top of the heap.
// O(1)
func (m *MaxHeap) Peek() *Node {
	return m.hd
}
