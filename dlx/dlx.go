package dlx

import (
	"bytes"
	"fmt"
	"math"
)

var debug bool

func SetDebug() {
	debug = true
}

type Node struct {
	L, R, U, D *Node
	C          *Column
	N          string
	x, y       int
}

func (n *Node) String() string {
	return n.N
}

type Column struct {
	Node
	S int
}

// given a boolean matrix, builds the corresponding dancing links matrix A
// for use in DLX search algorithm
func New(matrix [][]bool, columnNames []string, max int, countOnly bool) *DancingLinks {
	w, h := len(matrix[0]), len(matrix)
	if len(columnNames) != w {
		panic("number of column names doesn't match number of matrix columns")
	}

	root := &Column{Node: Node{N: "root"}, S: 0}
	root.C = root
	var cols []*Column
	var links []*Node // this is to carry forward the "last link" node

	// build the L/R columns row
	for x := 0; x < w; x++ {
		cols = append(cols, &Column{Node: Node{N: columnNames[x]}, S: 0})
		cols[x].C = cols[x]
		// link the previous col to this one
		if x == 0 {
			cols[x].L = &root.Node
			root.R = &cols[x].Node
		} else {
			cols[x].L = &cols[x-1].Node
			cols[x-1].R = &cols[x].Node
		}
		links = append(links, &cols[x].Node)
	}
	// comlpete the L/R circular links
	cols[len(cols)-1].R = &root.Node
	root.L = &cols[len(cols)-1].Node

	// build the nodes top to bottom and do U/D linking
	sizes := make([]int, w, w)
	var nodes [][]*Node // this is to complete lnking of L/R
	for y := range matrix {
		var row []*Node
		for x := range matrix[y] {
			if len(matrix[y]) != w {
				panic(fmt.Sprintf("bad matrix: w is %d but row %d len is %d", w, y, len(matrix[y])))
			}
			if !matrix[y][x] {
				row = append(row, nil)
				continue
			}
			node := &Node{
				C: cols[x],
				x: x,
				y: y,
				N: fmt.Sprintf("n(%s,%d)", cols[x].String(), y),
			}
			links[x].D = node
			node.U = links[x]
			links[x] = node
			sizes[x] += 1
			row = append(row, node)
		}
		nodes = append(nodes, row)
	}
	// comlpete the U/D circular links
	// and mark the sizes
	for x := 0; x < w; x++ {
		links[x].D = &cols[x].Node
		cols[x].U = links[x]
		cols[x].S = sizes[x]
	}

	// to complete the L/R links for nodes, we build a fake row header
	rowh := make([]*Node, 0, 0)
	links = make([]*Node, 0, 0)
	for y := 0; y < h; y++ {
		node := &Node{}
		rowh = append(rowh, node)
		links = append(links, node)
	}

	// link the L/R nodes left to right
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			node := nodes[y][x]
			if node == nil {
				continue
			}
			links[y].R = node
			node.L = links[y]
			links[y] = node
		}
	}
	// complete circular L/R links and delete rowh nodes
	for y := 0; y < h; y++ {
		rowh[y].R.L = links[y]
		links[y].R = rowh[y].R
		rowh[y].R = nil
		rowh[y].L = nil
	}

	dl := &DancingLinks{root: root, max: max, countOnly: countOnly}
	if debug {
		fmt.Println(dl)
	}
	return dl
}

type Solution []int

type DancingLinks struct {
	root      *Column
	o         []*Node
	Solutions []Solution
	max       int  // max solutions to search for (0 is all)
	N, S      int  // number of solutions found and steps taken
	countOnly bool // skip generation of Solutions
}

