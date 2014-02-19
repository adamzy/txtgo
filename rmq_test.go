package rmq

import (
	"testing"
)

func Test_log2(t *testing.T) {
	l := log2(10)
	if v := l[1]; v != 0 {
		t.Error("log2 fail")
	}

	if v := l[2]; v != 1 {
		t.Error("log2 fail")
	}

	if v := l[5]; v != 2 {
		t.Error("log2 fail")
	}

}

func Test_st(t *testing.T) {
	var a = []int64{2, 3, 1, 4, 5, 18, 0, 1, 2, 4, 8, 9, 22}
	//fmt.Println(a)
	ster := st(a)

	if p, v := ster(0, 4); p != 2 || v != 1 {
		t.Error("st fail")
	}

	if p, v := ster(2, 5); p != 2 || v != 1 {
		t.Error("st fail")
	}

	if p, v := ster(2, 6); p != 6 || v != 0 {
		t.Error("st fail")
	}

	if p, v := ster(4, 6); p != 6 || v != 0 {
		t.Error("st fail")
	}

	if p, v := ster(5, 6); p != 6 || v != 0 {
		t.Error("st fail")
	}

	if p, v := ster(0, 10); p != 6 || v != 0 {
		t.Error("st fail")
	}

	if p, v := ster(7, 12); p != 7 || v != 1 {
		t.Error("st fail")
	}

	if p, v := ster(7, 7); p != 7 || v != 1 {
		t.Error("st fail")
	}

	if p, v := ster(7, 8); p != 7 || v != 1 {
		t.Error("st fail")
	}

}

func Test_dynamic(t *testing.T) {
	var a = []int64{2, 3, 2, 3, 4, 5, 6, 5, 4, 5, 6, 7, 6}

	s1 := dynamic(arraytoint(a), 13)

	var p, v int64
	p = s1[0][4]
	v = a[p]
	if p != 0 || v != 2 {
		println(p, v)
		t.Error("st fail")
	}

	p = s1[0][2]
	v = a[p]
	if p != 0 || v != 2 {
		println(p, v)
		t.Error("st fail")
	}

	p = s1[2][6]
	v = a[p]
	if p != 2 || v != 2 {
		println(p, v)
		t.Error("st fail")
	}

	p = s1[1][6]
	v = a[p]
	if p != 2 || v != 2 {
		t.Error("st fail")
	}

	p = s1[5][12]
	v = a[p]
	if p != 8 || v != 4 {
		t.Error("st fail")
	}

}

func Test_processBlock(t *testing.T) {
	var a = []int64{2, 3, 2, 3, 4, 5, 6, 5, 4, 5, 6, 7, 6}
	pro := processBlock(int64(len(a)))
	s1 := pro(a)

	var p, v int64
	p = s1[0][4]
	v = a[p]
	if p != 0 || v != 2 {
		println(p, v)
		t.Error("st fail")
	}

	p = s1[0][2]
	v = a[p]
	if p != 0 || v != 2 {
		println(p, v)
		t.Error("st fail")
	}

	p = s1[2][6]
	v = a[p]
	if p != 2 || v != 2 {
		println(p, v)
		t.Error("st fail")
	}

	p = s1[1][6]
	v = a[p]
	if p != 2 || v != 2 {
		t.Error("st fail")
	}

	p = s1[5][12]
	v = a[p]
	if p != 8 || v != 4 {
		t.Error("st fail")
	}
}

func Test_ResRMQ(t *testing.T) {
	var a = []int64{2, 3, 2, 3, 4, 5, 6, 5, 4, 5, 6, 7, 6}
	ster := ResRMQ(a)

	if p, v := ster(0, 4); p != 0 || v != 2 {
		println(p, v)
		t.Error("st fail")
	}

	if p, v := ster(2, 5); p != 2 || v != 2 {
		t.Error("st fail")
	}

	if p, v := ster(1, 6); p != 2 || v != 2 {
		t.Error("st fail")
	}

	if p, v := ster(5, 10); p != 8 || v != 4 {
		t.Error("st fail")
	}

	if p, v := ster(5, 8); p != 8 || v != 4 {
		t.Error("st fail")
	}

	if p, v := ster(0, 12); p != 0 || v != 2 {
		t.Error("st fail")
	}

	if p, v := ster(10, 12); p != 10 || v != 6 {
		t.Error("st fail")
	}
}

func Test_ResRMQ2(t *testing.T) {
	var a = []int64{0,1,2,3,4,3,4,3,2,3,2,1,2,1,0,1,0}
	ster := ResRMQ(a)

	if p, v := ster(6, 12); p != 11 || v != 1 {
		println(p, v)
		t.Error("st fail")
	}

}
