package main

import (
	"errors"
	"github.com/banthar/Go-SDL/sdl"
	"path"
	"runtime"
)

var (
	// The path to the piece data.
	PiecePath = path.Join("data", "pieces")
)

// Stores information for the piece.
type Piece struct {
	img *sdl.Surface
}

// Initializes a new piece.
func NewPiece(t string) (*Piece, error) {
	p := new(Piece)

	p.img = sdl.Load(path.Join(PiecePath, t+".png"))
	if p.img == nil {
		return nil, errors.New(sdl.GetError())
	}

	runtime.SetFinalizer(p, (*Piece).free)

	return p, nil
}

// Frees resources associated with the piece.
func (p *Piece) free() {
	p.img.Free()
}

// Returns an image of the piece.
func (p *Piece) Image() *sdl.Surface {
	return p.img
}
