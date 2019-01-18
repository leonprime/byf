package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"
	"torres.guru/gagne/display"
	"torres.guru/gagne/dlx"
	"torres.guru/gagne/game"
)

type Game interface {
	Coverage() *game.Coverage
	Render(io.Writer, []int)
	String() string
}

type Game2D struct {
	pieceSpec string
	w, h      int
	board     *game.Board
}

func (g *Game2D) Coverage() *game.Coverage {
	g.board = game.NewBoard(g.w, g.h, g.pieceSpec)
	return g.board.Coverage
}

func (g *Game2D) Render(w io.Writer, rows []int) {
	plays := g.board.Play(rows)
	display.Render(g.board.W, g.board.H, plays, w)
}

func (g *Game2D) String() string {
	return fmt.Sprintf("%dx%d_%s", g.w, g.h, g.pieceSpec)
}

type Game3D struct {
	pieceSpec string
	w, h, d   int
	cube      *game.Cube
}

func (g *Game3D) Coverage() *game.Coverage {
	g.cube = game.NewCube(g.w, g.h, g.d, g.pieceSpec)
	return g.cube.Coverage
}

func (g *Game3D) Render(w io.Writer, rows []int) {
	plays := g.cube.Play(rows)
	display.Render3D(g.cube.W, g.cube.H, g.cube.D, plays, w)
}

func (g *Game3D) String() string {
	return fmt.Sprintf("%dx%dx%d_%s", g.w, g.h, g.d, g.pieceSpec)
}

func main() {
	path := flag.String("path", ".", "output path for game solutions.")
	max := flag.Int("max", 0, "maximum number of solutions to find.  the default, 0, means find all solutions")
	nprint := flag.Int("print", 10, "number of solutions to print")
	debug := flag.Bool("debug", false, "turn on debugging")
	pieces := flag.String("pieces", "data/pieces.txt", "path to pieces data file")
	show := flag.Bool("show", false, "print available pieces and quit")
	countOnly := flag.Bool("countOnly", false, "don't print, just count solutions")

	flag.Usage = func() {
		f := flag.CommandLine.Output()
		fmt.Fprintf(f, "Usage: %s [options] w h pieceSpec\n", os.Args[0])
		fmt.Fprintf(f, "  w and h are the board width and height\n")
		fmt.Fprintf(f, "  pieceSpec is the set of pieces to play with (see data/pieces.txt)\n")
		fmt.Fprintf(f, "Example: %s 5 3 otzvI\n", os.Args[0])
		fmt.Fprintf(f, "  the solutions are saved at ${path}/solutions/5x3_otzvI\n")
		fmt.Fprintf(f, "Options:\n")
		flag.PrintDefaults()
		os.Exit(2)
	}
	flag.Parse()

	game.LoadPieces(*pieces)
	if *show {
		for _, piece := range game.AllPieces() {
			fmt.Println(piece)
		}
		return
	}

	args := flag.Args()
	var (
		dim       int
		w, h, d   int
		pieceSpec string
		err       error
	)
	if len(args) == 3 {
		dim = 2
		pieceSpec = args[2]
	} else if len(args) == 4 {
		dim = 3
		d, err = strconv.Atoi(args[2])
		if err != nil {
			flag.Usage()
		}
		pieceSpec = args[3]
	} else {
		flag.Usage()
	}
	w, err = strconv.Atoi(args[0])
	if err != nil {
		flag.Usage()
	}
	h, err = strconv.Atoi(args[1])
	if err != nil {
		flag.Usage()
	}
	if w == 0 || h == 0 || len(pieceSpec) == 0 {
		flag.Usage()
	}
	if dim == 3 && d == 0 {
		flag.Usage()
	}

	if *debug {
		game.SetDebug()
		dlx.SetDebug()
	}

	var g Game
	if dim == 2 {
		g = &Game2D{w: w, h: h, pieceSpec: pieceSpec}
	} else {
		g = &Game3D{w: w, h: h, d: d, pieceSpec: pieceSpec}
	}
	run(g, *path, *nprint, *max, *countOnly)
}

func run(g Game, path string, nprint, max int, countOnly bool) {
	cov := g.Coverage()
	dl := dlx.New(cov.M.Cells, cov.Columns, max, countOnly)
	t := time.Now()
	dl.Search(0)
	dur := time.Now().Sub(t)

	fmt.Printf("game \"%s\" has %d solutions\n", g, dl.N)
	fmt.Printf("\ttime taken: %s\n", dur)
	fmt.Printf("\tsteps: %d\n", dl.S)
	if dl.N == 0 || countOnly {
		return
	}

	gamePath := fmt.Sprintf("%s/solutions/%s", path, g)
	os.RemoveAll(gamePath)
	os.MkdirAll(gamePath, os.ModePerm)

	if dl.N < nprint {
		nprint = dl.N
	}
	for i := 0; i < nprint; i++ {
		filename := fmt.Sprintf("%s/%d.png", gamePath, i)
		f, err := os.Create(filename)
		if err != nil {
			panic(err)
		}
		g.Render(f, dl.Solutions[i])
		f.Close()
	}
	quant := "the first"
	if dl.N == nprint {
		quant = "all"
	}
	fmt.Printf("wrote %s %d solutions to %s\n", quant, nprint, gamePath)
}
