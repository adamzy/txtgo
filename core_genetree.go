package tree

type NoTaxaError struct {
	Taxa string
}

func (err NoTaxaError) Error() string {
	return "Cannot find taxa '" + err.Taxa + "' in species tree."
}

type Lcamap struct {
	G   *Tree
	S   *SpeciesTree
	Map []*Node
}

// Generate `Lcamap` for gene tree and species tree.
// Assume that every internal genetree node has at least two children
func LcaMap(t *Tree, st *SpeciesTree) (*Lcamap, error) {
	maps, err := genLcaMap(t, st)
	if err != nil {
		return nil, err
	}
	return &Lcamap{t, st, maps}, nil
}

func genLcaMap(t *Tree, st *SpeciesTree) ([]*Node, error) {
	taxon := st.Taxon
	lca := st.Lca
	lcamap := make([]*Node, len(t.Nodes))

	for i, n := range t.Nodes {
		if n.IsLeaf() {
			sn, ok := taxon[n.Name]
			if !ok {
				return nil, NoTaxaError{n.Name}
			}
			lcamap[i] = sn
		} else {
			sn := lca(lcamap[n.Children[0].Id], lcamap[n.Children[1].Id])
			for j := 2; j < len(n.Children); j++ {
				sn = lca(sn, lcamap[n.Children[j].Id])
			}
			lcamap[i] = sn
		}
	}
	return lcamap, nil
}
