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
	contracter(t, checkByLength(length))
}

// Randomly contract tree branches according to the given rate.
func (t *Tree) RandomContract(rate float64) {
	contracter(t, checkRandom(rate))
}

func checkByLength(length float64) func(*Node) bool {
	return func(n *Node) bool {
		if n.Length < length {
			return true
		}
		return false
	}
}

func checkRandom(rate float64) func(*Node) bool {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return func(n *Node) bool {
		if r.Float64() < rate {
			return true
		}
		return false
	}
}

func contracter(t *Tree, checker func(*Node) bool) {
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
		if !n.IsLeaf() && checker(n) {
			Contract(n)
		}
	}
	t.Update()
}
