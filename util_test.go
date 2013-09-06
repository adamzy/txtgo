package tree

import (
	"fmt"
	"testing"
)

func Test_Prune(t *testing.T) {
	gs := "((a,a)A,(c,(c,d)D)D)R;"
	ss := "((a,b)A,(c,(d,e))D)R;"
	gt, err := Make(gs)
	if err != nil {
		t.Log(err)
	}
	st, err := Make(ss)
	if err != nil {
		t.Log(err)
	}
	st.PruneFromTree(gt)
	// fmt.Println(st)
	// for i, n := range st.Nodes {
	// 	fmt.Println(i, n.Name, len(n.Children))
	// }
	sst := st.SpeciesTree()
	_, err = LcaMap(gt, sst)
	if err != nil {
		t.Log(err)
	}

    _, _, _, err = BinaryCost(gt, sst)
    if err != nil {
        t.Log(err)
    }

}

func Test_Simu(t *testing.T) {
	size := 100
	ntaxon := 50
	tree := SimuTree(size)
	fmt.Println(tree)
	tree2 := SimuTreeRandomTaxon(size, ntaxon)
	fmt.Println(tree2)
	tree2.RandomContract(0.8)
	fmt.Println(tree2)
}

















