package main

import (
	"flag"
	"fmt"
	"log"
	"io/ioutil"
	"os"
	T "txtgo/tree"
)

const version = "20140210"

var (
	gf          = flag.String("g", "", "gene tree file")
	sf          = flag.String("s", "", "species tree file")
	m           = flag.String("m", "mutation", "method = mutation/duploss/lossdup/affine")
	wdup        = flag.Float64("wdup", 1.0, "weight of duplication, if method=affine")
	wloss       = flag.Float64("wloss", 1.0, "weight of loss, if method=affine")
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
    -m      method [mutation/duploss/lossdup/weighted/affine]
    -wdup   weight for duplication, only for weighted/affine
    -wloss  weight for loss, only for weighted/affine
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
		log.Fatal(err)
	}
	return string(data)
}

func checkerror(err error) {
    if err != nil {
        log.Fatal(err)
    }
}

func main() {
	flag.Usage = usage
	flag.Parse()
	if *showversion {
		fmt.Println("Version:", version)
		return
	}

	ss := read(*sf)
	gs := read(*gf)
	_wdup := *wdup
	_wloss := *wloss
	var _m int
	switch *m {
	case "mutation":
		_m = 1
	case "duploss":
		_m = 0
	case "lossdup":
		_m = 2
	case "weighted":
		_m = 3
		_wdup = *wdup + 2*(*wloss)
	case "affine":
		_m = 4
	default:
		log.Fatal("Invalid method.")
	}

	if *wdup <= 0 {
		log.Fatal("wdup must be non-negative.")
	}

	if *wloss <= 0 {
		log.Fatal("wdup must be non-negative.")
	}

	gt, err := T.Make(gs)
    checkerror(err)
		st, err := T.Make(ss)
        checkerror(err)
	sst := st.SpeciesTree()
	if !sst.IsBinary() {
		log.Fatal("Species must be binary.")
	}

    err = T.RefineGt(gt, sst, _m, _wdup, _wloss)
    checkerror(err)

	dup, loss, _, err := binaryCost(gt, sst)
    checkerror(err)

	fmt.Println("Refined Gene Tree:")
	fmt.Println(gt)
	fmt.Printf("\nDup: %d, Loss: %d\n", dup, loss)
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
			ld := (l.Level + r.Level - 2*m.Level)
			loss += ld
			dc += ld
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
