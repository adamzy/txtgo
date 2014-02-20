package tree

import (
	"fmt"
	"testing"
)

func Test_Lineartimegt(t *testing.T) {
	var gs, ss string
	//gs = "((a,b,c)A,(d,e)D)R;"
	//gs = "(a,a,(a,b),b,c,d,d);"
	//gs = "(a,b,a);"
	gs = "(a,a,a)"
	//gs = "(a,a,b,b,c,d,d)"
	ss = "((a,b)A,(c,(d,e)D)C)R;"

	gs = "((1,4),2,2,1,1)"
	ss = "(((0,1),2),(3,4))"

	gs = "(0,0,4,2,1,2)"
	ss = "((0,1)A,((2,3)B,4)C)R"

	gs = "((0,2,2),(1,0),3)"
	ss = "(0,((1,2)A,(3,4)B)C)R"

	//gs = "((1,1,3),(4,4,1))"
	gs = "(1,1,3)r"
	ss = "((((0,1)A,2)B,3)C,4)R"
	//gs = "(0,2,2)"
	//ss = "(0, 2)"

	gs = "((1,2),(3,1),3,1,3)"
	ss = "((0,((1,2)A,3)B)C,4)R;"

	gt, err := Make(gs)
	if err != nil {
		t.Log(err)
	}
	st, err := Make(ss)
	if err != nil {
		t.Log(err)
	}
	//st.PruneFromTree(gt)
	sst, err := st.SpeciesTree()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(sst.IsBinary())
	fmt.Println(gt)
	fmt.Println(sst)
	err = RefineGt(gt, sst, "mutation", 3, 1)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(gt)
	fmt.Println(BinaryCost(gt, sst))
}

func Test_Mutation1(t *testing.T) {
	gs := "(1,1,7,7,8,8)"
	ss := "(((((((1,2),3),4),5),6),7),8);"

	gt, _ := Make(gs)
	st, _ := Make(ss)
	sst, err := st.SpeciesTree()
	checkerror(t, err)

	RefineGt(gt, sst, "mutation")
	fmt.Println(gt)
	fmt.Println(BinaryCost(gt, sst))
	gt1, _ := Make(gs)
	RefineGt(gt1, sst, "weighted", 3, 1)
	fmt.Println(gt1)
	fmt.Println(BinaryCost(gt1, sst))
}

func Test_Weight(t *testing.T) {
	niter := 1000
	crate := 0.8
	ntaxon := 40
	nleaves := 80

	testone := func() {
		gt := SimuTreeRandomTaxon(nleaves, ntaxon)
		gt.RandomContract(crate)
		//st := SimuTree(ntaxon)
		st := LineTree2(ntaxon)
		gs := gt.String()
		ss := st.String()
		//fmt.Println(":")
		//fmt.Println(gs)
		//fmt.Println(ss)

		gt1, err := Make(gs)
		if err != nil {
			t.Log(err)
		}
		st1, err := Make(ss)
		if err != nil {
			t.Log(err)
		}
		sst1, err := st1.SpeciesTree()
		checkerror(t, err)
		RefineGt(gt1, sst1, "mutation")
		//fmt.Println(gt1)
		d1, l1, dc1, err := BinaryCost(gt1, sst1)
		if err != nil {
			t.Log(err)
		}

		gt2, err := Make(gs)
		if err != nil {
			t.Log(err)
		}
		st2, err := Make(ss)
		if err != nil {
			t.Log(err)
		}
		sst2, err := st2.SpeciesTree()
		checkerror(t, err)

		//RefineGt(gt2, sst2, 3, 2001, 1000)
		//RefineGt(gt2, sst2, 3, 3000, 1)
		RefineGt(gt2, sst2, "weighted", 3, 1)
		d2, l2, dc2, err := BinaryCost(gt2, sst2)
		if err != nil {
			t.Log(err)
		}

		if d1+l1 != d2+l2 {
			//if d1 != d2 || l1 != l2 {
			fmt.Println(gs)
			fmt.Println(ss)
			fmt.Println(gt1)
			fmt.Println(gt2)
			fmt.Println(d1, l1, dc1)
			fmt.Println(d2, l2, dc2)
			t.Fatal("Not equal, fail!")
			//panic("Fail!!!")
		}
	}

	for i := 0; i < niter/100; i++ {
		print("~")
		for j := 0; j < 100; j++ {
			testone()
		}
	}
	println()
}

func test_Integer(t *testing.T) {
	niter := 1000
	crate := 0.8
	ntaxon := 50
	nleaves := 100

	testone := func() {
		gt := SimuTreeRandomTaxon(nleaves, ntaxon)
		gt.RandomContract(crate)
		st := SimuTree(ntaxon)
		gs := gt.String()
		ss := st.String()
		//fmt.Println(":")
		//fmt.Println(gs)
		//fmt.Println(ss)

		gt1, err := Make(gs)
		if err != nil {
			t.Log(err)
		}
		st1, err := Make(ss)
		if err != nil {
			t.Log(err)
		}
		sst1, err := st1.SpeciesTree()
		checkerror(t, err)
		RefineGt(gt1, sst1, "weighted", 5, 1)
		//fmt.Println(gt1)
		d1, l1, dc1, err := BinaryCost(gt1, sst1)
		if err != nil {
			t.Log(err)
		}

		gt2, err := Make(gs)
		if err != nil {
			t.Log(err)
		}
		st2, err := Make(ss)
		if err != nil {
			t.Log(err)
		}
		sst2, err := st2.SpeciesTree()
		checkerror(t, err)
		//RefineGt(gt2, sst2, 3, 2001, 1000)
		RefineGt(gt2, sst2, "weighted", 4, 1)
		d2, l2, dc2, err := BinaryCost(gt2, sst2)
		if err != nil {
			t.Log(err)
		}

		gt3, err := Make(gs)
		if err != nil {
			t.Log(err)
		}
		st3, err := Make(ss)
		if err != nil {
			t.Log(err)
		}
		sst3, err := st3.SpeciesTree()
		checkerror(t, err)
		//RefineGt(gt2, sst2, 3, 2001, 1000)
		RefineGt(gt3, sst3, "weighted", 4.1, 1)
		d3, l3, dc3, err := BinaryCost(gt3, sst3)
		if err != nil {
			t.Log(err)
		}

		f := func(a, b int) float64 {
			return float64(a)*2.1 + float64(b)
		}
		//if d1+l1 != d2+l2 {
		//if d1 != d2 || l1 != l2 {
		if f(d1, l1) != f(d3, l3) && f(d2, l2) != f(d3, l3) {
			fmt.Println(gs)
			fmt.Println(ss)
			fmt.Println(gt1)
			fmt.Println(gt2)
			fmt.Println(gt3)
			fmt.Println(d1, l1, dc1)
			fmt.Println(d2, l2, dc2)
			fmt.Println(d3, l3, dc3)
			//t.Log("Not equal, fail!")
			//panic("Fail!!!")
			t.Fatal("Not equal, fail!")
		}
	}

	for i := 0; i < niter; i++ {
		//print("~")
		testone()
	}
	println()
}
