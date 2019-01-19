package game

var debug *debugOpts = &debugOpts{}

type debugOpts struct {
	cov, allPieces bool
	p              *Piece
}

func SetDebugAllPieces() {
	debug.allPieces = true
}

func SetDebugPiece(p string) {
	if piece, ok := allPieces[p]; ok {
		debug.p = piece
	} else {
		panic("piece not found: " + p)
	}
}

func SetDebugCoverage() {
	debug.cov = true
}

func (d *debugOpts) piece(p *Piece) bool {
	if d.allPieces {
		return true
	}
	return d.p == p
}

func (d *debugOpts) coverage() bool {
	return d.cov
}
