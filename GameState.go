package main 

import (
	"time"
	"math/rand"

	"github.com/rivo/tview"
	"github.com/gdamore/tcell/v2"
)

type GameState struct {
	frame int
	updateRate time.Duration
	bag *PieceBag

	// Piece field
	spawnPosition Coordinate
	width int
	height int
	field []rune
	displayField []rune
	current *Tetromino

	frozenFor float64  // Freeze frame
	pieceAdvanceDelay int
	lockDelay int

	upNext []Tetromino

	held *Tetromino
	switchedHeld bool
}

func NewGameState() *GameState {
	state := GameState{}
	state.frame = 0

	// would rather i could make these static or constant.
	state.spawnPosition = Coordinate{4, 0}
	state.updateRate = 1000 / 60 * time.Millisecond
	state.width = 12
	state.height = 20

	state.bag = NewPieceBag(&state)
	state.current = state.RandomPiece()
	state.current.position = state.spawnPosition

	state.held = nil
	state.upNext = nil
	state.switchedHeld = false

	state.frozenFor = 0
	state.pieceAdvanceDelay = state.GetSpeed()
	state.lockDelay = 60

	fieldContent := make([]rune, state.width * state.height)
	for index := range fieldContent {
		fieldContent[index] = ' '
	}
	state.field = fieldContent

	return &state
}

func (state *GameState) InputCapture(event *tcell.EventKey) *tcell.EventKey {
	// go's switch statement looks to be just an if statement anyway
	if event.Rune() == 'a' {

	} else if event.Rune() == 'd' {
		
	} else if event.Rune() == 'w' {  // Hard drop
		
	} else if event.Rune() == 's' {  // Soft drop
		state.current.Fall()
		state.pieceAdvanceDelay = state.GetSpeed()
	} else if event.Rune() == ' ' {  // hold.
		
	} else if event.Rune() == 'e' {
		
	} else if event.Key() == tcell.KeyLeft {

	} else if event.Key() == tcell.KeyRight {

	} else if event.Rune() == 'r' {

	}
	return nil
}

func (state *GameState) UpdateLoop(field_text_view *tview.TextView) {
	for {
		state.UpdateField(field_text_view)
		// TODO: Account for duration of frame (figure out how tf time module works)
		time.Sleep(state.updateRate)
	}
}

func (state *GameState) UpdateField(text_field *tview.TextView) {
	state.displayField = make([]rune, len(state.field))
	copy(state.displayField, state.field)

	if (!state.current.DoesPieceFit(Coordinate{0, 1})) {
		if state.lockDelay > 0 {
			state.lockDelay--
		} else {
			// TODO: Next piece.
		}
	}

	state.pieceAdvanceDelay--
	if state.pieceAdvanceDelay == 0 {
		state.current.Fall()
		state.pieceAdvanceDelay = state.GetSpeed()
	}
	state.current.Render()

	// field_content operations
	state.field[20 * 12 - 1] = 'S'
	state.field[20 * 12 - 2] = 'Z'
	state.field[20 * 12 - 3] = 'J'
	state.field[20 * 12 - 4] = 'L'
	state.field[20 * 12 - 5] = 'I'
	state.field[20 * 12 - 6] = 'O'
	state.field[20 * 12 - 7] = 'T'

	state.RenderField(text_field)
}

func (state *GameState) RenderField(text_field *tview.TextView) {
	/*
	The documentation says that this lists the valid colour names:
	https://www.w3schools.com/colors/colors_names.asp
	But multiple I've tried simply don't work.
	Make sure to double check colours used.

	markup reference:

	[yellow]Yellow text
	[yellow:red]Yellow text on red background
	[:red]Red background, text color unchanged
	[yellow::u]Yellow text underlined
	[::bl]Bold, blinking text
	[::-]Colors unchanged, flags reset
	[-]Reset foreground color
	[-:-:-]Reset everything
	[:]No effect
	[]Not a valid color tag, will print square brackets as they are
	*/
	var fieldWithBorder string
	// This stops us relying on text wrapping, as that breaks when using box borders.
	for y := 0; y < state.height; y++ {
		for x := 0; x < state.width; x++ {
			index := x + (y * state.width)
			fieldWithBorder += string(state.displayField[index])
		}
		fieldWithBorder += "\n"
	}

	// Colour field.
	var content string
	for index := range fieldWithBorder {
		currChar := string(fieldWithBorder[index])
		// content += "[yellow:red]" + currChar
		// Not efficient to do this for each character, but not the point of this code.
		switch (currChar) {
		case "S":
			content += "[:#800000]"
		case "Z":
			content += "[:#008080]"
		case "L":
			content += "[:#008000]"
		case "J":
			content += "[:#800080]"
		case "O":
			content += "[:#000080]"
		case "I":
			content += "[:#808000]"
		case "T":
			content += "[:#FF00FF]"
		default:
			content += "[-:-:-]"
		}

		content += currChar
	}
	text_field.SetText(content)
}

func (state *GameState) GetSpeed() int {
	// Calculate based on score/level/settings
	return 60
}

func (state *GameState) LockPiece() {
	state.lockDelay = 60
	state.pieceAdvanceDelay = state.GetSpeed()

	// Add piece to field
	for y := 0; y < 4; y++ {
		for x := 0; x < 4; x++ {
			index := x + (4 * y)
			symbol := rune(state.current.shape[index])
			if symbol == '.' {
				continue
			}
			
			fieldIndex := (state.current.position.x + x) + (state.current.position.y + y) * state.width
			state.field[fieldIndex] = state.current.character
		}
	}

	// Get next piece
	state.NextPiece();

	// TODO: Gameover
}

func (state *GameState) NextPiece() {
	state.current = state.bag.TakeTopPiece()
}

func (state *GameState) RandomPiece() *Tetromino {
	slice := []rune {'I', 'O', 'S', 'Z', 'J', 'L', 'T'}
	char := slice[rand.Intn(6)]
	// for i := range slice {
	// 	j := rand.Intn(i + 1)
	// 	slice[i], slice[j] = slice[j], slice[i]
	// }
	return NewTetromino(char, state)
}