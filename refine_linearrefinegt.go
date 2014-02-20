package tree

import (
	"errors"
	"sort"
)

type InvalidMethodError struct {
	Method string
}

func (err InvalidMethodError) Error() string {
	return "Invalid method: " + err.Method
}

var (
	NotEnoughWeightsError     = errors.New("Not enough weight parameters.")
	NotBinarySpeciesTreeError = errors.New("Species tree is not binary.")
)

// A dispatch function
func RefineGt(gt *Tree, st *SpeciesTree, method string, weights ...float64) error {
	var refine func(*Node)
	switch method {
	case "duploss":
		refine = minDupThenLoss
	case "mutation":
		refine = minDupPlusLoss
	case "lossdup":
		refine = minLossThenDup

    // deprecated, use `affine` instead.
	case "weighted":
		if len(weights) < 2 {
			return NotEnoughWeightsError
		}
		refine = weightedCost(weights[0], weights[1])
	case "affine":
		if len(weights) < 2 {
			return NotEnoughWeightsError
		}
		refine = affineCost(weights[0], weights[1])
	default:
		return InvalidMethodError{method} // This shouldn't happen.
	}
	err := refineGt(gt, st, refine)
	if err != nil {
		return err
	}
	return nil
}

func refineGt(gt *Tree, st *SpeciesTree, refine func(*Node)) error {
	if !st.IsBinary() {
		return NotBinarySpeciesTreeError
	}
	if gt.IsBinary() {
		return nil
	}

	sl := st.Nodes
	gl := gt.Nodes

	lm, err := LcaMap(gt, st)
	lca := st.Lca
	if err != nil {
		return err
	}
	M := lm.Map

	//------------------------------------------------------
	// 1) find the pre-image `pre[s]` of every species
	// node s under LCA mapping.
	pre := make([][]*Node, len(sl))

	// `d[node.Id]`: index of node in `node.Father.Children`
	d := make([]int, len(gl))
	var sid int
	for i, gn := range gl {
		if gn.IsInternal() {
			for j, c := range gn.Children {
				d[c.Id] = j
			}
		}

		sid = M[i].Id
		pre[sid] = append(pre[sid], gn)
	}

	//------------------------------------------------------
	// 2) reorder the children of every non-binary
	// gene tree node by swaping the children i and j
	swap := func(children []*Node, i, j int) {
		children[i], children[j] = children[j], children[i]
		d[children[i].Id] = i
		d[children[j].Id] = j
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
				ind[fn.Id]++
			}
		}
	}

	//------------------------------------------------------
	// 3) for every non-binary gene tree node g
	// find every node `s` in `I^*(g)`, make a copy `s'` of `s`,
	// and set `s'.Father = g`.
	//
	// If `g_i` maps to `s`, then append `g_i` to `s'.Ext`.
	//
	// Finally, append `s'` to `b[s]` (`b[s.Id]` in fact).
	// Notice that we just use `s'.Father` to store `g`.

	makeCopy := func(t *Node) *Node {
		n := &Node{&data{}, nil, nil}
		n.Name = t.Name
		n.Level = t.Level
		return n
	}

	appendCopy := func(t *Node, l []*Node, c *Node, leaf bool) []*Node {
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

	// 4) Using Euler Tour of species tree to find
	// the corresponding Euler Tour `a[g]` (`a[g.Id]` in fact)
	// of every `I^*(g)`.
	visit := make([]bool, len(sl))
	a := make([][]*Node, len(gl))
	var gn, sn *Node
	tn = st.Node
	for {
		tid = tn.Id
		visit[tid] = true
		for _, sn = range b[tid] {
			// In step 3), we use `sn.Father = g`
			// to indicates that `sn` is in the tree `I^*(g)`.
			gn = sn.Father

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

	// 5) Finally, reconstruct `I^*(g)` using its Euler Tour
	// `a[g]` (`a[g.Id]`). Notice that `I^*(g)` is just the
	// reconstructed tree with root `a[g][0]` (`a[g.Id][0]`).

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
			// If `x.Level == y.Level`
			// and `x, y` are consective in Euler Tour
			// thus `x == y`
		}

	}

	for i, gn = range gl {
		if !gn.IsBinary() && !gn.IsLeaf() {
			for _, n := range a[i] {
				n.Father = nil
			}
			// frist, reconstruct `I^*(g)`
			reconstruct(a[i])
			// then, refine it
			refine(a[i][0])
			// finally, replace the original subtree by the refined one.
			gn.replaceBy(a[i][0])
		}
	}

	// some clean up
	gt.Update()
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
	return nil
}

