package main

import(
	"os"
	"fmt"
	"sdl"
	"path"
)

var(
	BoardPath = path.Join("data", "boards")
)

type BoardSize int

const (
	Size9x9   BoardSize = 9
	Size19x19 BoardSize = 19
)

type Board struct {
	size int
	pieces []*Piece

	bg *sdl.Surface
}

func NewBoard(size BoardSize) (*Board, os.Error) {
	b := new(Board)

	b.size = int(size)
	b.pieces = make([]*Piece, b.size*b.size)

	b.bg = sdl.Load(path.Join(BoardPath, fmt.Sprintf("%v.png", b.size)))
	if b.bg == nil {
		return nil, os.NewError(sdl.GetError())
	}

	return b, nil
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
