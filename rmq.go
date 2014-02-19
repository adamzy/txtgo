// Package `rmq` implements Bender-Farach Algorithm 
// for RMQ (Range Minimum Query) problem.
// Assume size of input array <= 2^63-1, since we use int64.
// This should be OK for almost all cases.
// For the algorithm details, check 
//   Bender, Michael A., and Martin Farach-Colton. 
//   "The LCA problem revisited." LATIN 2000: Theoretical Informatics. 
//   Springer Berlin Heidelberg, 2000. 88-94.
package rmq

import "math"

func make2dslice(row, col int64) [][]int64 {
    S := make([][]int64, row)
    for i := range S {
        S[i] = make([]int64, col)
    }
    return S
}

// Sparse Table algorithm (time complexity O(n*log(n))
// Given a list `A`, `st(A)` return a function `f`,
// where `f(x,y)` return the position and value of the
// item with minimal value between location `x` and `y`.
func st(A []int64) func(x, y int64) (loci, value int64) {
    size := int64(len(A))
    _log := log2(size)
    s := _log[size]
    _pow := power2(s)
    
    M := make2dslice(size, s+1)
    N := make2dslice(size, s+1)

    for i:=int64(0); i< size; i++ {
        M[i][0] = A[i]
        N[i][0] = i
    }
    for j:=int64(1); j<=s; j++ {
        for i:=int64(0); i<=size-_pow[j]; i++ {
            k := i+_pow[j-1]
            if M[i][j-1] <= M[k][j-1] {
                M[i][j] = M[i][j-1]
                N[i][j] = N[i][j-1]
            } else {
                M[i][j] = M[k][j-1]
                N[i][j] = N[k][j-1]
            }
        }
    }

    return func (x, y int64) (posi, value int64) {
        if x==y {
            return x, A[x]
        }
        if x > y {
            x,y = y,x
        }
        r := _log[y-x] 
        k := y-_pow[r]+1
        if M[x][r] <= M[k][r] {
            return N[x][r], M[x][r]
        } else {
            return N[k][r], M[k][r]
        }
    }
}

func power2(size int64) []int64 {
	pow := make([]int64, size+1)
	p := int64(1)
	for i := range pow {
		pow[i] = p
		p <<= 1
	}
	return pow
}

// log2(size) return a list log[] which length size+1,
// where log[x] return n s.t. 2^n <= x < 2^(n+1),
// i.e. log[x] = floor(log_2(x)) for x > 0.
func log2(size int64) []int64 {
	log := make([]int64, size+1)
    var i,j,k int64
    i, j = 1, 2
	for j <= size {
		for i < j {
			log[i] = k
			i++
		}
		j <<= 1
		k++
	}
	for i <= size {
		log[i] = k
		i++
	}
	return log
}

// +-RMQ, the restricted RMQ
func ResRMQ(A []int64) func(x, y int64) (p, v int64) {
	length := int64(len(A))
    // if `A` has length less than 2, do nothing.
	if length < 2 {
		return func(x, y int64) (int64, int64) {
			return 0, 0
		}
    }

    size := int64(math.Ceil(math.Log2(float64(2*length))))

	// num is the number of blocks
	num := length / size
	if length%size != 0 {
		num++
	}
	// loc[i] is the block in which i is.
	loc := make([]int64, length)
	// res[i] is the offset of i in block loc[i].
	res := make([]int64, length)

	// pre[i] is the position of the first element in block i.
	pre := make([]int64, num)

	// B[i][j][k] is the optimal position in block i,
    // between position j and  k.
	B := make([][][]int64, num)

	// for sparse table algorithm
	C := make([]int64, num) // posi
	D := make([]int64, num) // value

	var i, j int64
	processer := processBlock(size)

	var pos int64
	for pos = 0; pos+size < length; pos += size {
		for j = 0; j < size; j++ {
			loc[pos+j] = i
			res[pos+j] = j
		}
		pre[i] = pos
		B[i] = processer(A[pos : pos+size])
		C[i] = B[i][0][size-1] + pos
		D[i] = A[C[i]]
		i++
	}
	if pos != length {
		// in this case, pos < length and i = num-1
		for j = 0; j < length-pos; j++ {
			loc[pos+j] = i
			res[pos+j] = j
		}
		pre[i] = pos
		B[i] = processer(A[pos:length])
		C[i] = B[num-1][0][length-pos-1] + pos
		D[i] = A[C[num-1]]
	}
	tabler := st(D)

	return func(x, y int64) (int64, int64) {
		// always assume 0 <= x <= y < length
        if x > y {
            x, y = y, x
        }
		u, v := loc[x], loc[y]
		ur, vr := res[x], res[y]
		up, vp := pre[u], pre[v]

		if u == v {
			pos := up + B[u][ur][vr]
			return pos, A[pos]
		} else {
			// u < v
			p1 := up + B[u][ur][size-1]
			v1 := A[p1]
			p2 := vp + B[v][0][vr]
			v2 := A[p2]

			var vv, pp int64
			if v1 <= v2 {
				vv = v1
				pp = p1
			} else {
				vv = v2
				pp = p2
			}

			if v-u > 1 {
				p3, v3 := tabler(u+1, v-1)
				if v3 < vv {
					vv = v3
					pp = C[p3]
				}
			}
			return pp, vv
		}
	}
}

// Convert a +-1 sequence to int, e.g.
// [1 2 3 2 1 0 -1 -2 -1 0 1] => 0-1 sequence 1100000111 => integer
func arraytoint(A []int64) int64 {
	if len(A) > 63 {
		panic("We cannot handle such long sequence")
	}
	var v int64
	for i := len(A) - 1; i > 0; i-- {
		if A[i] == A[i-1]+1 {
			v |= 1
		} else if A[i] != A[i-1]-1 {
			panic("Not a +-1 sequence")
		}
		v <<= 1
	}
	return v
}

func processBlock(size int64) func([]int64) [][]int64 {
    // number of all possible +-1 sequence begin with 0
	num := int64(1 << uint64(size))
	a := make([][][]int64, num)

	return func(A []int64) [][]int64 {
		v := arraytoint(A)
        // if a[v] hasn't been computed before, then compute it.
        // otherwise, use the cached result.
        if a[v] == nil {
            a[v] = dynamic(v, size)
        }
		return a[v]
	}
}

func dynamic(n, size int64) [][]int64 {
	if size > 63 {
		panic("We cannot handle such long sequence")
	}
	t := make([][]int64, size) //t[i][j] = position with min value in A[i,j]
	m := make([][]int64, size) //m[i][j] = min value in A[i,j]
	var a, b int64             // temp int(array)
	var u, v int64             // temp vaule
	var d int64                // direction
	var pos, opt int64         // position and optimal value
	a = n
	u = 0
	for i := int64(0); i < size; i++ {
		t[i] = make([]int64, size)
		m[i] = make([]int64, size)

		d = a & 1
		a >>= 1
		if d == 0 {
			u--
		} else {
			u++
		}
		t[i][i] = i
		m[i][i] = u

		b = a
		v = u
		pos = i
		opt = v

		//println("+")
		//if i==0 {
		//println("-", 0, v)
		/*}*/
		for j := i + 1; j < size; j++ {
			d = b & 1
			b >>= 1
			if d == 1 {
				v++
			} else {
				v--
				if v < opt {
					/*               if i==0 {*/
					//println(":", v, opt)
					/*}*/
					pos = j
					opt = v
				}
			}
			/* if i==0 {*/
			//println("-", j, v, pos, opt)
			/*}*/
			t[i][j] = pos
			m[i][j] = opt
		}
	}
	return t
}