// Refine gene tree node to achieve minimal duplication + loss cost
func minDupPlusLoss(root *Node) {
	// bottom up, update `a, b` on each edge with length `d`
	// i.e. `d = node.Level - father.Level`
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

	//bottom up, get `a, b` from `w, a1, b1, a2, b2`
	getab := func(w int, args ...int) (int, int) {
		sort.Ints(args)
		return w + args[1], w + args[2]
	}

	// don't update level!
	// keep the level value from original species tree!
	nl := root.Post2List()
	size := len(nl)
	A := make([]int, size)
	B := make([]int, size)
	W := make([]int, size)

	// use `P[i]` to denote the preimages of node `i`,
	// which are stored in `node.Ext`
	P := make([][]*Node, size)

	var a1, b1, a2, b2, d1, d2 int
	var node, father, lchild, rchild *Node
	var i, t int

	// Now calculate `a, b`, and `w` for each node.
	for i, node = range nl {
		node.Id = i
		if node.Ext != nil {
			P[i] = node.Ext.([]*Node)
			node.Ext = nil
		}

		W[i] = len(P[i])
		if node.IsLeaf() {
			A[i] = W[i] - 1 // remove background tree
			B[i] = A[i]
		} else {
			if len(node.Children) == 2 {
				rchild = node.Children[1]
				a2, b2 = A[rchild.Id], B[rchild.Id]
				d2 = rchild.Level - node.Level
				a2, b2 = update(a2, b2, d2)
			} else {
				a2, b2 = 0, 0
			}

			lchild = node.Children[0]
			a1, b1 = A[lchild.Id], B[lchild.Id]
			d1 = lchild.Level - node.Level
			a1, b1 = update(a1, b1, d1)

			A[i], B[i] = getab(W[i], a1, b1, a2, b2)
		}
	}

	// `K[i]`: the number of gene lineages entering the
	// incoming edge of `node[i]`.
	K := make([]int, size)

	// `T[i]`: the number of gene lineages leaving the
	// incoming branch of `node[i]`, and entering `node[i]`.
	// i.e. the number of gene copies that `node[i]` inherited.
	T := make([]int, size)

	// use `Fid[i]` to store the `nl[i].Father.Id`,
	// as `nl[i].Father` will change later.
	Fid := make([]int, size)

	// There is no extra lineage at root, thus `t = 0`
	T[size-1] = 0
	K[size-1] = project(0, A[size-1], B[size-1])
	// Now, compute the number of incoming and outgoing
	// extra lineages at each node.
	for i = size - 2; i >= 0; i-- {
		node = nl[i]
		father = node.Father
		Fid[i] = father.Id
		t = K[father.Id] - W[father.Id]

		// If level difference > 2, it's fine to
		// loss all extra lineages
		if node.Level-father.Level > 2 {
			t = 0
		}
		T[i] = t
		K[i] = project(t, A[i], B[i])
	}
	simpleConstruct(nl, K, T, W, Fid, P)
}

// Refine a gene tree node by minimizing
// duplication cost first, then gene loss cost.
func minDupThenLoss(root *Node) {
	nl := root.Post2List()
	size := len(nl)

	// `A[i] = W[i] - 1` for leaf, otherwise
	// `A[i]` is the min of `B[lchild]`, `B[rchild]`
	A := make([]int, size)
	// `B[i] = W[i] - 1` for leaf, otherwise
	// `B[i]` is the max of `B[lchild]`, `B[rchild]`
	B := make([]int, size)
	W := make([]int, size)
	P := make([][]*Node, size)
	for i, n := range nl {
		n.Id = i
		if n.Ext != nil {
			P[i] = n.Ext.([]*Node)
			n.Ext = nil
		}
		W[i] = len(P[i])
		nchild := len(n.Children)
		switch {
		case nchild == 0:
			A[i] = W[i] - 1
			B[i] = W[i] - 1
		case nchild == 1:
			A[i] = W[i] + 0
			B[i] = W[i] + B[n.Children[0].Id]
		case nchild == 2:
			b1 := B[n.Children[0].Id]
			b2 := B[n.Children[1].Id]
			if b1 > b2 {
				A[i] = b2
				B[i] = b1
			} else {
				A[i] = b1
				B[i] = b2
			}
			A[i] += W[i]
			B[i] += W[i]
		default:
			panic("Incorrect number of children.")
		}
	}
	Fid := make([]int, size)
	for i := 0; i < size-1; i++ {
		Fid[i] = nl[i].Father.Id
	}

	T := make([]int, size)
	K := make([]int, size)
	T[size-1] = 0
	K[size-1] = project(0, A[size-1], B[size-1])
	for i := size - 2; i >= 0; i-- {
		n := nl[i]
		T[i] = K[n.Father.Id] - W[n.Father.Id]
		K[i] = project(T[i], A[i], B[i])
	}

	simpleConstruct(nl, K, T, W, Fid, P)
}

