package main 

import (
	"time"
	"strings"

	"github.com/rivo/tview"
	"github.com/gdamore/tcell/v2"
)

type GameState struct {
	frame int
	updateRate time.Duration
	displayLetters bool

	// Piece field
	fieldView *tview.TextView
	spawnPosition Coordinate
	fieldWidth int
	fieldHeight int
	field []rune
	displayField []rune
	current *Tetromino

	pieceAdvanceDelay int
	lockDelay int
	maxLockDelay int
	lockResets int
	maxLockResets int

	freezeFrameDuration int
	freezeFrameMax int
	freezeFrameFlashRate int

	// Up next
	upNextView *tview.TextView
	upNextWidth int
	upNextHeight int
	bag *PieceBag
	
	// Held
	heldView *tview.TextView
	heldWidth int
	heldHeight int
	held *Tetromino
	switchedHeld bool

	// Score
	scoreView *tview.TextView
	score int
	scoreWidth int
	scoreHeight int
	scoreLog map[Score]ScoreEvent
}

func NewGameState(field_view, up_next_view, held_view, score_view *tview.TextView) *GameState {
	this := GameState{}
	this.fieldView = field_view
	this.upNextView = up_next_view
	this.frame = 0

	// would rather i could make these static or constant.
	this.spawnPosition = Coordinate{4, 0}
	this.updateRate = 1000 / 60 * time.Millisecond
	this.fieldWidth = 10
	this.fieldHeight = 20

	this.upNextView = up_next_view
	this.upNextWidth = 6
	this.upNextHeight = 20

	this.held = nil
	this.heldView = held_view
	this.heldWidth = 6
	this.heldHeight = 6
	this.switchedHeld = false

	this.scoreView = score_view
	this.score = 0
	this.scoreWidth = 10
	this.scoreHeight = 10
	this.scoreLog = make(map[Score]ScoreEvent)

	this.bag = NewPieceBag(&this)
	this.NextPiece()

	this.freezeFrameDuration = 0
	this.freezeFrameMax = 60
	this.freezeFrameFlashRate = 15
	this.pieceAdvanceDelay = this.GetSpeed()

	this.maxLockDelay = 60
	this.maxLockResets = 5

	this.lockDelay = this.maxLockDelay
	this.lockResets = this.maxLockResets

	fieldContent := make([]rune, this.fieldWidth * this.fieldHeight)
	for index := range fieldContent {
		fieldContent[index] = ' '
	}
	this.field = fieldContent

	this.displayLetters = true

	return &this
}

func (this *GameState) InputCapture(event *tcell.EventKey) *tcell.EventKey {
	// go's switch statement looks to be just an if statement anyway
	if event.Rune() == 'a' {
		this.current.MoveLeft()
	} else if event.Rune() == 'd' {
		this.current.MoveRight()
	} else if event.Rune() == 'w' {  // Hard drop
		for this.current.Fall() {
			this.AddScore(scoreHardDrop)
		}
		this.LockPiece()
	} else if event.Rune() == 's' {  // Soft drop
		this.current.Fall()
		this.AddScore(scoreSoftDrop)
		this.pieceAdvanceDelay = this.GetSpeed()
	} else if event.Rune() == ' ' && !this.switchedHeld {  // hold.
		if this.held != nil {
			t := this.held
			this.held = this.current
			this.current = t
		} else {
			this.held = this.current
			this.current = nil
			this.NextPiece()
		}
		this.SetPieceDefaults()
		this.switchedHeld = true
	} else if event.Key() == tcell.KeyLeft {
		if this.current.RotateLeft() {
			this.LockCountdownReset()
		}
	} else if event.Key() == tcell.KeyRight {
		if this.current.RotateRight() {
			this.LockCountdownReset()
		}
	} else if event.Rune() == 'r' {
		// *state = *NewGameState()
	} else if event.Rune() == 'q' {
		this.displayLetters = !this.displayLetters
	}
	return nil
}

func (this *GameState) UpdateLoop() {
	for {
		// Flash lines
		if this.freezeFrameDuration > 0 {
			this.freezeFrameDuration--
			if this.freezeFrameDuration == 0 {
				this.RemoveLineClears()
			}
		} else {
			this.UpdateField()
		}

		this.Render()
		// TODO: Account for duration of frame (figure out how tf time module works)
		time.Sleep(this.updateRate)
		this.frame++
	}
}

func (this *GameState) UpdateField() {
	// Update piece
	if (!this.current.DoesPieceFit(Coordinate{0, 1})) {
		if this.lockDelay > 0 {
			this.lockDelay--
		} else {
			this.LockPiece()
		}
	} else {
		this.LockCountdownReset()
	}
	this.pieceAdvanceDelay--
	if this.pieceAdvanceDelay == 0 {
		this.current.Fall()
		this.pieceAdvanceDelay = this.GetSpeed()
	}
}

func (this *GameState) Render() {
	this.RenderField()
	this.RenderHeld()
	this.RenderUpNext()
}

func (this *GameState) RenderField() {
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
	
	this.displayField = make([]rune, len(this.field))
	copy(this.displayField, this.field)
	
	this.current.Render()

	var fieldWithBorder string
	// This stops us relying on text wrapping, as that breaks when using box borders.
	for y := 0; y < this.fieldHeight; y++ {
		for x := 0; x < this.fieldWidth; x++ {
			index := x + (y * this.fieldWidth)
			fieldWithBorder += string(this.displayField[index])
		}
		fieldWithBorder += "\n"
	}

	// Colour field.
	this.fieldView.SetText(this.ApplyTetrominoColors(fieldWithBorder))
}

