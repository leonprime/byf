package display

import (
	. "image"
	"image/color"
	"image/png"
	"io"
	. "torres.guru/gagne/game"
)

const (
	tile   = 50
	pad    = 3
	border = 1
)

var borderColor = color.RGBA{0xBD, 0xBD, 0xBD, 0xBD}

// converts game w to img w
func imgw(cols int) int {
	return tile*cols + pad*(cols+1)
}

// converts game h to img h
func imgh(rows int) int {
	return tile*rows + pad*(rows+1)
}

// returns an img rectangle given board coords
func imgRect(x, y int) Rectangle {
	x0 := imgw(x)
	y0 := imgh(y)
	return Rect(x0, y0, x0+tile, y0+tile)
}

// input w and h of the grid and the plays
// renders the board to a png
func Render(w, h int, plays []*Play, out io.Writer) {
	g := &Graf{
		img: NewRGBA(Rect(0, 0, imgw(w), imgh(h))),
	}
	g.drawGrid()
	for _, play := range plays {
		g.drawPlay(play)
	}
	g.save(out)
}

type Graf struct {
	img *RGBA
	c   color.Color
}

func (g *Graf) save(w io.Writer) error {
	return png.Encode(w, g.img)
}

// Rect draws a filled rectangle
func (g *Graf) DrawRect(x1, y1, x2, y2 int) {
	for y := y1; y < y2; y++ {
		for x := x1; x <= x2; x++ {
			g.img.Set(x, y, g.c)
		}
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

type edges struct {
	u, d, l, r bool
}

// within a tile, there is a border along the edges and an interior
// a piece has a border = border
// the border is on the tile edges
func (g *Graf) drawTile(t Rectangle, tc color.Color, b edges) {
	g.c = tc
	g.DrawRect(t.Min.X, t.Min.Y, t.Max.X, t.Max.Y)

	// connected edges
	if !b.u {
		g.DrawRect(t.Min.X, t.Min.Y-pad, t.Max.X, t.Min.Y)
	}
	if !b.d {
		g.DrawRect(t.Min.X, t.Max.Y, t.Max.X, t.Max.Y+pad)
	}
	if !b.l {
		g.DrawRect(t.Min.X-pad, t.Min.Y, t.Min.X, t.Max.Y)
	}
	if !b.r {
		g.DrawRect(t.Max.X, t.Min.Y, t.Max.X+pad, t.Max.Y)
	}
}

func (g *Graf) drawBorders(t Rectangle, bc color.Color, b edges) {
	g.c = bc
	if b.u {
		g.DrawRect(t.Min.X, t.Min.Y, t.Max.X, t.Min.Y+border)
	} else {
		g.DrawRect(t.Min.X, t.Min.Y-pad-border, t.Min.X+border, t.Min.Y)
		g.DrawRect(t.Max.X-border, t.Min.Y-pad-border, t.Max.X, t.Min.Y)
	}
	if b.d {
		g.DrawRect(t.Min.X, t.Max.Y-border, t.Max.X, t.Max.Y)
	} else {
		g.DrawRect(t.Min.X, t.Max.Y, t.Min.X+border, t.Max.Y+pad+border)
		g.DrawRect(t.Max.X-border, t.Max.Y, t.Max.X, t.Max.Y+pad+border)
	}
	if b.l {
		g.DrawRect(t.Min.X, t.Min.Y, t.Min.X+border, t.Max.Y)
	} else {
		g.DrawRect(t.Min.X-pad-border, t.Min.Y, t.Min.X, t.Min.Y+border)
		g.DrawRect(t.Min.X-pad-border, t.Max.Y-border, t.Min.X, t.Max.Y)
	}
	if b.r {
		g.DrawRect(t.Max.X-border, t.Min.Y, t.Max.X, t.Max.Y)
	} else {
		g.DrawRect(t.Max.X, t.Min.Y, t.Max.X+pad+border, t.Min.Y+border)
		g.DrawRect(t.Max.X, t.Max.Y-border, t.Max.X+pad+border, t.Max.Y)
	}
}

func (g *Graf) drawPlay(play *Play) {
	pcol := color.RGBA{
		play.Piece.Color[0],
		play.Piece.Color[1],
		play.Piece.Color[2],
		255,
	}
	eachTile(play, pcol, g.drawTile)
	eachTile(play, borderColor, g.drawBorders)
}

func eachTile(play *Play, c color.Color, draw func(Rectangle, color.Color, edges)) {
	for x := 0; x < play.Grid.W; x++ {
		for y := 0; y < play.Grid.H; y++ {
			if play.Grid.Get(x, y) {
				b := edges{true, true, true, true}
				if play.Grid.IsSet(x, y-1) {
					b.u = false
				}
				if play.Grid.IsSet(x, y+1) {
					b.d = false
				}
				if play.Grid.IsSet(x-1, y) {
					b.l = false
				}
				if play.Grid.IsSet(x+1, y) {
					b.r = false
				}
				t := imgRect(play.X+x, play.Y+y)
				draw(t, c, b)
			}
		}
	}
}
