package tree

func weightedCostUncp(wdup, wdc float64) (func(*Node)) {
    return func (root *Node, st *SpeciesTree, m []*Node) {
        sl := st.Nodes
        W := make([]int, st.Size)
    
        for _, c := root.Children {
            W[m[c.Id]] ++
        }
        for i
    }

}
