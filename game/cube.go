package game

import "fmt"

type Play3D struct {
	Piece   *Piece
	Grid    *Grid3D
	X, Y, Z int
}

func (p *Play3D) String() string {
	return fmt.Sprintf("%s: (%d, %d, %d) w=%d, h=%d, d=%d\n%s", p.Piece.Name, p.X, p.Y, p.Z, p.Grid.W, p.Grid.H, p.Grid.D, p.Grid)
}

type Cube struct {
	W, H, D  int
	pieces   []*Piece
	Coverage *Coverage
}

func NewCube(w, h, d int, piecesSpec string) *Cube {
	c := &Cube{
		pieces: parsePiecesSpec(piecesSpec),
		W:      w,
		H:      h,
		D:      d,
	}
	c.Coverage = newCubeCoverage(c)
	return c
}

func (c *Cube) Play(rows []int) (plays []*Play3D) {
	for _, y := range rows {
		plays = append(plays, c.play(y))
	}
	return
}

func (c *Cube) play(y int) *Play3D {
	play := &Play3D{}
	p := len(c.pieces)
	row := c.Coverage.M.Row(y)
	//
	// scan for the piece
	for i := 0; i < p; i++ {
		if row[i] {
			play.Piece = c.pieces[i]
			break
		}
	}
	//
	// rebuild the grid from the coverage row
	grid := newEmptyGrid3D(c.W, c.H, c.D)
	for i := p; i < len(row); i++ {
		x := (i - p) % c.W
		y := ((i - p) % (c.W * c.H)) / c.W
		z := (i - p) / (c.W * c.H)
		grid.Set(x, y, z, row[i])
	}
	if debug.piece(play.Piece) {
		fmt.Printf("%s grid rebuilt from coverage row:\n", play.Piece.Name)
		fmt.Println(c.Coverage.RowString(y))
		fmt.Println(grid)
	}
	//
	// scan for piece location and extents
	w, h, d := 0, 0, 0
	found := false
	for y := 0; y < c.H; y++ {
		if found {
			if grid.IsPlaneEmptyY(y) {
				break
			} else {
				h++
			}
		} else {
			if grid.IsPlaneEmptyY(y) {
				play.Y++
			} else {
				h++
				found = true
			}
		}
	}
	found = false
	for x := 0; x < c.W; x++ {
		if found {
			if grid.IsPlaneEmptyX(x) {
				break
			} else {
				w++
			}
		} else {
			if grid.IsPlaneEmptyX(x) {
				play.X++
			} else {
				w++
				found = true
			}
		}
	}
	found = false
	for z := 0; z < c.D; z++ {
		if found {
			if grid.IsPlaneEmptyZ(z) {
				break
			} else {
				d++
			}
		} else {
			if grid.IsPlaneEmptyZ(z) {
				play.Z++
			} else {
				d++
				found = true
			}
		}
	}
	// trim the grid to the subgrid bounding the piece
	play.Grid = grid.GetSubgrid(play.X, play.Y, play.Z, w, h, d)
	if debug.piece(play.Piece) {
		fmt.Printf("play geometry: (%d, %d, %d) w=%d, h=%d, d=%d\n", play.X, play.Y, play.Z, w, h, d)
		fmt.Println(play)
	}
	return play
}
