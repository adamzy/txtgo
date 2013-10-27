package main

import (
	"flag"
	"fmt"
	"log"
	T "txtgo/tree"
	//"strings"
	"io/ioutil"
    "os"
)

const version = "201310027"

var (
	gf     = flag.String("g", "", "gene tree file")
	sf     = flag.String("s", "", "species tree file")
	m = flag.String("m", "mutation", "method = mutation/duploss/lossdup/affine")
	wdup   = flag.Float64("wdup", 1.0, "weight of duplication, if method=affine")
	wloss  = flag.Float64("wloss", 1.0, "weight of loss, if method=affine")
    showversion = flag.Bool("V", false, "show version")
)

func usage() {
    s := fmt.Sprintf(`TxT-RGT
    version %s
Usage:
    %s -g gene_tree -s species_tree -m method [-wdup a -wloss b]
Options:
    -g      gene tree file    
    -s      species tree file
    -m      method [mutation/duploss/lossdup/affine]
    -wdup   weight for duplication, only for affine
    -wloss  weight for loss, only for affine
    -V      show version
Example:
    %s -g gene_tree -s species_tree -m mutation
    %s -g gene_tree -s species_tree -m affine -wdup 2 -wloss 1
`, version, os.Args[0], os.Args[0], os.Args[0])
    fmt.Print(s)
}

func read(fname string) string {
    if fname == "" {
        flag.Usage()
        os.Exit(1)
    }
	data, err := ioutil.ReadFile(fname)
	if err != nil {
		//println("Cannot open", fname)
		log.Fatal(err)
	}
	return string(data)
}

//func usage() {
    //fmt.Fprintf(os.Stderr, "TxT-RGT\n")
    //fmt.Fprintf(os.Stderr, "usage:\n")
    //flag.PrintDefaults()
    //os.Exit(2)
//}


func main() {
    flag.Usage = usage
	flag.Parse()
    if *showversion {
        fmt.Println("Version:", version)
        return
    }

	ss := read(*sf)
	gs := read(*gf)

	var _m int
	switch *m {
	case "mutation":
		_m = 1
	case "duploss":
		_m = 0
	case "lossdup":
		_m = 2
	case "affine":
		_m = 3
	default:
		log.Fatal("Invalid method.")
	}

	if *wdup <= 0 {
		log.Fatal("wdup must be non-negative.")
	}

	if *wloss <= 0 {
		log.Fatal("wdup must be non-negative.")
	}

	_wdup := *wdup + 2*(*wloss)
	_wdc := *wloss

	gt, err := T.Make(gs)
	if err != nil {
		log.Fatal(err)
	}
	st, err := T.Make(ss)
	if err != nil {
		log.Fatal(err)
	}
	//st.PruneFromTree(gt)
	sst := st.SpeciesTree()
	//_, err = T.LinearRefineGt(gt, st)
    if !sst.IsBinary() {
        log.Fatal("Species must be binary.")
    }
	T.RefineGt(gt, sst, _m, _wdup, _wdc)
	dup, loss, _, err := binaryCost(gt, sst)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Refined Gene Tree:")
	fmt.Println(gt)
	fmt.Printf("\nDup: %d, Loss: %d\n", dup, loss)
	//fmt.Print(dup, loss, dc, err)
	//fmt.Println(method, wdup, wdc)
}

func binaryCost(gt *T.Tree, st *T.SpeciesTree) (dup, loss, dc int, err error) {
	lm, err := T.LcaMap(gt, st)
    
	if err != nil {
		return 0, 0, 0, err
	}

	for _, n := range gt.Nodes {
		if !n.IsLeaf() {
			m := lm.Map[n.Id]
			l := lm.Map[n.Children[0].Id]
			r := lm.Map[n.Children[1].Id]
			loss += (l.Level + r.Level - 2*m.Level)
			dc += (l.Level + r.Level - 2*m.Level)
			if m == l || m == r {
				dup++
                n.Name = n.Name + "*"
			} else {
				loss -= 2
			}
		}
	}
	dc -= len(st.Nodes) - 1
	return
}
