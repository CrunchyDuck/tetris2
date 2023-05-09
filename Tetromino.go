package main

// In lieu of constant arrays/dictionaries/etc, I'm using a function.
func GetPieceShape(char rune) string {
	switch(char) {
	case 'I':
		return "...." +
		"XXXX" +
		"...." +
		"...."
	case 'O':
		return ".XX." +
		".XX." +
		"...." +
		"...."
	case 'L':
		return "..X." +
		"XXX." +
		"...." +
		"...."
	case 'J':
		return "X..." +
		"XXX." +
		"...." +
		"...."
	case 'Z':
		return "XX.." +
		".XX." +
		"...." +
		"...."
	case 'S':
		return ".XX." +
		"XX.." +
		"...." +
		"...."
	case 'T':
		return ".X.." +
		"XXX." +
		"...." +
		"...."
	default:
		panic("No piece with symbol " + string(char))
	}
}

type Tetromino struct {
	position Coordinate  // Top left of piece.
	character rune
	rotation int
	state *GameState
	shape string
}

func NewTetromino(char rune, state *GameState) *Tetromino {
	t := Tetromino{}
	t.position = Coordinate{-1, -1}
	t.character = char
	t.shape = GetPieceShape(char)
	t.rotation = 0
	t.state = state
	return &t
}

// Returns false if it was blocked.
func (piece *Tetromino) Fall() bool {
	piece.position.y += 1
	if !piece.DoesPieceFit(Coordinate{}) {
		piece.position.y -= 1
		return false
	}
	return true
}

func (piece *Tetromino) Render() {
	for y := 0; y < 4; y++ {
		for x := 0; x < 4; x++ {
			index := x + (4 * y)
			symbol := rune(piece.shape[index])
			if symbol == '.' {
				continue
			}
			
			fieldIndex := (piece.position.x + x) + (piece.position.y + y) * piece.state.width
			piece.state.displayField[fieldIndex] = piece.character
		}
	}
}

func (piece *Tetromino) DoesPieceFit(offset Coordinate) bool {
	for y := 0; y < 4; y++ {
		for x := 0; x < 4; x++ {
			index := x + (4 * y)
			symbol := rune(piece.shape[index])
			if symbol == '.' {
				continue
			}
			
			fieldIndex := (piece.position.x + offset.x + x) + (piece.position.y + offset.y + y) * piece.state.width
			if piece.state.field[fieldIndex] == ' ' {
				continue
			}
			return false
		}
	}
	return true
}

// don't look at this please
func (piece *Tetromino) GetPieceKickOffsets(rotate_clockwise bool) [4]Coordinate {
	// Based on: https://tetris.wiki/Super_Rotation_System#Wall_Kicks
	// Instead of 0, R, L and 2, I'm using 0-3, where 0 is spawn and rotation is clockwise.
	// Therefore, 0 = 0, R = 1, 2 = 2, L = 3

	// can't put this in a const, can't put it in an object, guess i'm putting it here
	// Still haven't figured out the method to this madness.
	wallKicks := [4][4]Coordinate {
		// 0>1 or 2>1
		{
			Coordinate {-1, 0},
			Coordinate {-1, -1},
			Coordinate {0, 2},
			Coordinate {-1, 2},
		},
		// 1>0 or 1>2
		{
			Coordinate {1, 0},
			Coordinate {1, 1},
			Coordinate {1, -2},
			Coordinate {1, -2},
		},
		// 2>3 or 0>3
		{
			Coordinate {1, 0},
			Coordinate {1, -1},
			Coordinate {0, 2},
			Coordinate {1, 2},
		},
		// 3>2 or 3>0
		{
			Coordinate {-1, 0},
			Coordinate {-1, 1},
			Coordinate {0, -2},
			Coordinate {-1, -2},
		},
	}
	wallKicksI := [4][4]Coordinate {
		// 0>1 or 3>2
		{
			Coordinate {-2, 0},
			Coordinate {1, 0},
			Coordinate {-2, 1},
			Coordinate {1, -2},
		},
		// 1>2 or 0>3
		{
			Coordinate {-1, 0},
			Coordinate {2, 0},
			Coordinate {-1, -2},
			Coordinate {2, 1},
		},
		// 1>0 or 2>3
		{
			Coordinate {2, 0},
			Coordinate {-1, 0},
			Coordinate {2, -1},
			Coordinate {-1, 2},
		},
		// 2>1 or 3>0
		{
			Coordinate {1, 0},
			Coordinate {-2, 0},
			Coordinate {1, 2},
			Coordinate {-2, -1},
		},
	}

	if rotate_clockwise {
		if piece.character != 'i' {
			switch piece.rotation {
			case 0:
				return wallKicks[3]
			case 1:
				return wallKicks[0]
			case 2:
				return wallKicks[1]
			default:
				return wallKicks[2]
			}
		}
		// i piece
		switch piece.rotation {
		case 0:
			return wallKicksI[3]
		case 1:
			return wallKicksI[0]
		case 2:
			return wallKicksI[1]
		default:
			return wallKicksI[2]
		}
	}
	if piece.character != 'i' {
		switch piece.rotation {
		case 0:
			return wallKicks[1]
		case 1:
			return wallKicks[0]
		case 2:
			return wallKicks[3]
		default:
			return wallKicks[2]
		}
	}
	switch piece.rotation {
	case 0:
		return wallKicks[2]
	case 1:
		return wallKicks[3]
	case 2:
		return wallKicks[0]
	default:
		return wallKicks[1]
	}
}