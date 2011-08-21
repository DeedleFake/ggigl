package main

import (
	"os"
	"fmt"
	"sdl"
	"path"
	"runtime"
)

var (
	BoardPath = path.Join("data", "boards")
)

type BoardSize int

const (
	Size9x9   BoardSize = 9
	Size19x19 BoardSize = 19
)

type Board struct {
	size   int
	pieces []*Piece

	bg  *sdl.Surface
	img *sdl.Surface
}

func NewBoard(size BoardSize) (*Board, os.Error) {
	b := new(Board)

	b.size = int(size)
	b.pieces = make([]*Piece, b.size*b.size)

	b.bg = sdl.Load(path.Join(BoardPath, fmt.Sprintf("%v.png", b.size)))
	if b.bg == nil {
		return nil, os.NewError(sdl.GetError())
	}

	//b.img = sdl.CreateRGBSurface(sdl.HWSURFACE,
	//	int(b.bg.W),
	//	int(b.bg.H),
	//	int(b.bg.Format.BitsPerPixel),
	//	b.bg.Format.Rmask,
	//	b.bg.Format.Gmask,
	//	b.bg.Format.Bmask,
	//	b.bg.Format.Amask,
	//)
	//if b.img == nil {
	//	return nil, os.NewError(sdl.GetError())
	//}

	b.img = sdl.Load(path.Join(BoardPath, fmt.Sprintf("%v.png", b.size)))
	if b.img == nil {
		return nil, os.NewError(sdl.GetError())
	}

	runtime.SetFinalizer(b, (*Board).free)

	return b, nil
}

func (b *Board) free() {
	b.bg.Free()
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

func (b *Board) drawPiece(x, y int, p *Piece) {
	pimg := p.Image()

	switch BoardSize(b.size) {
	case Size19x19:
		x = (x * 25) + 14
		y = (y * 25) + 14
	}

	x -= int(pimg.W / 2)
	y -= int(pimg.H / 2)

	b.img.Blit(&sdl.Rect{X: int16(x), Y: int16(y)}, pimg, nil)
}

func (b *Board) Image() *sdl.Surface {
	b.img.FillRect(nil, sdl.MapRGB(b.img.Format, 0, 0, 0))

	b.img.Blit(nil, b.bg, nil)

	for y := 0; y < b.size; y++ {
		for x := 0; x < b.size; x++ {
			if p := b.At(x, y); p != nil {
				b.drawPiece(x, y, p)
			}
		}
	}

	return b.img
}
