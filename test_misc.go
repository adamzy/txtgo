package tree

import "testing"

func checkerror(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}
