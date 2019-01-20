package display

import (
	. "image"
	"image/color"
	"io"
	. "torres.guru/gagne/game"
)

// lay out the z dimension as additional grids going down
// this returns height of image in that representation
// we'll separate the grids by one tile of spacing
func imgh3D(h, d int) int {
	return d*imgh(h) + (d-1)*tile
}

// returns the image y at the grid (y, z) with height h in our flat representation
func y3D(h, y, z int) int {
	stride := z*(h+1) + y
	return tile*stride + pad*(stride+1)
}

func imgRect3D(h, x, y, z int) Rectangle {
	x0 := imgw(x)
	y0 := y3D(h, y, z)
	return Rect(x0, y0, x0+tile, y0+tile)
}

func Render3D(w, h, d int, plays []*Play3D, out io.Writer) {
	g := &Graf{
		img: NewRGBA(Rect(0, 0, imgw(w), imgh3D(h, d))),
	}
	g.drawGrid3D(h, d)
	for _, play := range plays {
		g.drawPlay3D(h, play)
	}
	g.save(out)
}

func (g *Graf) drawGrid3D(h, d int) {
	g.drawGrid()
	// erase in-between lines
	g.c = color.Black
	for z := 1; z < d; z++ {
		stride := z*h + z - 1
		y := tile*stride + pad*(stride+1)
		g.DrawRect(0, y, g.img.Bounds().Max.X, y+tile)
	}
}

func (g *Graf) drawPlay3D(h int, play *Play3D) {
	pcol := color.RGBA{
		play.Piece.Color[0],
		play.Piece.Color[1],
		play.Piece.Color[2],
		255,
	}
	eachTile3D(h, play, pcol, g.drawTile)
	eachTile3D(h, play, borderColor, g.drawBorders)
}

func eachTile3D(h int, play *Play3D, c color.Color, draw func(Rectangle, color.Color, edges)) {
	for z := 0; z < play.Grid.D; z++ {
		for y := 0; y < play.Grid.H; y++ {
			for x := 0; x < play.Grid.W; x++ {
				if play.Grid.Get(x, y, z) {
					b := edges{true, true, true, true}
					if play.Grid.IsSet(x, y-1, z) {
						b.u = false
					}
					if play.Grid.IsSet(x, y+1, z) {
						b.d = false
					}
					if play.Grid.IsSet(x-1, y, z) {
						b.l = false
					}
					if play.Grid.IsSet(x+1, y, z) {
						b.r = false
					}
					t := imgRect3D(h, play.X+x, play.Y+y, play.Z+z)
					draw(t, c, b)
				}
			}
		}
	}
}
