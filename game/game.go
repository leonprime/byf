package game

import (
	"fmt"
)

// a 2D game is a set of pieces and the dimensions of the playing area
type Game2D struct {
	Pieces     []*Piece
	rows, cols int
}

func New2D(rows, cols int, pieces_spec string) *Game2D {
	if Pieces == nil {
		panic("ensure LoadPieces(file) is called first")
	}
	var pieces []*Piece
	for _, char := range pieces_spec {
		if piece, ok := Pieces[string(char)]; ok {
			pieces = append(pieces, piece)
		} else {
			panic(fmt.Sprintf("no piece \"%c\" defined", char))
		}
	}
	return &Game2D{
		Pieces: pieces,
		rows:   rows,
		cols:   cols,
	}
}
