package tree

import (
	"fmt"
	"testing"
)

func Test_Node(t *testing.T) {
	n := newNode()
	n.Name = "hello"
}

func Test_MakePrint(t *testing.T) {
	//s := "((a,b ), ( c:3.2, (d , e , fff:-23 ) D ): 23 ) root"
	s := "((a,b )A, ( c:3.2, (d , eee :-23 ) D )C: 23 ) root"
	tree, err := Make(s)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(tree)
	fmt.Println(tree.toString(true))
	fmt.Println("size:", tree.Size)

	tree.Nodes = tree.Post2List()
	for i := range tree.Nodes {
		n := tree.Nodes[i]
		fmt.Println(i, n.Name, n.IsBinary())
	}
	fmt.Println(tree.IsBinary())
	fmt.Println(tree.IsRoot())
	sss := new(Node)
	fmt.Println(sss.IsRoot())
	fmt.Println(sss.IsLeaf())
	fmt.Println(sss.IsBinary())
}

func Test_In2List(t *testing.T) {
	s := "((a,b )A, ( c:3.2, (d , eee :-23 ) D )C: 23 ) root"
	tree, err := Make(s)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(tree)
	il := tree.In2List()
	for i := range il {
		n := il[i]
		fmt.Println(i, n.Name)
	}
	fmt.Println(tree.Name)
}

func Test_LcaMap(t *testing.T) {
	gs := "((a,b)A,(c,(c,d)D)D)R;"
	ss := "((a,b)A,(c,d)D)R;"
	gt, err := Make(gs)
	if err != nil {
		t.Log(err)
	}
	st, err := Make(ss)
	if err != nil {
		t.Log(err)
	}

	sst := st.SpeciesTree()
	lm, err := LcaMap(gt, sst)
	if err != nil {
		t.Log(err)
	}
	for i, sn := range lm.Map {
		fmt.Println(i, gt.Nodes[i].Name, sn.Name)
		if gt.Nodes[i].Name != sn.Name {
			t.Log("Incorrect LCA!")
		}
	}

	dup, loss, dc, err := BinaryCost(gt, sst)
	if err != nil {
		t.Log(err)
	}
	fmt.Println(dup, loss, dc)
	if dup != 1 || loss != 1 || dc != 1 {
		t.Log("Incorrect binry cost.")
	}
	/*fmt.Println(st.Size)*/
	//for i, n := range st.Nodes {
	//fmt.Println(i, n.Name, n.Id)
	/*}*/
}

func Test_Nodeex(t *testing.T) {
	n := newNode()
	pre := make([]int, 2)
	pre[0] = 11
	pre[1] = 21
	n.Ext = pre
	// n.Ext = append(n.Ext, 31) // doesn't work
	p := n.Ext.([]int)
	fmt.Println(p)
}

func Test_RF(t *testing.T) {
	gs := "(((1,2),(3,4)),(5,6))"
	ss := "(((1,3),(2,6)),(4,5))"
	gt, err := Make(gs)
	if err != nil {
		fmt.Println(err)
	}
	st, err := Make(ss)
	if err != nil {
		fmt.Println(err)
	}
	rf, err := RF(gt, st.SpeciesTree())
	if err != nil {
		fmt.Println(err)
	}
	if rf != 4 {
		t.Fatal("Incorrect RF")
	}
}

func Test_RF2(t *testing.T) {
	gs := "(((1,2),(3,4)),(5,6))"
	ss := "(((1,3),(2,4)),(6,5))"
	gt, err := Make(gs)
	if err != nil {
		fmt.Println(err)
	}
	st, err := Make(ss)
	if err != nil {
		fmt.Println(err)
	}
	rf, err := RF(gt, st.SpeciesTree())
	if err != nil {
		fmt.Println(err)
	}
	if rf != 2 {
		t.Fatal("Incorrect RF")
	}
}
