package tree

import (
	"fmt"
	"testing"
)

func Test_restrict(t *testing.T) {
	f1 := &f{[]int{1, 2, 3}, []float64{-3, -1, 1, 2}}
	fmt.Println(f1)

	restrictF(f1, 3, 2)
	fmt.Println(f1)

	f2 := &f{[]int{2}, []float64{-3, 2}}
	fmt.Println(f2)

	restrictF(f2, 3, 2)
	fmt.Println(f2)
}

func Test_MergeAndRestrict(t *testing.T) {
	f1 := &f{[]int{1, 3, 5}, []float64{-3, -1, 1, 2}}
	f2 := &f{[]int{2, 4, 5}, []float64{-1, 0, 1, 2}}
	fmt.Println(f1, f2)

	f := mergeF(f1, f2)
	fmt.Println(f)
}

func Test_affineCost(t *testing.T) {
	wdup := 3.0
	wloss := 1.0
	gs := "(1,1,7,7,8,8)"
	ss := "(((((((1,2),3),4),5),6),7),8);"

	gs = "((3,3,0),0,1)"
	ss = "(((0,1),2),3)"

	gs = "(((1,2),2),(2,3))"
	ss = "(((0,1),2),3)"

	gs = "((2,1),2,2,3,3)"
	ss = "(((0,1),2),3)"

	gs = "(((7,1,5),2),0,0,7,6,(0,6,3),3,7,6,7)"
	//gs = "(0,0,7,6,(0,6,3),3,7,6,7)"
	//gs = "(0,0,6,3,7,6,7)"
	ss = "(((((((0,1),2),3),4),5),6),7)"

	//gs = "(1,1,2,2,3,3,4,4,5,5,6,6,7,7,8,8,9,9,10,10,11,11,12,12)"
	//ss = "((((((((((((0,1),2),3),4),5),6),7),8),9),10),11),12)"

	//gs = "((1,2),1,1)"
	//ss = "((0,1),2)"

	st, _ := Make(ss)
	sst := st.SpeciesTree()

	gt, _ := Make(gs)
	//affine := affineCost(wdup, wloss)
	//linearRefineGt(gt, sst, affine)
	RefineGt(gt, sst, "affine", wdup, wloss)
	fmt.Println(gt)
	fmt.Println(BinaryCost(gt, sst))

	gt1, _ := Make(gs)
	RefineGt(gt1, sst, "weighted", 2*wloss+wdup, wloss)
	fmt.Println(gt1)
	fmt.Println(BinaryCost(gt1, sst))
}

func Test_AffineCorrect(t *testing.T) {
	niter := 10000
	crate := 0.7
	ntaxon := 50
	nleaves := 100

	wdup := 7.0
	wloss := 1.5

	testone := func() {
		//fmt.Println()
		gt := SimuTreeRandomTaxon(nleaves, ntaxon)
		gt.RandomContract(crate)
		//st := SimuTree(ntaxon)
		st := LineTree2(ntaxon)
		gs := gt.String()
		ss := st.String()

		gt1, err := Make(gs)
		if err != nil {
			t.Log(err)
		}
		st1, err := Make(ss)
		if err != nil {
			t.Log(err)
		}
		sst1 := st1.SpeciesTree()
		RefineGt(gt1, sst1, "affine", wdup, wloss)
		//fmt.Println(gt1)
		d1, l1, dc1, err := BinaryCost(gt1, sst1)
		if err != nil {
			fmt.Println()
			fmt.Println(gs)
			fmt.Println()
			fmt.Println(ss)
			fmt.Println()
			fmt.Println(gt1)
			t.Fatal(err)
		}

		gt2, err := Make(gs)
		if err != nil {
			t.Log(err)
		}
		st2, err := Make(ss)
		if err != nil {
			t.Log(err)
		}
		sst2 := st2.SpeciesTree()
		RefineGt(gt2, sst2, "weighted", 2*wloss+wdup, wloss)
		d2, l2, dc2, err := BinaryCost(gt2, sst2)
		if err != nil {
			fmt.Println(gt2)
			t.Fatal(err)
		}

		if wdup*float64(d1)+wloss*float64(l1) != wdup*float64(d2)+wloss*float64(l2) {
			fmt.Println()
			fmt.Println(gs)
			fmt.Println()
			fmt.Println(ss)
			fmt.Println()
			fmt.Println(gt1)
			fmt.Println()
			fmt.Println(gt2)
			fmt.Println()
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
