package tree

import (
	"fmt"
)

// A configuration consists of `x` and `slop`.
// where `len(slop) = len(x) + 1`.
// For example:
//
//      s[0]   s[1]   s[2] ...  s[k]   s[k+1]
//          x[0]   x[1]  x[2] ...   x[k]
type f struct {
	x    []int
	slop []float64
}

// Merge configurations.
func mergeF(fa, fb *f) *f {
	xa, sa := fa.x, fa.slop
	xb, sb := fb.x, fb.slop

	xc := make([]int, 0, len(xa)+len(xb))
	sc := make([]float64, 0, len(xa)+len(xb)+1)

	var i1, i2 int
	// Merge `x` and `slop` for each interval.
	for i1 < len(xa) || i2 < len(xb) {
		sc = append(sc, sa[i1]+sb[i2])
		switch {
		case i1 == len(xa):
			xc = append(xc, xb[i2])
			i2++
		case i2 == len(xb):
			xc = append(xc, xa[i1])
			i1++
		case xa[i1] < xb[i2]:
			xc = append(xc, xa[i1])
			i1++
		case xa[i1] > xb[i2]:
			xc = append(xc, xb[i2])
			i2++
		default:
			xc = append(xc, xa[i1])
			i1++
			i2++
		}
	}
	// Append the last slop.
	sc = append(sc, sa[i1]+sb[i2])
	return &f{xc, sc}
}

// Find `x^+` and `x^-`.
func getX(c *f, wdup, wloss float64) (a, b int) {
	var i int
	s := c.slop
	for i = 0; i < len(s); i++ {
		if s[i] > -wdup {
			a = i
			break
		}
	}
	for ; i < len(s); i++ {
		if s[i] >= wloss {
			b = i
			break
		}
	}
	return
}

func getJK(c *f, wdup, wloss float64) (j, k int) {
	a, b := getX(c, wdup, wloss)
	return c.x[a-1], c.x[b-1]
}

// Restrict the configuration.
func restrictF(c *f, wdup, wloss float64) *f {
	a, b := getX(c, wdup, wloss)

	// If cannot find such `j` and `k`,
	// just panic!
	// This shouldn't happen.
	if a < 0 || b < 0 {
		fmt.Println(c, wdup, wloss)
		fmt.Println("j k:", a, b)
		panic("Cannot restrict!")
	}

	// Restrict `x` and `slop`.
	c.x = c.x[a-1 : b]
	c.slop = c.slop[a-1 : b+1]
	c.slop[0] = -wdup
	c.slop[len(c.slop)-1] = wloss
	return c
}

// Refine gene tree node by minimizing
// the affine cost defined by `wdup`, `wloss`.
func affineCost(wdup, wloss float64) func(*Node) {
	return func(root *Node) {
		// Make configuration for leaf.
		// Put this inside as it needs `wdup` and `wloss`.
		makeLeafF := func(w int) *f {
			if w == 0 {
				panic("Incorrect Leaf.")
			}
			return &f{[]int{w - 1}, []float64{-wdup, wloss}}
		}

		// Lift configuration along a path with length `d`.
		liftF := func(c *f, d int) *f {
			if d <= 1 {
				return c
			}
			// `(d-1)` gene losses on the path.
			for i, _ := range c.slop {
				c.slop[i] += float64(d-1) * wloss
			}
			// If the smallest `x` is greater than `0`,
			// we bound the configuration below by `0`.
			if c.x[0] != 0 {
				c.x = append([]int{0}, c.x...)
				c.slop = append([]float64{-wdup}, c.slop...)
			} else {
				c.slop[0] = -wdup
			}
			return restrictF(c, wdup, wloss)
		}

		// Shift configuration if `W[i]>0`.
		shiftF := func(c *f, w int) *f {
			// If shift is `0`, just use the original configuration.
			if w == 0 {
				return c
			}
			// Otherwise, make a new configuration.
			nc := new(f)
			nc.x = make([]int, len(c.x))
			nc.slop = make([]float64, len(c.slop))
			// Copy the original one.
			copy(nc.x, c.x)
			copy(nc.slop, c.slop)
			// Increase each `x` by `w`.
			for i, _ := range nc.x {
				nc.x[i] += w
			}
			return nc
		}

		nl := root.Post2List()
		// Number of nodes in the gene tree.
		size := len(nl)
		// List of configuration.
		F := make([]*f, size)
		// List of the pre-images
		P := make([][]*Node, size)
		// List of the size of pre-images.
		W := make([]int, size)
		// List of bounds `(x^+, x^-)`.
		A := make([][2]int, size)

		for i, n := range nl {
			n.Id = i
			if n.Ext != nil {
				P[i] = n.Ext.([]*Node)
				n.Ext = nil
			}
			W[i] = len(P[i])
			switch len(n.Children) {
			case 0:
				F[i] = makeLeafF(W[i])
			case 1:
				lc := n.Children[0]
				// One more loss!
				F[i] = shiftF(liftF(F[lc.Id], lc.Level-n.Level+1), W[i])
			case 2:
				lc := n.Children[0]
				fa := liftF(F[lc.Id], lc.Level-n.Level)

				rc := n.Children[1]
				fb := liftF(F[rc.Id], rc.Level-n.Level)

				F[i] = shiftF(restrictF(mergeF(fa, fb), wdup, wloss), W[i])

			default:
				panic("Not a binary node!")
			}
			A[i][0], A[i][1] = getJK(F[i], wdup, wloss)
		}

		Fid := make([]int, size)
		for i := 0; i < size-1; i++ {
			Fid[i] = nl[i].Father.Id
		}

		// `T[i]` : number of genes flow into a branch (or a path) `i`.
		T := make([]int, size)
		// `K[i]` : number of genes flow out of branch (or a path) `i`.
		K := make([]int, size)
		T[size-1] = 0
		K[size-1] = project(0, A[size-1][0], A[size-1][1])
		for i := size - 2; i >= 0; i-- {
			n := nl[i]
			d := n.Level - n.Father.Level

			T[i] = K[Fid[i]] - W[Fid[i]]
			// Check how many genes flow into a long path.
			if d > 1 {
				j, k := getJK(F[i], wdup, wloss)
				T[i] = project(T[i], j, k)
			}
			K[i] = project(T[i], A[i][0], A[i][1])
		}
		// Now reconstruct a refined gene tree.
		simpleConstruct(nl, K, T, W, Fid, P)
	}
}

