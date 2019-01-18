package game

var debug debugOpts

type debugOpts struct {
	p   *Piece
	all bool
}

func SetDebug(p string) {
	if piece, ok := allPieces[p]; ok {
		debug = debugOpts{p: piece}
	} else {
		panic("piece not found: " + p)
	}
}

func SetDebugAll() {
	debug = debugOpts{all: true}
}

func (d debugOpts) piece(p *Piece) bool {
	if d.all {
		return true
	}
	return d.p == p
}