// Refine a gene tree node by minimizing
// gene loss cost first, then duplication cost.
func minLossThenDup(root *Node) {
	nl := root.Post2List()
	size := len(nl)
	K := make([]int, size)
	W := make([]int, size)
	P := make([][]*Node, size)
	for i, n := range nl {
		n.Id = i
		if n.Ext != nil {
			P[i] = n.Ext.([]*Node)
			n.Ext = nil
		}
		W[i] = len(P[i])
		nchild := len(n.Children)
		switch nchild {
		case 0:
			K[i] = W[i] - 1
		case 1:
			K[i] = W[i]
		case 2:
			var v1, v2 int
			if n.Children[0].Level-n.Level == 1 {
				v1 = K[n.Children[0].Id]
			} // else v1 = 0
			if n.Children[1].Level-n.Level == 1 {
				v2 = K[n.Children[1].Id]
			} // else v2 = 0
			if v1 > v2 {
				K[i] = v2
			} else {
				K[i] = v1
			}
			K[i] += W[i]
		}
	}

	Fid := make([]int, size)
	for i := 0; i < size-1; i++ {
		Fid[i] = nl[i].Father.Id
	}

	T := make([]int, size)
	T[size-1] = 0
	for i := size - 2; i >= 0; i-- {
		n := nl[i]
		T[i] = K[n.Father.Id] - W[n.Father.Id]
	}

	simpleConstruct(nl, K, T, W, Fid, P)
}

// A simple reconstruction of subtree with
// the information on `I^*(g)`.
//
// All duplication occurs on the background tree.
func simpleConstruct(nl []*Node, K, T, W, Fid []int, P [][]*Node) {
	// insert nodes onto the edge `(node.Father, node)`
	insertNode := func(node *Node, nodes []*Node) {
		if len(nodes) == 0 {
			return
		}
		lchild := newNode()
		lchild.replaceBy(node)
		node.Children = nil
		for i := 0; i < len(nodes)-1; i++ {
			nnode := newNode()
			nnode.AddChild(lchild)
			nnode.AddChild(nodes[i])
			lchild = nnode
		}
		node.AddChild(lchild)
		node.AddChild(nodes[len(nodes)-1])
	}

	// every `node = nl[i]` in the tree `I^*(g)`
	// may has multiple extra copies,
	// we store them in the list `nodes[i]`.
	nodes := make([][]*Node, len(K))

	for i := len(K) - 1; i >= 0; i-- {
		// `K[i]` extra copies.
		nodes[i] = make([]*Node, K[i])
		for j := range nodes[i] {
			nodes[i][j] = newNode()
		}

		node := nl[i]

		// we need to replace the leaves with original
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
					n.replaceBy(P[i][j+1])
				}
				node.replaceBy(P[i][0])
			} else {
				for j, pnode := range P[i] {
					nodes[i][K[i]-1-j].replaceBy(pnode)
				}
			}
		}

		if node.Father != nil {
			fid := Fid[i]
			// Connect extra nodes at `n` to the nodes
			// at `n.Father` as much as possible.
			// There are `T[i]` lineages inherited from `n.father`.
			for j := 0; j < T[i] && j < K[i]; j++ {
				nodes[fid][j].AddChild(nodes[i][j])
			}
		}
		// the rest nodes are just duplicated from background tree.
		// i.e. from the `node = nl[i]`.
		if T[i] < K[i] {
			insertNode(node, nodes[i][T[i]:K[i]])
		}
	}

	for _, n := range nl {
		if len(n.Children) == 1 {
			n.replaceBy(n.Children[0])
		}
	}

	// This step has been moved to the caller function!

	// // now replace the original non-binary tree node with root of I^*(g)

	// replaceNode(root.Map, root)
}
