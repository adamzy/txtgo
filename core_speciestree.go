package tree

import (
	"txtgo/rmq"
)

type SpeciesTree struct {
	*Tree
	//Taxon map[string]*Node // a map that maps taxa to the unique leaf node
    Taxon Taxonmap
	Lca   func(a, b *Node) *Node
}

func (t *Tree) SpeciesTree() *SpeciesTree {
	st := new(SpeciesTree)
	st.Tree = t
	st.Taxon = t.TaxonMap()
	st.Lca = t.LCAer()
	return st
}

// Make a LCAer function for species tree t.
// Usage:
//      Lca := st.LCAer()
//      Lca(a, b)
func (t *Tree) LCAer() func(a, b *Node) *Node {
	// Euler tour visit each internal node 3 times,
	// and leaf node once, thus the total length is at most
	// len(nodes)*3.
	array := make([]int64, t.Size*3)

	// each node in t has an id in the rmq array
	mapid := make([]int64, t.Size)

	// each id in rmq array also associate a node
	nodes := make([]*Node, t.Size*3)
	visit := make([]int, t.Size)

	//fmt.Println("rmq")
	var index int64
	n := t.Node
	for n != nil {
		array[index] = int64(n.Level)
		nodes[index] = n

		id := n.Id
		mapid[id] = index //mapid[id] may update 3 times, it's fine

		//visit[id] = true
		//fmt.Println(index, n.Name, n)
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

	rmqer := rmq.ResRMQ(array)
	return func(a, b *Node) *Node {
		aid := mapid[a.Id]
		bid := mapid[b.Id]
		//fmt.Println(a.Name, aid, b.Name, bid, len(array))
		p, _ := rmqer(aid, bid)
		return nodes[p]
	}
}

// This one only works for binary tree, has been replaced by LCAer().
// Make a LCAer function for species tree t.
// Usage:
//      Lca := st.LCAer()
//      Lca(a, b)
func (t *Tree) lCAerOld() func(a, b *Node) *Node {
	// Euler tour visit each internal node 3 times,
	// and leaf node once, thus the total length is at most
	// len(nodes)*3.
	array := make([]int64, t.Size*3)

	// each node in t has an id in the rmq array
	mapid := make([]int64, t.Size)

	// each id in rmq array also associate a node
	nodes := make([]*Node, t.Size*3)
	visit := make([]bool, t.Size)

	var index int64
	n := t.Node
	for n != nil {
		array[index] = int64(n.Level)
		nodes[index] = n

		id := n.Id
		mapid[id] = index //mapid[id] may update 3 times, it's fine

		visit[id] = true
		index++

		switch {
		case n.IsLeaf():
			n = n.Father
		case !visit[n.Children[0].Id]:
			n = n.Children[0]
		case !visit[n.Children[1].Id]:
			n = n.Children[1]
		default:
			n = n.Father
		}
	}
	array = array[:index]

	rmqer := rmq.ResRMQ(array)
	return func(a, b *Node) *Node {
		aid := mapid[a.Id]
		bid := mapid[b.Id]
		p, _ := rmqer(aid, bid)
		return nodes[p]
	}
}
