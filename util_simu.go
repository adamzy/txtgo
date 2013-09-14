package tree

import (
	"math/rand"
	"strconv"
	"time"
)

func simuTree(size int) *Tree {
	if size < 1 {
		return nil
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	root := newNode()
	leaves := []*Node{root}

	duplicate := func(n int) []*Node {
		c1 := newNode()
		c2 := newNode()
		node := leaves[n]
		node.AddChild(c1)
		node.AddChild(c2)
		leaves[n] = c1
		leaves = append(leaves, c2)
		return leaves
	}

	for len(leaves) < size {
		leaves = duplicate(r.Intn(len(leaves)))
	}

	// root.UpdateIndex()
	t := new(Tree)
	t.Node = root
	t.Update()
	return t
}

func SimuTree(size int) *Tree {
	t := simuTree(size)
	if t.Node != nil {
		t.AssignLeafName()
	}
	return t
}

func SimuTreeRandomTaxon(size, ntaxon int) *Tree {
	t := simuTree(size)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	//for n:=root.PostStart();n!=nil;n=n.PostNext(){
	for _, n := range t.Nodes {
		if n.IsLeaf() {
			n.Name = strconv.Itoa(r.Intn(ntaxon))
			//n.Name = "l_"+strconv.Itoa(r.Intn(ntaxon))
		}
	}
	return t
}

func (t *Tree) RandomContract(rate float64) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	d := make([]int, t.Size)
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
	for i := 0; i < t.Size-1; i++ {
		n := t.Nodes[i]
		if !n.IsLeaf() && r.Float64() < rate {
			Contract(n)
		}
	}
	t.Update()
}

func (t *Tree) AssignLeafName() {
	count := 0
	for _, node := range t.Nodes {
		if node.IsLeaf() {
			node.Name = strconv.Itoa(count)
			count++
		}
	}
}
