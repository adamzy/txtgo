package tree

import (
	//"errors"
	"fmt"
)

// A configuration consists of x and slop.
// where len(slop) = len(x) + 1.
// For example:
// s[0]   s[1]   s[2] ...  s[k]   s[k+1]
//    x[0]   x[1]  x[2] ...   x[k]
type f struct {
	x    []int
	slop []float64
}

// Merge configurations from children.
func mergeF(fa, fb *f) *f {
	//fa = restrictF(fa, wdup, wloss)
	//fb = restrictF(fb, wdup, wloss)
	xa, sa := fa.x, fa.slop
	xb, sb := fb.x, fb.slop

	xc := make([]int, 0, len(xa)+len(xb))
	sc := make([]float64, 0, len(xa)+len(xb)+1)

	var i1, i2 int
	//x1, x2 := xa[i1], xb[i2]
	//s1, s2 := sa[i1], sb[i2]
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
	// append the last slop
	sc = append(sc, sa[i1]+sb[i2])
	return &f{xc, sc}
}

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

	// If cannot find such j and k,
	// just panic!
	// This shouldn't happen.
	if a < 0 || b < 0 {
		fmt.Println(c, wdup, wloss)
		fmt.Println("j k:", a, b)
		panic("Cannot restrict!")
	}

	// Restrict x and slop.
	c.x = c.x[a-1 : b]
	c.slop = c.slop[a-1 : b+1]
	c.slop[0] = -wdup
	c.slop[len(c.slop)-1] = wloss
	return c
}

// Refine gene tree W.R.T. the affine cost defined by wdup, wloss
func affineCost(wdup, wloss float64) func(*Node) {
	return func(root *Node) {
		// Make configuration for leaf.
		// Put this inside as it needs wdup and wloss.
		makeLeafF := func(w int) *f {
			if w == 0 {
				panic("Incorrect Leaf.")
			}
			return &f{[]int{w - 1}, []float64{-wdup, wloss}}
		}

		// Lift configuration along a path with length d.
		liftF := func(c *f, d int) *f {
			if d <= 1 {
				return c
			}
			// (d-1) gene losses on the path.
			for i, _ := range c.slop {
				c.slop[i] += float64(d-1) * wloss
			}
			// If the smallest x is greater than 0,
			// we bound the configuration below by 0.
			if c.x[0] != 0 {
				c.x = append([]int{0}, c.x...)
				c.slop = append([]float64{-wdup}, c.slop...)
			} else {
				c.slop[0] = -wdup
			}
			return restrictF(c, wdup, wloss)
		}

		// Shift configuration if W[i]>0.
		shiftF := func(c *f, w int) *f {
			// If shift is 0, just use the original configuration.
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
			// Increase each x by w.
			for i, _ := range nc.x {
				nc.x[i] += w
			}
			return nc
		}

		// Project in onto interval [a,b].
		project := func(i, a, b int) int {
			switch {
			case i <= a:
				return a
			case i >= b:
				return b
			default:
				return i
			}
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
		// List of bounds (x+, x-).
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
				//fmt.Println("leaf:", i, F[i])
			case 1:
				lc := n.Children[0]
				// One more loss!
				F[i] = shiftF(liftF(F[lc.Id], lc.Level-n.Level+1), W[i])
				//fmt.Println("inte:", i, F[i])
			case 2:
				lc := n.Children[0]
				fa := liftF(F[lc.Id], lc.Level-n.Level)
				//fmt.Println(" child:", lc.Id, fa)

				rc := n.Children[1]
				fb := liftF(F[rc.Id], rc.Level-n.Level)
				//fmt.Println(" child:", rc.Id, fb)

				F[i] = shiftF(restrictF(mergeF(fa, fb), wdup, wloss), W[i])

				//fmt.Println("inte:", i, F[i])

			default:
				panic("Not a binary node!")
			}
			A[i][0], A[i][1] = getJK(F[i], wdup, wloss)
		}
		//fmt.Println("F done")

		Fid := make([]int, size)
		for i := 0; i < size-1; i++ {
			Fid[i] = nl[i].Father.Id
		}

		// T[i] : number of genes flow into a branch (or a path) i.
		T := make([]int, size)
		// K[i] : number of genes flow out of branch (or a path) i.
		K := make([]int, size)
		T[size-1] = 0
		K[size-1] = project(0, A[size-1][0], A[size-1][1])
		//fmt.Println("i, t, k:", size-1, T[size-1], K[size-1])
		for i := size - 2; i >= 0; i-- {
			//fmt.Println(i, F[i])
			n := nl[i]
			d := n.Level - n.Father.Level

			T[i] = K[Fid[i]] - W[Fid[i]]
			// Check how many genes flow into a long path.
			if d > 1 {
				j, k := getJK(F[i], wdup, wloss)
				T[i] = project(T[i], j, k)
			}
			K[i] = project(T[i], A[i][0], A[i][1])
			//fmt.Println("i, t, k:", i, T[i], K[i])
		}
		// Now reconstruct a refined gene tree.
		simpleConstruct(nl, K, T, W, Fid, P)
	}
}
