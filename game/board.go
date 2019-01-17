package game

import (
	"fmt"
)

var debug bool

func SetDebug() {
	debug = true
}

type Play struct {
	piece *Piece
	grid  *Grid
	x, y  int
}

func (p *Play) String() string {
	return fmt.Sprintf("%s: (%d, %d) w=%d, h=%d\n%s", p.piece.Name, p.x, p.y, p.grid.w, p.grid.h, p.grid)
}

type BoardGame struct {
	w, h     int
	pieces   []*Piece
	Coverage *Coverage
}

func NewBoardGame(w, h int, pieces_spec string) *BoardGame {
	if allPieces == nil {
		panic("ensure LoadPieces(file) is called first")
	}
	var pieces []*Piece
	for _, char := range pieces_spec {
		if piece, ok := allPieces[string(char)]; ok {
			pieces = append(pieces, piece)
		} else {
			panic(fmt.Sprintf("no piece \"%c\" defined", char))
		}
	}
	b := &BoardGame{
		pieces: pieces,
		w:      w,
		h:      h,
	}
	b.Coverage = newCoverage(b)
	return b
}

// play the solution by reading the selected rows
// from the coverage matrix and related data
// returns a list of play objects that represent the placement of each piece
func (b *BoardGame) Play(rows []int) (plays []*Play) {
	for _, y := range rows {
		plays = append(plays, b.play(y))
	}
	return
}

func (b *BoardGame) play(y int) *Play {
	play := &Play{}
	p := len(b.pieces)
	row := b.Coverage.M.Row(y)
	//
	// scan for the piece
	for i := 0; i < p; i++ {
		if row[i] {
			play.piece = b.pieces[i]
			break
		}
	}
	//
	// rebuild the grid from the row
	grid := newEmptyGrid(b.w, b.h)
	for i := p; i < len(row); i++ {
		x := (i - p) % b.w
		y := (i - p) / b.w
		grid.Set(x, y, row[i])
	}
	if debug {
		fmt.Println(grid)
	}
	//
	// scan for piece location and extents
	w, h := 0, 0
	found := false
	for y := 0; y < b.h; y++ {
		if found {
			if grid.IsRowEmpty(y) {
				break
			} else {
				h++
			}
		} else {
			if grid.IsRowEmpty(y) {
				play.y++
			} else {
				h++
				found = true
			}
		}
	}
	found = false
	for x := 0; x < b.w; x++ {
		if found {
			if grid.IsColEmpty(x) {
				break
			} else {
				w++
			}
		} else {
			if grid.IsRowEmpty(x) {
				play.x++
			} else {
				w++
				found = true
			}
		}
	}
	// trim the grid to the subgrid bounding the piece
	play.grid = grid.GetSubgrid(play.x, play.y, w, h)
	if debug {
		fmt.Println(play)
	}
	return play
}
