package tree

import (
	"math/rand"
	"time"
)

// Contract tree branches which have only one child.
func (t *Tree) ContractSingleChild() {
	for _, n := range t.Nodes {
		if len(n.Children) == 1 {
			n.replaceBy(n.Children[0])
		}
	}
	t.Update()
}

// Contract tree branches whose lengthes are less than the given one.
func (t *Tree) ContractByLength(length float64) {
	size := len(t.Nodes)
	d := make([]int, size)
	for _, n := range t.Nodes {
		if n.IsInternal() {
			for j, c := range n.Children {
				d[c.Id] = j
			}
		}
	}

	removeChild := func(n *Node, i int) {
		a := n.Children

		// update the index of last child
		d[a[len(a)-1].Id] = i // important!

		a[len(a)-1], a[i], a = nil, a[len(a)-1], a[:len(a)-1]
		n.Children = a
	}

	Contract := func(n *Node) {
		f := n.Father
		cn := n.Children
		removeChild(n.Father, d[n.Id])
		for _, c := range cn {
			f.AddChild(c)
		}
	}

	// avoid root
	for i := 0; i < size-1; i++ {
		n := t.Nodes[i]
		if !n.IsLeaf() && n.Length < length {
			Contract(n)
		}
	}
	t.Update()
}

// Randomly contract tree branches according to the given rate.
func (t *Tree) RandomContract(rate float64) {
	size := len(t.Nodes)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	d := make([]int, size)
	for _, n := range t.Nodes {
		if n.IsInternal() {
			for j, c := range n.Children {
				d[c.Id] = j
			}
		}
	}

	removeChild := func(n *Node, i int) {
		a := n.Children

		// update the index of last child
		d[a[len(a)-1].Id] = i // important!

		a[len(a)-1], a[i], a = nil, a[len(a)-1], a[:len(a)-1]
		n.Children = a
	}

	Contract := func(n *Node) {
		f := n.Father
		cn := n.Children
		removeChild(n.Father, d[n.Id])
		for _, c := range cn {
			f.AddChild(c)
		}
	}

	// avoid root
	for i := 0; i < size-1; i++ {
		n := t.Nodes[i]
		if !n.IsLeaf() && r.Float64() < rate {
			Contract(n)
		}
	}
	t.Update()
}
