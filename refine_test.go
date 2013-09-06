package tree

import (
	"fmt"
	"testing"
)

func Test_Lineartimegt(t *testing.T) {
	gs := "((a,b,c)A,(d,e)D)R;"
	ss := "((a,b)A,(c,(d,e)D)C)R;"
	gt, err  := Make(gs)
	if err != nil {
		t.Log(err)
	}
	st, err := Make(ss)
	if err != nil {
		t.Log(err)
	}
	sst := st.SpeciesTree()
	LinearRefineGt(gt, sst)
	fmt.Println(gt)
}
