package tree

import (
//	. "txtgo/tree"
	"errors"
	"sort"
)

func LinearRefineGt(gt *Tree, st *SpeciesTree) ([][]*Node, error) {
	if !st.IsBinary() {
		return nil, errors.New("Species tree is not binary.")
	}
	if gt.IsBinary() {
		return nil, nil
	}

	sl := st.Nodes
	gl := gt.Nodes

	lm, err := LcaMap(gt, st)
	lca := st.Lca
	if err != nil {
		return nil, err
	}
	M := lm.Map

	//--------------------------------------------------------------
	// 1) find the pre-image pre[s] of every species node s under
	// LCA mapping.
	pre := make([][]*Node, len(sl))

	// d[node.Id] : the index of node in node.Father.Children
	d := make([]int, len(sl))
	var sid int
	for i, gn := range gl {
		sid = M[i].Id
		pre[sid] = append(pre[sid], gn)
		if gn.IsInternal() {
			for j, c := range gn.Children {
				d[c.Id] = j
			}
		}
	}

	//--------------------------------------------------------------
	// 2) reorder the children of every non-binary gene tree node

	// swap the children i and j
	swap := func(children []*Node, i, j int) {
		children[i], children[j] = children[j], children[i]
	}


	ind := make([]int, len(gl))
	var fn *Node
	for _, sn := range sl {
		for _, gn := range pre[sn.Id] {
			fn = gn.Father
			if fn == nil {
				break
			}
			if !fn.IsBinary() {
				swap(fn.Children, ind[fn.Id], d[gn.Id])
				//fmt.Println("put", gn.Name, ind[fn.Id])
				ind[fn.Id]++
			}
		}
	}

	//fmt.Println(gt)

	//--------------------------------------------------------------
	// 3) for every non-binary gene tree node g
	// find every node s in I^*(g) and make a copy s' of s, then
	// set s'.Father = g. If g_i maps to s, then append g_i to s'.Ext.
	// Finally, append s' to b[s] (b[s.Id] in fact).
	// Notice that we just use s'.Father to store g.

	makeCopy := func(t *Node) *Node {
		n := &Node{&data{}, nil, nil}
		n.Name = t.Name
		n.Level = t.Level
		return n
	}

	appendCopy := func(t *Node, l []*Node, c *Node, leaf bool) []*Node{
		if len(l) == 0 || l[len(l)-1].Father != c.Father {
			n := makeCopy(t)
			n.Father = c.Father
			l = append(l, n)
		}
		if leaf {
			n := l[len(l)-1]
			if n.Ext != nil {
				n.Ext = append(n.Ext.([]*Node), c)				
			} else {
				n.Ext = []*Node{c}
			}
		}
		if l[len(l)-1].Father != c.Father {
			if c.Father == nil {
				panic("What!")
			}
		}
		return l
	}

	b := make([][]*Node, len(sl))
	var tn *Node
	var i, tid int

	for _, gn := range gl {

		if !gn.IsBinary() && !gn.IsLeaf() {

			for i = 0; i < len(gn.Children)-1; i++ {
				tn = M[gn.Children[i].Id]
				tid = tn.Id
				b[tid] = appendCopy(tn, b[tid], gn.Children[i], true)
				tn = lca(tn, M[gn.Children[i+1].Id])
				tid = tn.Id
				b[tid] = appendCopy(tn, b[tid], gn.Children[i], false)
			}
			tn = M[gn.Children[i].Id]
			tid = tn.Id
			b[tid] = appendCopy(tn, b[tid], gn.Children[i], true)
		}
	}
	//println("Step 3 done.")

	// 4) Using Euler Tour of species tree to find the corresponding
	// Euler Tour a[g] (a[g.Id] in fact) of every I^*(g).

	visit := make([]bool, len(sl))
	a := make([][]*Node, len(gl))
	var gn, sn *Node
	tn = st.Node
	for {
		tid = tn.Id
		visit[tid] = true

		for _, sn = range b[tid] {
			// In step 3), we use sn.Father = g
			// to indicates that sn is in the tree I^*(g).
			gn = sn.Father 
			
			//fmt.Println(sn.Name, "->", gn)
			a[gn.Id] = append(a[gn.Id], sn)
		}
		switch {
		case tn.IsLeaf():
			tn = tn.Father
		case !visit[tn.Children[0].Id]:
			tn = tn.Children[0]
		case !visit[tn.Children[1].Id]:
			tn = tn.Children[1]
		default:
			tn = tn.Father
		}

		if tn == nil {
			break
		}
	}
	//println("Step 4 done.")

	// 5) Finally, reconstruct I^*(g) using its Euler Tour a[g] (a[g.Id]).
	// Notice that I^*(g) is just the reconstructed tree with
	// root a[g][0] (a[g.Id][0]).

	// use Eulor Tour to reconstruct the sub-tree
	reconstruct := func(l []*Node) {
		addchild := func(f, n *Node) {
			if len(f.Children) == 0 || f.Children[len(f.Children)-1] != n {
				f.AddChild(n)
			}
		}
		var x, y *Node
		for i := 0; i < len(l)-1; i++ {
			x = l[i]
			y = l[i+1]
			switch {
			case x.Level < y.Level:
				addchild(x, y)
			case x.Level > y.Level:
				addchild(y, x)
			}
			// if x.Level == y.Level
			// and x y are consective in Euler Tour
			// thus x == y
		}

	}

	for i, gn = range gl {
		if !gn.IsBinary() && !gn.IsLeaf() {
			for _, n := range a[i] {
				n.Father = nil
			}
			reconstruct(a[i])

			linearRefineGeneNode(a[i][0])
			gn.replaceBy(a[i][0])
		}
	}

	gt.Update()
	// some clean up
	for _, n := range gt.Nodes {
		// remove node with out-degree 1.
		// just by replace node with its single children
		if len(n.Children) == 1 {
			n.replaceBy(n.Children[0])
		}

		if !n.IsLeaf() {
			n.Name = ""
		}
	}
	gt.Update()
	//println("Step 5 done.")
	//fmt.Println(gt)

	return a, nil
}

