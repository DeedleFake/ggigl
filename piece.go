package main

import (
	"fmt"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

// Stores information for the piece.
type Piece struct {
	img *ebiten.Image
}

var (
	Black = &Piece{img: pieceBlackImage}
	White = &Piece{img: pieceWhiteImage}
)

// Returns an image of the piece.
func (p *Piece) Image() *ebiten.Image {
	return p.img
}

func (p *Piece) HighlightColor() *image.Uniform {
	switch p {
	case Black:
		return blackHighlight
	case White:
		return whiteHighlight
	default:
		panic(fmt.Errorf("invalid piece: %v", p))
	}
}
