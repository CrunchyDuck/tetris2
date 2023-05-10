package main

// In lieu of constant arrays/dictionaries/etc, I'm using a function.
func GetPieceShape(char rune) (string, int) {
	switch(char) {
	case 'I':
		return ""+
		"...." +
		"XXXX" +
		"...." +
		"....", 4
	case 'O':
		return ""+
		"XX" +
		"XX", 2
	case 'L':
		return ""+
		"..X" +
		"XXX" +
		"...", 3
	case 'J':
		return ""+
		"X.." +
		"XXX" +
		"...", 3
	case 'Z':
		return ""+
		"XX." +
		".XX" +
		"...", 3
	case 'S':
		return ""+
		".XX" +
		"XX." +
		"...", 3
	case 'T':
		return ""+
		".X." +
		"XXX" +
		"...", 3
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
	size int
}

func NewTetromino(char rune, state *GameState) *Tetromino {
	t := Tetromino{}
	t.position = Coordinate{-1, -1}
	t.character = char
	t.shape, t.size = GetPieceShape(char)
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

func (piece *Tetromino) MoveRight() bool {
	piece.position.x += 1
	if !piece.DoesPieceFit(Coordinate{}) {
		piece.position.x -= 1
		return false
	}
	return true
}

func (piece *Tetromino) MoveLeft() bool {
	piece.position.x -= 1
	if !piece.DoesPieceFit(Coordinate{}) {
		piece.position.x += 1
		return false
	}
	return true
}

func (piece *Tetromino) RotateRight() bool {
	piece.Rotate(true)
	if piece.DoesPieceFit(Coordinate{}) {
		return true
	}

	kicks := piece.GetPieceKickOffsets(true)
	for _, kick := range kicks {
		kick.y = -kick.y  // My y is opposite to normal y.
		if piece.DoesPieceFit(kick) {
			piece.position.x += kick.x
			piece.position.y += kick.y
			return true
		}
	}
	
	// Return to original position
	piece.Rotate(false)
	return false
}

func (piece *Tetromino) RotateLeft() bool {
	piece.Rotate(false)
	if piece.DoesPieceFit(Coordinate{}) {
		return true
	}

	kicks := piece.GetPieceKickOffsets(false)
	for _, kick := range kicks {
		kick.y = -kick.y  // My y is opposite to normal y.
		if piece.DoesPieceFit(kick) {
			piece.position.x += kick.x
			piece.position.y += kick.y
			return true
		}
	}

	// Return to original position
	piece.Rotate(true)
	return false
}

func (piece *Tetromino) Rotate(clockwise bool) {
	rotated := ""
	s := piece.size
	
	for x := 0; x < s; x++ {
		for y := 0; y < s; y++ {
			var index int
			if clockwise {
				index = x + ((s - y - 1) * s)
			} else {
				index = (s - x - 1) + (y * s)
			}
			rotated += string(piece.shape[index])
		}
	}
	piece.shape = rotated
	
	if clockwise {
		piece.rotation++
		if piece.rotation >= s {
			piece.rotation = 0
		}
	} else {
		piece.rotation--
		if piece.rotation < 0 {
			piece.rotation = s - 1
		}
	}
}

func (piece *Tetromino) Render() {
	s := piece.size
	for y := 0; y < s; y++ {
		for x := 0; x < s; x++ {
			index := x + (s * y)
			symbol := rune(piece.shape[index])
			if symbol == '.' {
				continue
			}
			
			fieldIndex := (piece.position.x + x) + (piece.position.y + y) * piece.state.fieldWidth
			piece.state.displayField[fieldIndex] = piece.character
		}
	}
}

func (this *Tetromino) GetShapeString() string {
	ret := ""
	s := this.size
	for y := 0; y < s; y++ {
		for x := 0; x < s; x++ {
			index := x + (s * y)
			symbol := rune(this.shape[index])
			if symbol == '.' {
				ret += " "
			} else {
				ret += string(this.character)
			}
		}
	}
	return ret
}

func (piece *Tetromino) DoesPieceFit(offset Coordinate) bool {
	s := piece.size
	for y := 0; y < s; y++ {
		for x := 0; x < s; x++ {
			index := x + (s * y)
			symbol := rune(piece.shape[index])
			if symbol == '.' {
				continue
			}
			
			// Bounds checks
			xPos := piece.position.x + offset.x + x
			yPos := (piece.position.y + offset.y + y)
			fieldIndex := xPos + (yPos * piece.state.fieldWidth)
			if xPos >= piece.state.fieldWidth || xPos < 0 {
				return false
			} else if yPos >= piece.state.fieldHeight || yPos < 0 {
				return false
			} else if piece.state.field[fieldIndex] != ' ' { // Collision check
				return false
			}
		}
	}
	return true
}

// don't look at this please
func (piece *Tetromino) GetPieceKickOffsets(rotated_clockwise bool) [4]Coordinate {
	// Based on: https://tetris.wiki/Super_Rotation_System#Wall_Kicks
	// Instead of 0, R, L and 2, I'm using 0-3, where 0 is spawn and rotation is clockwise.
	// Therefore, 0 = 0, R = 1, 2 = 2, L = 3

	// can't put this in a const, can't put it in an object, guess i'm putting it here
	// Still haven't figured out the method to this madness.
	wallKicks := [4][4]Coordinate {
		// 0>1 or 2>1
		{
			Coordinate {-1, 0},
			Coordinate {-1, 1},
			Coordinate {0, -2},
			Coordinate {-1, -2},
		},
		// 1>0 or 1>2
		{
			Coordinate {1, 0},
			Coordinate {1, -1},
			Coordinate {0, 2},
			Coordinate {1, 2},
		},
		// 2>3 or 0>3
		{
			Coordinate {1, 0},
			Coordinate {1, 1},
			Coordinate {0, -2},
			Coordinate {1, -2},
		},
		// 3>2 or 3>0
		{
			Coordinate {-1, 0},
			Coordinate {-1, -1},
			Coordinate {0, 2},
			Coordinate {-1, 2},
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

	if rotated_clockwise {
		if piece.character != 'i' {
			switch piece.rotation {
			// 3>0
			case 0:
				return wallKicks[3]
			// 0>1
			case 1:
				return wallKicks[0]
			// 1>2
			case 2:
				return wallKicks[1]
			// 2>3
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
		// 1>0
		case 0:
			return wallKicks[1]
		// 2>1
		case 1:
			return wallKicks[0]
		// 3>2
		case 2:
			return wallKicks[3]
		// 0>3
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