// ........................................................................

// DP + C
// Now should be replaced by Affine method, which is faster in extreme case.
func weightedCost(wdup, wdc float64) func(*Node) {
	return func(root *Node) {
		nl := root.Post2List()
		size := len(nl)
		M := make([]int, size)
		W := make([]int, size)
		P := make([][]*Node, size)
		U := make([][]float64, size)
		I := make([][]int, size)
		for i, n := range nl {
			n.Id = i
			if n.Ext != nil {
				P[i] = n.Ext.([]*Node)
				n.Ext = nil
			}
			W[i] = len(P[i])
			switch {
			case len(n.Children) == 0:
				M[i] = W[i] - 1
			case len(n.Children) == 1:
				M[i] = M[n.Children[0].Id] + W[i]
			case len(n.Children) == 2:
				v1 := M[n.Children[0].Id]
				v2 := M[n.Children[1].Id]
				if v1 > v2 {
					M[i] = v1
				} else {
					M[i] = v2
				}
				M[i] += W[i]
			default:
				panic("Not a binary node!")
			}
		}

		// A special case is that there is
		// only one node in I^*(g).
		if size == 1 {
			K := make([]int, size)
			T := make([]int, size)
			T[0] = 0
			K[0] = W[0] - 1
			simpleConstruct(nl, K, T, W, nil, P)
			return
		}

		mul := M[size-1]
		C := make([]float64, mul+1)

		min := func(args ...float64) float64 {
			m := args[0]
			for i := 1; i < len(args); i++ {
				if args[i] < m {
					m = args[i]
				}
			}
			return m
		}

		getCost := func(in, out, d int) float64 {
			m := min(float64(d)*wdc, wdup)
			if in >= out {
				return m * float64(out)
			} else {
				return m*float64(in) + wdup*float64(out-in)
			}
		}

		getMin := func(ind, in, d int) {
			// multiplicity
			w := W[ind]
			// list of cost
			l := C
			// (p, m): the cost m with p outgoing lineages
			// in branch ind with in incoming lineages.
			p, m := w, getCost(in, w, d)+l[w]
			for out := w + 1; out <= M[ind]; out++ {
				t := getCost(in, out, d) + l[out]
				if t < m {
					p, m = out, t
				}
			}
			U[ind][in] = m
			I[ind][in] = p
		}

		getC := func(i, w int) {
			n := nl[i]
			switch {
			case len(n.Children) == 1:
				lu := U[n.Children[0].Id]
				for j := w; j <= M[i]; j++ {
					C[j] = lu[j-w]
				}
			case len(n.Children) == 2:
				lu := U[n.Children[0].Id]
				ru := U[n.Children[1].Id]
				for j := w; j <= M[i]; j++ {
					C[j] = lu[j-w] + ru[j-w]
				}
			}
		}

		Fid := make([]int, size)
		for i := 0; i < size-1; i++ {
			Fid[i] = nl[i].Father.Id
		}

		for i, n := range nl {
			w := W[i]
			switch {
			case n.IsLeaf():
				U[i] = make([]float64, M[Fid[i]]-W[Fid[i]]+1)
				I[i] = make([]int, M[Fid[i]]-W[Fid[i]]+1)
				d := n.Level - n.Father.Level
				for j := 0; j <= M[Fid[i]]-W[Fid[i]]; j++ {
					U[i][j] = getCost(j, w-1, d)
					I[i][j] = w - 1
				}
			case n.IsRoot():
				U[i] = make([]float64, 1)
				I[i] = make([]int, 1)
				getC(i, w)
				getMin(i, 0, 0)
			default: // internal but not root
				U[i] = make([]float64, M[Fid[i]]-W[Fid[i]]+1)
				I[i] = make([]int, M[Fid[i]]-W[Fid[i]]+1)
				getC(i, w)
				d := n.Level - n.Father.Level
				for j := 0; j <= M[Fid[i]]-W[Fid[i]]; j++ {
					getMin(i, j, d)
				}
			}
		}

		K := make([]int, size)
		T := make([]int, size)
		T[size-1] = 0
		K[size-1] = I[size-1][0]
		for i := size - 2; i >= 0; i-- {
			n := nl[i]
			T[i] = K[n.Father.Id] - W[n.Father.Id]

			// Important!
			d := n.Level - n.Father.Level
			if float64(d)*wdc > wdup {
				T[i] = 0
			}
			K[i] = I[i][T[i]]
		}

		simpleConstruct(nl, K, T, W, Fid, P)
	}
}


