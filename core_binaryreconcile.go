package tree

// BinaryCost finds the duplication, loss, and deep coalescence cost 
// for input binary gene tree and species tree.
func BinaryCost(gt *Tree, st *SpeciesTree) (dup, loss, dc int, err error){
    lm, err := LcaMap(gt, st)

    if err!=nil {return 0,0,0, err}

    for _, n:=range gt.Nodes {
        if !n.IsLeaf() {
            m := lm.Map[n.Id]
            l := lm.Map[n.Children[0].Id]
            r := lm.Map[n.Children[1].Id]
            loss+=(l.Level + r.Level - 2*m.Level)
            dc += (l.Level + r.Level - 2*m.Level)
            if m==l || m==r {
                dup++
            } else {
                loss-=2
            }
        }    
    }
    dc -= len(st.Nodes)-1
    return
}
