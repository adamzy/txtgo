package tree

import (
	"strconv"
	//"strings"
	//"errors"
	"bytes"
)

type data struct {
	// built-in values
	// enough for most of time
	Name   string
	Length float64
	Id     int
	Level  int

	// Ext: for any runtime extension.
	// e.g. one can use node.Ext = []int{1,2,3},
	// to bind extra data to that node.
	// The data can be accessed by node.Ext.(type),
	// e.g. a := node.Ext.([]int).
	Ext   interface{} 
}

type Node struct {
	*data
	Father   *Node
	Children []*Node
}

/*func newNode2() *Node{*/
//node := new(Node)
//node.data = new(data)
//return node
/*}*/

// newNode will allocate space for data while new(Node) not
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
			//b.WriteByte(byte(':'))
			b.WriteRune(':')
			length := strconv.FormatFloat(n.Length, 'f', -1, 64)
			b.WriteString(length)
		}
	}

	children := n.Children
	if len(children) > 0 {
		//b.WriteByte(byte('('))
		b.WriteRune('(')
		children[0].toByte(b, showlength)
		for i := 1; i < len(children); i++ {
			//b.WriteByte(byte(','))
			b.WriteRune(',')
			//writeNode(children[i])
			children[i].toByte(b, showlength)
		}
		//b.WriteByte(byte(')'))
		b.WriteRune(')')
	}
	writeNode(n)
}

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

func leftmost(node *Node) *Node {
	for len(node.Children) > 0 {
		node = node.Children[0]
	}
	return node
}

// Postorder iterate the tree and put nodes into a list
func (node *Node) Post2List() []*Node {
	//var nl []*Node
	nl := make([]*Node, 0, 20)

	// A function that appends node with its descandent to nl
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

// Inorder iterate the tree and put nodes into a list.
// It's only defined for binary tree.
// It will cause panic for non-bianry tree.
// Note: Not finished yet.
// TODO Fix this. Morri's iteration?
func (node *Node) in2List() []*Node {
	nl := make([]*Node, 0, 20)

	flag := false
	for n := leftmost(node); n.Father != nil; {
		if flag {
			n = n.Children[1]
		} else {
			n = n.Father
		}
	}
	nl = append(nl, node)
	return nl
}

// Inorder iterate the tree and put nodes into a list.
// It's only defined for binary tree.
// It will cause panic for non-bianry tree.
// Ugly!
func (node *Node) In2List() []*Node {
	//var nl []*Node
	nl := make([]*Node, 0, 20)

	// A function that appends node with its descandent to nl
	var ap func(n *Node) //necessary for making a recursive function here
	ap = func(n *Node) {
		if n.IsInternal() {
			ap(n.Children[0])
		}
		nl = append(nl, n)
		if n.IsInternal() {
			ap(n.Children[1])
		}
	}

	ap(node)
	return nl
}

func (n *Node) replaceBy(m *Node) {
	n.Children = m.Children
	for _, c := range n.Children {
		c.Father = n
	}
	n.Name = m.Name // necessary
}
