package main

type BoardSize int

const (
	Size9x9   BoardSize = 9
	Size19x19 BoardSize = 19
)

type Board struct {
	pieces []*Piece

	size int
}

func NewBoard(size BoardSize) (b *Board) {
	b = new(Board)

	b.size = int(size)
	b.pieces = make([]*Piece, b.size*b.size)

	return
}

func (b *Board) At(x, y int) *Piece {
	return b.pieces[(y*b.size)+x]
}

func (b *Board) place(x, y int, p *Piece) {
	b.pieces[(y*b.size)+x] = p
}

func (b *Board) Place(x, y int, p *Piece) bool {
	if b.At(x, y) != nil {
		return false
	}

	b.place(x, y, p)

	return true
}
