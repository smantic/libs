package theap

type Interface interface {
	Push(n *Node)
	Pop(n *Node)
	Peek() *Node
}
