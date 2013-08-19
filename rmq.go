// Assume size of input array <= 2*64-1
package rmq

//func log2(x int64) int64 {
//n := 0
//for ;x!=0;x>>=1 {
//n++
//}
//return n
//}

// Sparse Table algorithm (time complexity O(n*log(n))
func st(A []int64) func(x, y int64) (posi, min int64) {
	size := int64(len(A))
	log := log2(size) // log_2
	s := log[len(A)]
	pow := power2(s) // pow_2

	M := make([][]int64, size)
	N := make([][]int64, size)
	for i := range M {
		M[i] = make([]int64, s+1)
		N[i] = make([]int64, s+1)
	}

	for i := int64(0); i < size-1; i++ {
		if A[i] <= A[i+1] {
			M[i][0] = A[i]
			N[i][0] = i
		} else {
			M[i][0] = A[i+1]
			N[i][0] = i + 1
		}
	}

	for j := int64(1); j <= s; j++ {
		for i := int64(0); i < size-pow[j]; i++ {
			if M[i][j-1] <= M[i+pow[j-1]][j-1] {
				M[i][j] = M[i][j-1]
				N[i][j] = N[i][j-1]
			} else {
				M[i][j] = M[i+pow[j-1]][j-1]
				N[i][j] = N[i+pow[j-1]][j-1]
			}
		}
	}

	return func(x, y int64) (int64, int64) {
		//fmt.Println(len(log), x, y)
		if x == y { // if y-x == 0, don't use log and pow
			return N[x][0], M[x][0]
		}

		r := log[y-x]
		if M[x][r] > M[y-pow[r]][r] {
			return N[y-pow[r]][r], M[y-pow[r]][r]
		} else {
			return N[x][r], M[x][r]
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

// log[x] return n s.t. 2^n <= x < 2^(n+1), i.e. log[x] = floor(log_2(x))
func log2(size int64) []int64 {
	log := make([]int64, size+1)
	log[0] = 0
	i := int64(1)
	j := int64(2)
	k := int64(0)
	for j <= size+1 {
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
	if len(A) < 2 {
		return func(x, y int64) (int64, int64) {
			return 0, 0
		}
	}
	length := int64(len(A))

	// n is ceil(log(2))
	n := int64(0)
	for x := length; x != 0; x >>= 1 {
		n++
	}

	//taken from graphics.stanford.edu/~seander/bithacks.html
	// check if size is a power of 2
	if (length & (length - 1)) == 0 {
		n--
	}

	// size is the size of each block
	size := n / 2
	if n%2 != 0 {
		size++
	}

	// num is the number of blocks
	num := length / size
	if length%size != 0 {
		num++
	}
	// loc[i] is the block in which i is.
	loc := make([]int64, length)
	// res[i] is the offset of i in block loc[i].
	res := make([]int64, length)

	// rre[i] is the position of the first element in block i.
	pre := make([]int64, num)

	// B[i][j][k] is the optimal position in block i, from position j to k.
	B := make([][][]int64, num)

	// for sparse table algorithm
	C := make([]int64, num) // posi
	D := make([]int64, num) // value

	var i, j int64
	processer := processBlock(size) //, 1<<uint64(size-1))

	var pos int64
	for pos = 0; pos+size < length; pos += size {
		for j = 0; j < size; j++ {
			loc[pos+j] = i
			res[pos+j] = j
		}
		pre[i] = pos
		B[i] = processer(A[pos : pos+size])
		C[i] = B[i][0][size-1] + pos
		//D[i] = E[i][0][size-1]
		D[i] = A[C[i]]
		//fmt.Println("_", A[pos:pos+size], B[i][0][size-1], D[i])
		i++
	}
	if pos != length {
		// pos < length
		// i = num-1
		for j = 0; j < length-pos; j++ {
			loc[pos+j] = i
			res[pos+j] = j
		}
		pre[i] = pos
		B[i] = processer(A[pos:length])
		C[i] = B[num-1][0][length-pos-1] + pos
		//D[num-1] = E[num-1][0][length-pos-1]
		D[i] = A[C[num-1]]
	}
	tabler := st(D)

	/*fmt.Println(C)*/
	/*fmt.Println(D)*/
	//for i:=int64(0);i<num;i++ {
	//for j:=i; j<num; j++ {
	//a, b := tabler(i,j)
	//fmt.Println(i,j,a,b)
	//}
	//}

	return func(x, y int64) (int64, int64) {
		// always assume 0 <= x <= y < length
		u, v := loc[x], loc[y]
		ur, vr := res[x], res[y]
		up, vp := pre[u], pre[v]

		//fmt.Println("~", x, y, ":", u, v, ur, vr, up, vp)
		if u == v {
			//return B[u][ur][vr], E[u][ur][vr]
			pos := up + B[u][ur][vr]
			return pos, A[pos]
		} else {
			// u < v
			p1 := up + B[u][ur][size-1]
			//v1 := E[u][ur][size-1]
			v1 := A[p1]
			p2 := vp + B[v][0][vr]
			//v2 := E[v][0][vr]
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
				//println(v1, v2, v3)
			}
			return pp, vv
		}
	}
}

func arraytoint(A []int64) int64 {
	if len(A) > 64 {
		panic("Overflow")
	}
	var v int64
	for i := len(A) - 1; i > 0; i-- {
		if A[i] == A[i-1]+1 {
			v |= 1
		} else if A[i] != A[i-1]-1 {
			panic("Not +-1 array")
		}
		v <<= 1
	}
	return v
}

func processBlock(size int64) func([]int64) [][]int64 {
        // num = 2^(size)
	num := int64(1 << uint64(size))
	a := make([][][]int64, num)
	for i := int64(0); i < num; i++ {
		a[i] = dynamic(i, size)
	}
	return func(A []int64) [][]int64 {
		v := arraytoint(A)
		return a[v]
	}
}

func dynamic(n, size int64) [][]int64 {
	if size > 64 {
		panic("Cannot")
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

func dynamicArray(A []int64) [][]int64 {
	size := int64(len(A))
	t := make([][]int64, size) //t[i][j] = position with min value in A[i,j]
	m := make([][]int64, size) //m[i][j] = min value in A[i,j]

	for i := int64(0); i < size; i++ {
		t[i] = make([]int64, size)
		m[i] = make([]int64, size)
		t[i][i] = i
		m[i][i] = A[i]
		for j := i + 1; j < size; j++ {
			if A[j] < m[i][j-1] {
				t[i][j] = j
				m[i][j] = A[j]
			} else {
				t[i][j] = t[i][j-1]
				m[i][j] = m[i][j-1]
			}
		}
	}
	return t
}
