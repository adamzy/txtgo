package tree

// RF finds the Robinson-Folds distance for uniquely
// labeled input binary gene tree and species tree.
func RF(gt *Tree, st *SpeciesTree) (rf int, err error) {
	lm, err := LcaMap(gt, st)

	if err != nil {
		return 0, err
	}

	sd := make([]int, len(st.Nodes))
	for i, s := range st.Nodes {
		if s.IsLeaf() {
			sd[i] = 1
		} else {
			a, b := s.Children[0], s.Children[1]
			sd[i] = sd[a.Id] + sd[b.Id]
		}
	}

	gd := make([]int, len(gt.Nodes))
	for i, g := range gt.Nodes {
		if g.IsLeaf() {
			gd[i] = 1
			rf++
		} else {
			a, b := g.Children[0], g.Children[1]
			gd[i] = gd[a.Id] + gd[b.Id]
			if gd[i] == sd[lm.Map[g.Id].Id] {
				rf--
			}
		}
	}
	rf--
	return
}
