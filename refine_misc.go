package tree

// Project i onto interval [a,b].
func project(i, a, b int) int {
	switch {
	case i <= a:
		return a
	case i >= b:
		return b
	default:
		return i
	}
}
