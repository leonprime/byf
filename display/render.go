package display

import (
	. "image"
	"image/color"
	"image/png"
	"os"
	. "torres.guru/gagne/game"
)

const (
	tile   = 50
	pad    = 2
	border = 1
)

func width(cols int) int {
	return tile*cols + pad*(cols+1)
}

func height(rows int) int {
	return tile*rows + pad*(rows+1)
}

// returns the tile rectangle given board coords
func tileRect(x, y int) Rectangle {
	x0 := width(x)
	y0 := height(y)
	return Rect(x0, y0, x0+tile, y0+tile)
}

// input w and h of the grid and the plays
// renders the board to a png
func Render(w, h int, plays []*Play) {
	g := &Graf{
		img: NewRGBA(Rect(0, 0, width(w), height(h))),
	}
	g.drawGrid()
	for _, play := range plays {
		g.drawPlay(play)
	}
	g.save()
}

type Graf struct {
	img *RGBA
	c   color.Color
}

func (g *Graf) save() {
	f, err := os.Create("/Users/leon/Desktop/game.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if err := png.Encode(f, g.img); err != nil {
		panic(err)
	}
}

func (g *Graf) HLine(x1, y, x2 int) {
	for ; x1 <= x2; x1++ {
		g.img.Set(x1, y, g.c)
	}
}

func (g *Graf) VLine(x, y1, y2 int) {
	for ; y1 <= y2; y1++ {
		g.img.Set(x, y1, g.c)
	}
}

// Rect draws a filled rectangle
func (g *Graf) DrawRect(x1, y1, x2, y2 int) {
	for ; y1 < y2; y1++ {
		g.HLine(x1, y1, x2)
	}
}

// start by drawing a grid with padding = pad
// this is a bunch of lines
// each grid cell has width = tile and is padded on either side
func (g *Graf) drawGrid() {
	g.c = color.White
	for x := 0; x < g.img.Bounds().Max.X; x += tile + pad {
		g.DrawRect(x, 0, x+pad, g.img.Bounds().Max.Y)
	}
	for y := 0; y < g.img.Bounds().Max.Y; y += tile + pad {
		g.DrawRect(0, y, g.img.Bounds().Max.X, y+pad)
	}
}

// within a tile, there is a border along the edges and an interior
// a piece has a border = border
// the border is on the tile edges
func (g *Graf) drawTile(x, y int, c, b color.Color, u, d, l, r bool) {
	t := tileRect(x, y)
	g.c = c
	g.DrawRect(t.Min.X, t.Min.Y, t.Max.X, t.Max.Y)
	g.c = b
	if u {
		g.DrawRect(t.Min.X, t.Min.Y, t.Max.X, t.Min.Y+border)
	}
	if d {
		g.DrawRect(t.Min.X, t.Max.Y-border, t.Max.X, t.Max.Y)
	}
	if l {
		g.DrawRect(t.Min.X, t.Min.Y, t.Min.X+border, t.Max.Y)
	}
	if r {
		g.DrawRect(t.Max.X-border, t.Min.Y, t.Max.X, t.Max.Y)
	}
}

func (g *Graf) drawPlay(play *Play) {
	pcol := color.RGBA{
		play.Piece.Color[0],
		play.Piece.Color[1],
		play.Piece.Color[2],
		255,
	}
	for x := 0; x < play.Grid.Width(); x++ {
		for y := 0; y < play.Grid.Height(); y++ {
			if play.Grid.Get(x, y) {
				g.drawTile(play.X+x, play.Y+y, pcol, color.Black, false, false, false, false)
			}
		}
	}
}

// to draw a piece we go left to right then top to bottom
// we flip state in then out of piece
// during transition is when we draw the border
// otherwise when state is in, the padding and border are both the piece color
