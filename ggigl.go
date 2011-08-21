package main

var Pieces map[string]*Piece

func main() {
	Pieces = make(map[string]*Piece)
	Pieces["black"] = NewPiece("black")
	Pieces["white"] = NewPiece("white")

	println("There is no response...")
}
