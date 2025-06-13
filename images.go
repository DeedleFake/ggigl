package main

import (
	"embed"
	"image"
	"image/color"
	"image/png"
	"path"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	//go:embed data
	data embed.FS

	board9x9Image   = loadImage("boards/9")
	board19x19Image = loadImage("boards/19")

	pieceBlackImage = loadImage("pieces/black")
	pieceWhiteImage = loadImage("pieces/white")

	blackHighlight = image.NewUniform(&color.NRGBA{0, 0, 0, 128})
	whiteHighlight = image.NewUniform(&color.NRGBA{255, 255, 255, 128})
)

func loadImage(name string) *ebiten.Image {
	file, err := data.Open(path.Join("data", name) + ".png")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	img, err := png.Decode(file)
	if err != nil {
		panic(err)
	}

	return ebiten.NewImageFromImage(img)
}
