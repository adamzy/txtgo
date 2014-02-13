package tree

import (
	"bytes"
	"strconv"
)

type data struct {
	// built-in values are
	// enough for most cases
	Name   string
	Length float64
	Id     int
	Level  int

	// `Ext`: for any runtime extension.
	// e.g. one can use `node.Ext = []int{1,2,3}`,
	// to bind extra data to that node.
	// The data can be accessed by `node.Ext.(type)`,
	// e.g. `node.Ext.([]int)`.
	Ext interface{}
}

type Node struct {
	*data
	Father   *Node
	Children []*Node
}

// `newNode` will allocate space for data while `new(Node)` not
func newNode() *Node {
	return &Node{&data{}, nil, nil}
}

func (n *Node) AddChild(c *Node) {
	n.Children = append(n.Children, c)
	c.Father = n
}

func (n *Node) toByte(b *bytes.Buffer, showlength bool) {
	writeNode := func(n *Node) {
		b.WriteString(n.Name)
		if showlength {
			b.WriteRune(':')
			length := strconv.FormatFloat(n.Length, 'f', -1, 64)
			b.WriteString(length)
		}
	}

	children := n.Children
	if len(children) > 0 {
		b.WriteRune('(')
		children[0].toByte(b, showlength)
		for i := 1; i < len(children); i++ {
			b.WriteRune(',')
			children[i].toByte(b, showlength)
		}
		b.WriteRune(')')
	}
	writeNode(n)
}

// Represent the tree (rooted at `n`) in newick format.
func (n *Node) String() string {
	b := new(bytes.Buffer)
	n.toByte(b, false)
	return b.String()
}

func (n *Node) toString(showlength bool) string {
	b := new(bytes.Buffer)
	n.toByte(b, showlength)
	return b.String()
}

// Find the most left node of the tree (rooted at `node`).
// It is useful for post-order-iteration.
func leftmost(node *Node) *Node {
	for len(node.Children) > 0 {
		node = node.Children[0]
	}
	return node
}

// Postorder iterate the tree and put nodes into a list
func (node *Node) Post2List() []*Node {
	nl := make([]*Node, 0, 20)

	// A function that appends node with its descendent to `nl`
	var ap func(n *Node) //necessary for making a recursive function here
	ap = func(n *Node) {
		for _, c := range n.Children {
			ap(c)
		}
		nl = append(nl, n)
	}

	ap(node)
	return nl
}

func (node *Node) IsLeaf() bool {
	return len(node.Children) == 0
}

func (node *Node) IsRoot() bool {
	return node.Father == nil
}

func (node *Node) IsInternal() bool {
	return !node.IsLeaf()
}

func (node *Node) IsBinary() bool {
	return len(node.Children) == 2
}

func (n *Node) replaceBy(m *Node) {
	n.Children = m.Children
	for _, c := range n.Children {
		c.Father = n
	}
	n.Name = m.Name // necessary
}