func (dl *DancingLinks) Search(k int) {
	dl.S++
	if dl.root.R == &dl.root.Node {
		dl.printSolution()
		return
	}
	if dl.max > 0 && dl.N >= dl.max {
		return
	}
	dl.o = append(dl.o, nil)
	c := dl.chooseColumn()
	if debug {
		fmt.Printf("k is %d\n", k)
		fmt.Printf("column choice is %s\n", c)
	}
	dl.cover(c)
	for r := c.D; r != &c.Node; r = r.D {
		dl.o[k] = r
		for j := r.R; j != r; j = j.R {
			dl.cover(j.C)
		}
		dl.Search(k + 1)
		r = dl.o[k]
		c = r.C
		for j := r.L; j != r; j = j.L {
			dl.uncover(j.C)
		}
	}
	dl.uncover(c)
}

// method that minimizes branching
func (dl *DancingLinks) chooseColumn() (c *Column) {
	dl.S++
	s := math.MaxInt32
	for col := dl.root.R; col != &dl.root.Node; col = col.R {
		if col.C.S < s {
			c = col.C
			s = col.C.S
		}
	}
	return
}

// remove c from the header list and remove all rows in c's own list from other column lists they are in
func (dl *DancingLinks) cover(c *Column) {
	dl.S++
	if debug {
		fmt.Printf("covering %s\n", c)
	}
	c.R.L = c.L
	c.L.R = c.R
	for i := c.D; i != &c.Node; i = i.D {
		for j := i.R; j != i; j = j.R {
			j.D.U = j.U
			j.U.D = j.D
			j.C.S -= 1
		}
	}
}

// the meat of the dancing links
func (dl *DancingLinks) uncover(c *Column) {
	dl.S++
	if debug {
		fmt.Printf("uncovering %s\n", c)
	}
	for i := c.U; i != &c.Node; i = i.U {
		for j := i.L; j != i; j = j.L {
			j.C.S += 1
			j.D.U = j
			j.U.D = j
		}
	}
	c.R.L = &c.Node
	c.L.R = &c.Node
}

func (dl *DancingLinks) String() string {
	var s bytes.Buffer
	s.WriteString("Columns:\n")
	str := fmt.Sprintf("%s: R=%s, L=%s, U=%s, D=%s", dl.root, dl.root.R.C, dl.root.L.C, dl.root.U, dl.root.D)
	s.WriteString(str)
	s.WriteRune('\n')
	for c := dl.root.R; c != &dl.root.Node; c = c.R {
		str := fmt.Sprintf("%s: R=%s, L=%s, U=%s, D=%s", c.C, c.R.C, c.L.C, c.U, c.D)
		s.WriteString(str)
		s.WriteRune('\n')
	}
	s.WriteString("\nNodes (by column top to bottom):\n")
	for c := dl.root.R; c != &dl.root.Node; c = c.R {
		s.WriteString(c.C.String())
		s.WriteRune('\n')
		for i := c.D; i != c; i = i.D {
			s.WriteRune('\t')
			str := fmt.Sprintf("%s: R=%s, L=%s, U=%s, D=%s", i, i.R, i.L, i.U, i.D)
			s.WriteString(str)
			s.WriteRune('\n')
		}
	}
	return s.String()
}

func (dl *DancingLinks) printSolution() {
	if debug {
		var buf bytes.Buffer
		buf.WriteString("========\n")
		buf.WriteString("solution\n")
		buf.WriteString("========\n")
		for _, o := range dl.o {
			if o == nil {
				break
			}
			//
			// print the row that includes node o
			for i := o.R; i != o; i = i.R {
				buf.WriteString(fmt.Sprintf("%-12s", i.String()))
				buf.WriteString("\t")
			}
			buf.WriteRune('\n')
		}
		fmt.Println(buf.String())
	}
	dl.N++
	if dl.N%1000 == 0 {
		fmt.Printf("\rfound %d solutions", dl.N)
	}
	if dl.countOnly {
		return
	}
	soln := Solution{}
	for _, o := range dl.o {
		if o == nil {
			break
		}
		soln = append(soln, o.y)
	}
	dl.Solutions = append(dl.Solutions, soln)
}
