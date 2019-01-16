package game

// converts a game into a coverage matrix for solving with DLX
func Game2DCoverageMatrix(g *Game2D) {
}

// given the dimensions of a 2d game board, returns all uniquely
// oriented positions of a piece on the game board
func (p *Piece) Positions(w, h int) []*Grid {
	var grids []*Grid
	for _, shape := range p.Shapes {
		grids = append(grids, perms(w, h, shape)...)
		shape = shape.Rotate()
		grids = append(grids, perms(w, h, shape)...)
		shape = shape.Rotate()
		grids = append(grids, perms(w, h, shape)...)
		shape = shape.Rotate()
		grids = append(grids, perms(w, h, shape)...)
	}
	return grids
}

func perms(w, h int, shape *Grid) []*Grid {
	var grids []*Grid
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			// make a grid with piece at this position
			grid := newEmptyGrid(w, h)
			// note that SetSubgrid handles oob conditions with a noop
			grid.SetSubgrid(x, y, shape)
		}
	}
	return grids
}
