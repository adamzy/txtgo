package tree

import (
    "strconv"
    "strings"
    "errors"
    "regexp"
)

type Tree struct {
    *Node
    Nodes []*Node
    Size int
}

// Maybe using regexp is not fast enough.
// TODO Write a native one.
func tokenize(s string) []string {
    s = strings.Replace(s, " ", "", -1) 
    reg := `[\(\),;:]|[^\(\),;:]+`
    re := regexp.MustCompile(reg)
    return re.FindAllString(s, -1)
}

// Dont know which one is faster
func tokenize2(s string) []string {
    s = strings.Replace(s, " ", "", -1) 
    r := strings.NewReplacer("(", " ( ", ")", " ) ", ":", " : ", ";", " ; ", ",", " , ")
    return strings.Fields(r.Replace(s))
}

// Make a tree from a string
func Make(s string) (*Tree, error) {
    token := tokenize(s)
    if len(token) == 0 {return nil, errors.New("Invalid tree.")}
    t := new(Tree)
    t.Node = newNode()
    
    size := 1
    flag := false
    var n, c *Node
    n = t.Node

    for _, v := range token {
        switch v {
        case "(":
            c = newNode()
            size++
            n.AddChild(c)
            n = c
            flag = false
        case ")":
            n = n.Father
            flag = false
        case ",":
            c = newNode()
            size++
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
                    return nil, errors.New("Invalid branch length.")
                }
                n.Length = length
            } else {
                n.Name = v
                flag = false
            }
        }
    }

    if t.Node != n {
        return nil, errors.New("Invalid tree.")
    }

    t.Size = size
    // Update t.Nodes once the tree is successfully built.
    t.Nodes = t.Post2List()
    t.UpdateInfo()
    return t, nil
}

// Update node.Level and node.Id for tree
func (t *Tree) UpdateInfo() {
    t.Node.Level = 0
    t.Node.Id = len(t.Nodes)-1
    for i:=len(t.Nodes)-2; i>=0; i-- {
        n := t.Nodes[i]
        n.Level = n.Father.Level + 1
        n.Id = i
    }
}

// A tree is binary tree iff every internal node has two children.
// Overwirte node.IsBinary method.
func (t *Tree) IsBinary() bool {
    for i := range t.Nodes {
        n := t.Nodes[i]
        if n.IsInternal() && !n.IsBinary() {return false}
    }
    return true
}

//func (t *Tree) String () string {
    //return t.Root.String()
//}

//func (t *Tree) toString (showlength bool) string {
    //return t.Root.toString(showlength)
//}


// The following is uncessary.

// Postorder iterate the tree and put nodes into a list
// Rewrited method of node.Post2List, because the size is known for tree.
func (tree *Tree) Post2List () []*Node {
    //var nl []*Node
    nl := make([]*Node, 0, tree.Size)

    // A function that appends node with its descandent to nl
    var ap func (n *Node) //necessary for making a recursive function here
    ap = func (n *Node) {
        for _, c :=range n.Children {
            ap(c)
        }
        nl = append(nl, n)
    }

    ap(tree.Node)
    return nl
}

// Inorder iterate the tree and put nodes into a list.
// It's only defined for binary tree.
// It will cause panic for non-bianry tree.
// Rewrited method of node.Post2List, because the size is known for tree.
// Ugly!
func (tree *Tree) In2List() []*Node {
    //var nl []*Node
    nl := make([]*Node, 0, tree.Size)

    // A function that appends node with its descandent to nl
    var ap func (n *Node) //necessary for making a recursive function here
    ap = func (n *Node) {
        if n.IsInternal() {
            ap(n.Children[0])
        }
        nl = append(nl, n)
        if n.IsInternal() {
            ap(n.Children[1])
        }
    }

    ap(tree.Node)
    return nl
}
