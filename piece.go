package main

import (
	"os"
	"sdl"
	"path"
)

var (
	PiecePath = path.Join("data", "pieces")
)

type Piece struct {
	img *sdl.Surface
}

func NewPiece(t string) (*Piece, os.Error) {
	p := new(Piece)

	p.img = sdl.Load(path.Join(PiecePath, t+".png"))
	if p.img == nil {
		return nil, os.NewError(sdl.GetError())
	}

	return p, nil
}
