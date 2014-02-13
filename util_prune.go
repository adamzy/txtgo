package tree

// Prune tree according to the taxonmap
func (t *Tree) PruneFromTaxon(taxon Taxonmap) {
	// the node with m[node.Id] = true is marked
	// as leaf in original tree.
	m := make([]bool, len(t.Nodes))

	// the index of each node
	// node.Father.Children[node.Index] = node
	d := make([]int, len(t.Nodes))
	for i, n := range t.Nodes {
		if n.IsLeaf() {
			m[i] = true
		} else {
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

	for i, n := range t.Nodes {
		if m[i] {
			if _, ok := taxon[n.Name]; !ok {
				removeChild(n.Father, d[i])
			}
		} else if len(n.Children) == 1 {
			//contractNode(n)
			n.replaceBy(n.Children[0])
		} else if len(n.Children) == 0 {
			if n.IsRoot() {
				// if all leaves are pruned
				// the result tree will be empty.
				t.Node = nil
			}
			removeChild(n.Father, d[i])
		}
	}
	t.Update()
}

// Prune tree according to another tree
func (t *Tree) PruneFromTree(gt *Tree) {
	taxon, _ := gt.TaxonMap()
	t.PruneFromTaxon(taxon)
}
