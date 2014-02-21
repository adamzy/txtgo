package tree

import (
	"regexp"
	"strconv"
	"strings"
)

type Tree struct {
	*Node
	Nodes []*Node
}

type InvalidTreeError struct {
	TreeString string
}

func (err InvalidTreeError) Error() string {
	return "Invalid tree: " + err.TreeString
}

type InvalidEdgeLengthError struct {
	EdgeLength string
}

func (err InvalidEdgeLengthError) Error() string {
	return "Invalid edge length: " + err.EdgeLength
}

// Maybe using regexp is not fast enough.
//
// **TODO** Write a native one.
func tokenize3(s string) []string {
	//s = strings.Replace(s, " ", "", -1)
	reg := `[\(\),;:]|[^\(\),;:\s]+`
	re := regexp.MustCompile(reg)
	return re.FindAllString(s, -1)
}

// This one seems faster.
func tokenize(s string) []string {
	//s = strings.Replace(s, " ", "", -1)
	r := strings.NewReplacer("(", " ( ", ")", " ) ", ":", " : ", ";", " ; ", ",", " , ")
	return strings.Fields(r.Replace(s))
}

// Make a tree from a string
func Make(s string) (*Tree, error) {
	token := tokenize(s)
	if len(token) == 0 {
		return nil, InvalidTreeError{s}
	}
	t := new(Tree)
	t.Node = newNode()

	var flag bool
	var n, c *Node
	n = t.Node

	for _, v := range token {
		switch v {
		case "(":
			c = newNode()
			n.AddChild(c)
			n = c
			flag = false
		case ")":
			if n.Father == nil {
				return nil, InvalidTreeError{s}
			}
			n = n.Father
			flag = false
		case ",":
			c = newNode()
			if n.Father == nil {
				return nil, InvalidTreeError{s}
			}
			n.Father.AddChild(c)
			n = c
			flag = false
		case ":":
			flag = true
		case ";":
			break
		default:
			if flag {
				length, err := strconv.ParseFloat(v, 64)
				if err != nil {
					return nil, InvalidEdgeLengthError{v}
				}
				n.Length = length
			} else {
				n.Name = v
				flag = false
			}
		}
	}

	if t.Node != n {
		return nil, InvalidTreeError{s}
	}

	// Update t.Nodes once the tree is successfully built.
	t.Nodes = t.Post2List()
	t.UpdateInfo()
	return t, nil
}

// Update node.Level and node.Id for tree `t`
func (t *Tree) UpdateInfo() {
	size := len(t.Nodes)
	t.Node.Level = 0
	t.Node.Id = size - 1
	for i := size - 2; i >= 0; i-- {
		n := t.Nodes[i]
		n.Level = n.Father.Level + 1
		n.Id = i
	}
}

// Update everything of a tree from `tree.Node`,
// including `tree.Nodes, size, node.id, level`.
// Useful after the tree was manually edited
func (t *Tree) Update() {
	t.Nodes = t.Node.Post2List()
	t.UpdateInfo()
}

// A tree is binary tree <=> every internal node has two children.
// Overwrite `node.IsBinary` method.
func (t *Tree) IsBinary() bool {
	for i := range t.Nodes {
		n := t.Nodes[i]
		if n.IsInternal() && !n.IsBinary() {
			return false
		}
	}
	return true
}

// A map that maps taxa to (one of) the corresponding leaf node
type Taxonmap map[string]*Node

// `TaxonMap` return the `Taxonmap` for a tree, and the uniqueness of its taxon.
func (t *Tree) TaxonMap() (taxon Taxonmap, unique bool) {
	taxon = make(Taxonmap, len(t.Nodes))
	i := 0
	for _, n := range t.Nodes {
		if n.IsLeaf() {
			taxon[n.Name] = n
			i++
		}
	}
	if i == len(taxon) {
		unique = true
	}
	return
}
