package main

import (
	"math/rand"
)

type PieceBag struct {
	state *GameState
	pieces []Tetromino
	minSize int
}

func NewPieceBag(state *GameState) *PieceBag {
	bag := PieceBag{}
	bag.minSize = 7
	bag.pieces = make([]Tetromino, 0, bag.minSize * 2)
	for i := 0; i < bag.minSize; i++ {
		
	}

	return &bag
}

func (bag *PieceBag) TakeTopPiece() *Tetromino {
	x, t := bag.pieces[0], bag.pieces[1:]
	if len(t) < 7 {
		bag.GetNewPieceStack()
	}
	return &x
}

func (bag *PieceBag) GetNewPieceStack() {
	chars := []rune {'O', 'I', 'T', 'S', 'Z', 'J', 'L'}
	for i := range chars {
		j := rand.Intn(i + 1)
		chars[i], chars[j] = chars[j], chars[i]
	}

	for _, char := range chars {
		bag.pieces = append(bag.pieces, *NewTetromino(char, bag.state))
	}
}