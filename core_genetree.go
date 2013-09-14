package tree

import (
	"errors"
)

/*type GeneTree struct {*/
//*Tree
//Map []*Node
//}

type Lcamap struct {
	G   *Tree
	S   *SpeciesTree
	Map []*Node
}

// Assume that every internal genetree node has at least two children
func LcaMap(t *Tree, st *SpeciesTree) (*Lcamap, error) {
	maps, err := genLcaMap(t, st)
	if err != nil {
		return nil, err
	}
	return &Lcamap{G: t, S: st, Map: maps}, nil
}

func genLcaMap(t *Tree, st *SpeciesTree) ([]*Node, error) {
	taxon := st.Taxon
	//fmt.Println(taxon)
	lca := st.Lca
	lcamap := make([]*Node, t.Size)

	for i, n := range t.Nodes {
		if n.IsLeaf() {
			sn, ok := taxon[n.Name]
			if !ok {
				return nil, errors.New("Cannot find taxa '" + n.Name + "' in species tree.")
			}
			lcamap[i] = sn
		} else {
			sn := lca(lcamap[n.Children[0].Id], lcamap[n.Children[1].Id])
			//println("lca: ~~ ", lcamap[n.Children[0].Id].Name, lcamap[n.Children[1].Id].Name, sn.Name)
			for j := 2; j < len(n.Children); j++ {
				//println(sn.Name)
				//ssn := sn
				sn = lca(sn, lcamap[n.Children[j].Id])
				//println("lca: ~~ ", ssn.Name, lcamap[n.Children[j].Id].Name, sn.Name)
			}
			lcamap[i] = sn
		}
		//println("lca: ", n.Name, lcamap[i].Name)
	}
	return lcamap, nil
}
