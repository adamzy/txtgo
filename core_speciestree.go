package tree

import (
	"errors"
	//Use `rmq` package which inplements the linear-time algorithm for solving RMQ problem.
	"txtgo/rmq"
)

var (
	SpeciesTreeNotUniquelyLabeledError = errors.New("Species tree is not uniquely labeled.")
)

type SpeciesTree struct {
	*Tree
	Taxon Taxonmap
	Lca   func(a, b *Node) *Node
}

func (t *Tree) SpeciesTree() (*SpeciesTree, error) {
	st := new(SpeciesTree)
	st.Tree = t
	var unique bool
	st.Taxon, unique = t.TaxonMap()
	if !unique {
		return nil, SpeciesTreeNotUniquelyLabeledError
	}
	st.Lca = t.LCAer()
	return st, nil
}

// Make a LCAer function for species tree t.
//
//      Usage:
//          Lca := st.LCAer()
//          Lca(a, b)
func (t *Tree) LCAer() func(a, b *Node) *Node {
	size := len(t.Nodes)
	// Euler tour visit each internal node 3 times,
	// and leaf node once, thus the total length is at most
	// `len(nodes)*3`.
	array := make([]int64, size*3)

	// each node in `t` has an id in the rmq array
	mapid := make([]int64, size)

	// each id in rmq array also associate a node
	nodes := make([]*Node, size*3)
	visit := make([]int, size)

	// ##### Translate the LCA problem to a RMQ problem
	// Firstly, generate the Euler tour of the tree.
	var index int64
	n := t.Node
	for n != nil {
		array[index] = int64(n.Level)
		nodes[index] = n
		id := n.Id

		//`mapid[id]` may be updated 3 times, it's fine
		mapid[id] = index

		index++

		switch {
		case n.IsLeaf():
			n = n.Father
		case visit[n.Id] < len(n.Children):
			n = n.Children[visit[n.Id]]
			visit[n.Father.Id]++
		default:
			n = n.Father
		}
	}
	array = array[:index]

	// Now, solve the RMQ problem.
	rmqer := rmq.ResRMQ(array)
	// Return the solver as a function.
	return func(a, b *Node) *Node {
		aid := mapid[a.Id]
		bid := mapid[b.Id]
		p, _ := rmqer(aid, bid)
		return nodes[p]
	}
}