func (this *GameState) RenderHeld() {
	heldContent := " HELD \n"

	if this.held == nil {
		this.heldView.SetText(heldContent)
		return
	}

	var heldPieceText string
	heldPieceShape := this.held.GetShapeString()
	for y := 0; y < this.held.size; y++ {
		heldPieceText += "\n "
		for x := 0; x < this.held.size; x++ {
			index := y * this.held.size + x
			heldPieceText += string(heldPieceShape[index])
		}
	}
	heldContent += this.ApplyTetrominoColors(heldPieceText)
	this.heldView.SetText(heldContent)
}

func (this *GameState) RenderUpNext() {
	header := " NEXT \n"
	var body string
	for i := 0; i < 3; i++ {
		t := this.bag.pieces[i]
		tShape := t.GetShapeString()
		for y := 0; y < t.size; y++ {
			body += "\n "
			for x := 0; x < t.size; x++ {
				index := y * t.size + x
				body += string(tShape[index])
			}
		}
	}

	content := header + this.ApplyTetrominoColors(body)
	this.upNextView.SetText(content)
}

func (this *GameState) RenderScore() {
	header := " SCORE " + "\n\n"
	
	var body string
	for _, val := range this.scoreLog {
		body += val.GetText() + "\n"
	}

	this.scoreView.SetText(header + body)
}

func (this *GameState) ApplyTetrominoColors(input string) string {
	var content string
	for index := range input {
		currChar := string(input[index])
		// Not efficient to do this for each character, could add chains.
		switch (currChar) {
		case "S":
			content += "[:#800000]"
		case "Z":
			content += "[:#00C0C0]"
		case "L":
			content += "[:#008000]"
		case "J":
			content += "[:#8000C0]"
		case "O":
			content += "[:#000080]"
		case "I":
			content += "[:#A0A000]"
		case "T":
			content += "[:#FF00FF]"
		case "=":
			if (this.frame / this.freezeFrameFlashRate) % 2 == 0 {
				content += "[:#CCCCCC]"
			} else {
				content += "[-:-:-]"
			}
		default:
			content += "[-:-:-]"
		}

		if this.displayLetters {
			content += currChar
		} else {
			if currChar == "\n" {
				content += "\n"
			} else {
				content += " "
			}
		}
	}

	return content
}

func (this *GameState) GetSpeed() int {
	// Calculate based on score/level/settings
	return 60
}

func (this *GameState) LockCountdownReset() {
	if this.lockResets > 0 && this.lockDelay < this.maxLockDelay {
		this.lockDelay = this.maxLockDelay
		this.lockResets--
	}
}

func (this *GameState) LockPiece() {
	p := this.current
	// Add piece to field
	for y := 0; y < p.size; y++ {
		for x := 0; x < p.size; x++ {
			index := x + (p.size * y)
			symbol := rune(p.shape[index])
			if symbol == '.' {
				continue
			}
			
			fieldIndex := (this.current.position.x + x) + (this.current.position.y + y) * this.fieldWidth
			this.field[fieldIndex] = this.current.character
		}
	}

	this.TryLineClears()

	// Get next piece
	this.NextPiece()

	// TODO: Check gameover
}

func (this *GameState) SetPieceDefaults() {
	this.current.position = this.spawnPosition
	this.lockDelay = this.maxLockDelay
	this.lockResets = this.maxLockResets
	this.pieceAdvanceDelay = this.GetSpeed()
	this.switchedHeld = false
}

func (this *GameState) TryLineClears() {
	clearedLineText := []rune(strings.Repeat("=", this.fieldWidth))

	clearedAnyLines := false
	for y := 0; y < this.fieldHeight; y++ {
		clearedLine := true
		for x := 0; x < this.fieldWidth; x++ {
			index := x + (y * this.fieldWidth)
			if this.field[index] == ' ' {
				clearedLine = false
				break
			}
		}
		if !clearedLine {
			continue
		}

		clearedAnyLines = true
		startIndex := y * this.fieldWidth
		endIndex := startIndex + this.fieldWidth
		copy(this.field[startIndex:endIndex], clearedLineText)
	}

	if clearedAnyLines {
		this.freezeFrameDuration = this.freezeFrameMax
	}
}

func (this *GameState) RemoveLineClears() {
	blankLine := strings.Repeat(" ", this.fieldWidth)[:]

	for y := 0; y < this.fieldHeight; y++ {
		startIndex := y * this.fieldWidth
		endIndex := startIndex + this.fieldWidth
		// Remove line and shuffle other lines down.
		if this.field[y * this.fieldWidth] == '=' {
			newField := blankLine + string(this.field[:startIndex]) + string(this.field[endIndex:])
			this.field = []rune(newField)
		}
	}
}

func (this *GameState) NextPiece() {
	this.current = this.bag.TakeTopPiece()
	this.current.position = this.spawnPosition
	this.SetPieceDefaults()
}

func (this *GameState) AddScore(score_type Score) {
	log, hasLog := this.scoreLog[score_type]
	if !hasLog {
		log = *NewScoreEvent(score_type)
		this.scoreLog[score_type] = log
	}
	before := log.GetScore()
	log.count++
	after := log.GetScore()
	this.score += after - before
}