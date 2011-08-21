package main

import (
	"sdl"
	"path"
)

var (
	PiecePath = path.Join("data", "pieces")
)

type Piece struct {
	img *sdl.Surface
}

func NewPiece(t string) *Piece {
	p := new(Piece)

	p.img = sdl.Load(path.Join(PiecePath, t))

	return p
}