// Refine gene tree node to achieve minimal duplication + loss cost
func linearRefineGeneNode(root *Node) {
	// bottom up, update a, b on each edge with length d
	// i.e. d = node.Level - father.Level
	update := func(a, b, d int) (int, int) {
		switch {
		case d > 2:
			return 0, 0
		case d == 2:
			return 0, a
		case d == 1:
			return a, b
		default:
			panic("Incorrect level difference.")
		}
	}

	//bottom up, get a, b from w, a1, b1, a2, b2
	getab := func(w int, args ...int) (int, int) {
		sort.Ints(args)
		return w + args[1], w + args[2]
	}

	//top down, get k from t, a, b
	getk := func(a, b, t int) int {
		switch {
		case t >= b:
			return b
		case t <= a:
			return a
		default:
			return t
		}
	}

	// don't update level!
	// keep the level value from original species tree!

	nl := root.Post2List()
	size := len(nl)
	A := make([]int, size)
	B := make([]int, size)
	W := make([]int, size)

	// use P[i] to denote the preimages of node i,
	// which are stored in node.Ext
	P := make([][]*Node, size)

	// Not necessary, so remove it since Go's GC is slow. 
	// level difference, i.e. node.Level-father.Level
	//d := make([]int, size)

	var a1, b1, a2, b2, d1, d2, w, nchild int
	var node, father, lchild, rchild *Node
	var i int

	// Now calculate a, b, and w for each node.
	for i, node = range nl {
		node.Id = i
		if node.Ext != nil {
			P[i] = node.Ext.([]*Node)			
		} else {
			P[i] = nil
		}

		w = len(P[i])
		W[i] = w
		if node.IsLeaf() {
			A[i] = w - 1 // remove background tree
			B[i] = A[i]
		} else {
			nchild = len(node.Children)
			switch {
			case nchild > 2 || nchild == 0:
				//fmt.Println(node)
				panic("Incorrect children list.")
			case nchild == 2:
				rchild = node.Children[1]
				a2 = A[rchild.Id]
				b2 = B[rchild.Id]
				d2 = rchild.Level - node.Level
				a2, b2 = update(a2, b2, d2)
			case nchild == 1:
				a2, b2 = 0, 0
			}

			lchild = node.Children[0]
			a1 = A[lchild.Id]
			b1 = B[lchild.Id]
			d1 = lchild.Level - node.Level
			a1, b1 = update(a1, b1, d1)

			A[i], B[i] = getab(w, a1, b1, a2, b2)

		}
		//fmt.Println(A[i], B[i], W[i], a1, b1, d1, a2, b2, d2, node.Name)
	}

	K := make([]int, size)
	T := make([]int, size)
	// there is no extra lineage at root, thus t = 0
	T[size-1] = 0
	K[size-1] = getk(A[size-1], B[size-1], 0)

	//fmt.Println(K[size-1], nl[size-1].Name)

	// use Fid[i] to store the nl[i].Father.Id,
	// as nl[i].Father will change later.
	Fid := make([]int, size)

	// the number of incoming extra lineages at node
	var t int
	for i = size - 2; i >= 0; i-- {
		node = nl[i]
		father = node.Father
		Fid[i] = father.Id
		//fmt.Println(":", node.Name, father.Name)
		t = K[father.Id] - W[father.Id]
		// if level difference > 2, it's fine to loss all extra lineages
		if node.Level-father.Level > 2 {
			t = 0
		}
		T[i] = t
		K[i] = getk(A[i], B[i], t)
	}

	// insert nodes onto the edge (node.Father, node)
	insertNode := func(node *Node, nodes []*Node) {
		if len(nodes) == 0 {
			return
		}
		lchild := new(Node)
		//replaceNode(lchild, node)
		lchild.replaceBy(node)
		node.Children = nil
		for i := 0; i < len(nodes)-1; i++ {
			nnode := new(Node)
			nnode.AddChild(lchild)
			nnode.AddChild(nodes[i])
			lchild = nnode
		}
		node.AddChild(lchild)
		node.AddChild(nodes[len(nodes)-1])
	}
	//fmt.Println()

	// every node = nl[i] in the tree I^*(g) 
	// may has multiple extra copies,
	// we store them in the list nodes[i].
	nodes := make([][]*Node, len(K))

	for i := len(K) - 1; i >= 0; i-- {
		// K[i] extra copies.
		nodes[i] = make([]*Node, K[i])
		for j := range nodes[i] {
			nodes[i][j] = new(Node)
		}

		node = nl[i]
		//fmt.Println("in  ", root)
		//length := 0

		//fmt.Println(node.Name, length, K[i])

		// we need to replace the leave with original 
		// children of non-binary gene tree node, 
		// such that we can embed the refined tree into
		// the original gene tree.

		if W[i] > 0 {
			if node.IsLeaf() && K[i] != W[i]-1 {
				panic("Something wrong.")
			} else if !node.IsLeaf() && K[i] < W[i] {
				panic("Something wrong too.")
			}
			if node.IsLeaf() {
				for j, n := range nodes[i] {
					//replaceNode(n, P[i][j+1])
					n.replaceBy(P[i][j+1])
				}
				//replaceNode(node, P[i][0])
				node.replaceBy(P[i][0])
			} else {
				for j, pnode := range P[i] {
					//replaceNode(nodes[i][K[i]-1-j], pnode)
					nodes[i][K[i]-1-j].replaceBy(pnode)
				}
			}
		}

		if node.Father != nil {
			fid := Fid[i]
			// connect extra nodes at n to the nodes at n.Father
			// as much as possible
			//extra := nodes[fid][0 : K[fid]-W[fid]]

			// there are T[i] lineages inherited from n.father
			//extra := nodes[fid][0:T[i]]
			//length = len(extra)
			//fmt.Println(length, K[i])
			//for j:=0; j<length && j<K[i]; j++{
			//for j, n := range nodes[fid][0:T[i]] {
			for j := 0; j < T[i] && j < K[i]; j++ {
				nodes[fid][j].AddChild(nodes[i][j])
			}
		}
		// the rest nodes are just duplicated from background tree.
		// i.e. from the node = nl[i].
		if T[i] < K[i] {
			insertNode(node, nodes[i][T[i]:K[i]])
		}

		//fmt.Println("out ", root)
	}

	for _, n := range nl {
		if len(n.Children) == 1 {
			//replaceNode(n, n.Children[0])
			n.replaceBy(n.Children[0])
		}
	}

	// // now replace the original non-binary tree node with root of I^*(g)
	// replaceNode(root.Map, root)
}
