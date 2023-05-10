package main 

import (
	"fmt"
	"time"
	"strings"
	"sort"

	"github.com/rivo/tview"
	"github.com/gdamore/tcell/v2"
)

type GameState struct {
	programState *ProgramState
	frame int
	updateRate time.Duration
	stopGame bool
	gameOver bool

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

	// Level/Score
	levelView *tview.TextView
	linesCleared int
	level int
	score int
	combo int
	wasLastInputRotation bool

	// Score log
	scoreView *tview.TextView
	// scoreWidth int  // Defined in the UI code instead.
	scoreHeight int
	scoreLog map[Score]*ScoreEvent
}

func NewGameState(state *ProgramState, field_view, up_next_view, held_view, score_view, level_view *tview.TextView) *GameState {
	this := GameState{}
	this.fieldView = field_view
	this.upNextView = up_next_view
	this.frame = 0
	this.programState = state
	this.stopGame = false
	this.gameOver = false

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

	this.levelView = level_view
	this.linesCleared = 0
	this.level = 1
	this.score = 0
	this.combo = 0
	this.wasLastInputRotation = false

	this.scoreView = score_view
	// this.scoreWidth = 5
	this.scoreHeight = 5
	this.scoreLog = make(map[Score]*ScoreEvent)

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

	this.DebugInitPerfectTriple()
	return &this
}

func (this *GameState) DebugInitTSpinTriple() {
	this.field[len(this.field)-1] = 'O'
	this.field[len(this.field)-2] = 'O'
	this.field[len(this.field)-3] = 'O'
	this.field[len(this.field)-4] = 'O'
	this.field[len(this.field)-5] = 'O'
	this.field[len(this.field)-6] = 'O'
	this.field[len(this.field)-7] = 'O'
	this.field[len(this.field)-8] = 'O'

	this.field[len(this.field)-10] = 'O'

	this.field[len(this.field)-11] = 'O'
	this.field[len(this.field)-12] = 'O'
	this.field[len(this.field)-13] = 'O'
	this.field[len(this.field)-14] = 'O'
	this.field[len(this.field)-15] = 'O'
	this.field[len(this.field)-16] = 'O'
	this.field[len(this.field)-17] = 'O'


	this.field[len(this.field)-20] = 'O'
	
	this.field[len(this.field)-21] = 'O'
	this.field[len(this.field)-22] = 'O'
	this.field[len(this.field)-23] = 'O'
	this.field[len(this.field)-24] = 'O'
	this.field[len(this.field)-25] = 'O'
	this.field[len(this.field)-26] = 'O'
	this.field[len(this.field)-27] = 'O'
	this.field[len(this.field)-28] = 'O'

	this.field[len(this.field)-30] = 'O'

	this.field[len(this.field)-40] = 'O'

	this.field[len(this.field)-49] = 'O'
	this.field[len(this.field)-50] = 'O'
}

func (this *GameState) DebugInitPerfectTetris() {
	this.field[len(this.field)-1] = 'O'
	this.field[len(this.field)-2] = 'O'
	this.field[len(this.field)-3] = 'O'
	this.field[len(this.field)-4] = 'O'
	this.field[len(this.field)-5] = 'O'
	this.field[len(this.field)-6] = 'O'
	this.field[len(this.field)-7] = 'O'
	this.field[len(this.field)-8] = 'O'
	this.field[len(this.field)-9] = 'O'

	this.field[len(this.field)-11] = 'O'
	this.field[len(this.field)-12] = 'O'
	this.field[len(this.field)-13] = 'O'
	this.field[len(this.field)-14] = 'O'
	this.field[len(this.field)-15] = 'O'
	this.field[len(this.field)-16] = 'O'
	this.field[len(this.field)-17] = 'O'
	this.field[len(this.field)-18] = 'O'
	this.field[len(this.field)-19] = 'O'

	this.field[len(this.field)-21] = 'O'
	this.field[len(this.field)-22] = 'O'
	this.field[len(this.field)-23] = 'O'
	this.field[len(this.field)-24] = 'O'
	this.field[len(this.field)-25] = 'O'
	this.field[len(this.field)-26] = 'O'
	this.field[len(this.field)-27] = 'O'
	this.field[len(this.field)-28] = 'O'
	this.field[len(this.field)-29] = 'O'

	this.field[len(this.field)-31] = 'O'
	this.field[len(this.field)-32] = 'O'
	this.field[len(this.field)-33] = 'O'
	this.field[len(this.field)-34] = 'O'
	this.field[len(this.field)-35] = 'O'
	this.field[len(this.field)-36] = 'O'
	this.field[len(this.field)-37] = 'O'
	this.field[len(this.field)-38] = 'O'
	this.field[len(this.field)-39] = 'O'
}

func (this *GameState) DebugInitPerfectTriple() {
	this.field[len(this.field)-1] = 'O'
	this.field[len(this.field)-2] = 'O'
	this.field[len(this.field)-3] = 'O'
	this.field[len(this.field)-4] = 'O'
	this.field[len(this.field)-5] = 'O'
	this.field[len(this.field)-6] = 'O'
	this.field[len(this.field)-7] = 'O'
	this.field[len(this.field)-8] = 'O'
	this.field[len(this.field)-9] = 'O'

	this.field[len(this.field)-11] = 'O'
	this.field[len(this.field)-12] = 'O'
	this.field[len(this.field)-13] = 'O'
	this.field[len(this.field)-14] = 'O'
	this.field[len(this.field)-15] = 'O'
	this.field[len(this.field)-16] = 'O'
	this.field[len(this.field)-17] = 'O'
	this.field[len(this.field)-18] = 'O'
	this.field[len(this.field)-19] = 'O'

	this.field[len(this.field)-21] = 'O'
	this.field[len(this.field)-22] = 'O'
	this.field[len(this.field)-23] = 'O'
	this.field[len(this.field)-24] = 'O'
	this.field[len(this.field)-25] = 'O'
	this.field[len(this.field)-26] = 'O'
	this.field[len(this.field)-27] = 'O'
	this.field[len(this.field)-28] = 'O'
}

func (this *GameState) InputCapture(event *tcell.EventKey) *tcell.EventKey {
	// go's switch statement looks to be just an if statement anyway
	if event.Rune() == 'a' {
		if (this.current.MoveLeft()) {
			this.wasLastInputRotation = false
		}
	} else if event.Rune() == 'd' {
		if (this.current.MoveRight()) {
			this.wasLastInputRotation = false
		}
	} else if event.Rune() == 'w' {  // Hard drop
		for this.current.Fall() {
			this.wasLastInputRotation = false
			this.AddScore(scoreHardDrop)
		}
		this.LockPiece()
	} else if event.Rune() == 's' {  // Soft drop
		if (this.current.Fall()) {
			this.wasLastInputRotation = false
			this.AddScore(scoreSoftDrop)
			this.pieceAdvanceDelay = this.GetSpeed()
		}
	} else if event.Rune() == ' ' && !this.switchedHeld {  // hold.
		this.wasLastInputRotation = false
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
			this.wasLastInputRotation = true
		}
	} else if event.Key() == tcell.KeyRight {
		if this.current.RotateRight() {
			this.LockCountdownReset()
			this.wasLastInputRotation = true
		}
	} else if event.Rune() == 'r' {
		this.stopGame = true
	} else if event.Rune() == 'q' {
		this.programState.ascii = !this.programState.ascii
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
		} else if this.gameOver {
			this.programState.gracefulExit = true
			this.programState.targetState = programGameOver
			this.programState.app.Stop()
			break
		} else {
			this.UpdateField()
		}

		this.Render()
		// TODO: Account for duration of frame (figure out how tf time module works)
		time.Sleep(this.updateRate)
		this.frame++
		if this.stopGame {
			this.programState.gracefulExit = true
			this.programState.targetState = programMenu
			this.programState.app.Stop()
			break
		}
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
		if (this.current.Fall()) {
			this.pieceAdvanceDelay = this.GetSpeed()
			this.wasLastInputRotation = false
		}
	}
}

func (this *GameState) Render() {
	this.RenderField()
	this.RenderHeld()
	this.RenderUpNext()
	this.RenderScoreLog()
	this.RenderLevel()
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
	for i := 0; i < 4; i++ {
		t := this.bag.pieces[i]
		tShape := t.GetShapeString()
		for y := 0; y < 4; y++ {
			body += "\n "
			for x := 0; x < t.size; x++ {
				index := y * t.size + x
				if index >= t.size * 2 {
					break
				}
				body += string(tShape[index])
			}
		}
	}

	content := header + this.ApplyTetrominoColors(body)
	this.upNextView.SetText(content)
}

func (this *GameState) RenderScoreLog() {
	/// why isn't there a builtin way to convert a map to a slice
	/// am i just missing it, SO seems to say there isn't
	scoreEvents := make([]*ScoreEvent, 0, len(this.scoreLog))
	for  _, value := range this.scoreLog {
		scoreEvents = append(scoreEvents, value)
	}

	sort.Slice(scoreEvents, func(i, j int) bool {
		return scoreEvents[i].totalScore < scoreEvents[j].totalScore
	})

	var body string
	for _, scoreEvent := range scoreEvents {
		scoreEvent.disappearDelay--
		if (scoreEvent.disappearDelay <= 0) {
			delete(this.scoreLog, scoreEvent.scoreType)
			continue
		}

		body += fmt.Sprintf("%s x%d = %d\n", scoreEvent.name, scoreEvent.count, scoreEvent.totalScore)
	}

	this.scoreView.SetText(body)
}

func (this *GameState) RenderLevel() {
	body := fmt.Sprintf("LEVEL:\n%d\n\nSCORE:\n%d\n\nLINES:\n%d", this.level, this.score, this.linesCleared)
	this.levelView.SetText(body)
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

		if currChar == "\n" {
			content += "\n"
		} else if this.programState.ascii {
			content += currChar
		} else {
			content += " "
		}
	}

	return content
}

func (this *GameState) GetSpeed() int {
	// Based on: https://tetris.fandom.com/wiki/Tetris_(NES,_Nintendo)
	// Calculate based on score/level/settings
	switch this.level {
	case 1, 2, 3, 4, 5, 6, 7, 8, 9:
		return 48 - ((this.level - 1) * 5)
	case 10:
		return 6
	case 11, 12, 13:
		return 5
	case 14, 15, 16:
		return 4
	case 17, 18, 19:
		return 3
	case 20, 21, 22, 23, 24, 25, 26, 27, 28, 29:
		return 2
	default:
		return 1
	}
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

	if !this.current.DoesPieceFit(Coordinate{}) {
		this.gameOver = true
	}
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

	totalLines := 0  // used to check for perfect clear
	clearLineCount := 0
	for y := 0; y < this.fieldHeight; y++ {
		clearedLine := true
		lineHadPieces := false
		for x := 0; x < this.fieldWidth; x++ {
			index := x + (y * this.fieldWidth)
			char := this.field[index]

			switch char {
			case ' ':
				clearedLine = false
			default:
				lineHadPieces = true
			}
		}
		if lineHadPieces {
			totalLines++
		}
		if !clearedLine {
			continue
		}

		// Replace with "cleared line" symbol.
		clearLineCount++
		startIndex := y * this.fieldWidth
		endIndex := startIndex + this.fieldWidth
		copy(this.field[startIndex:endIndex], clearedLineText)
	}
	if clearLineCount == 0 {
		this.combo = 0
		return
	}

	wasTSpin := this.CheckTSpinValidity()
	
	if clearLineCount == totalLines {  // Perfect clear
		switch clearLineCount {
		case 1:
			this.AddScore(scorePerfectSingle)
		case 2:
			this.AddScore(scorePerfectDouble)
		case 3:
			this.AddScore(scorePerfectTriple)
		case 4:
			this.AddScore(scorePerfectQuad)
		default:
			return
		}
	} 
	if wasTSpin {
		switch clearLineCount {
		case 1:
			this.AddScore(scoreTSpinSingle)
		case 2:
			this.AddScore(scoreTSpinDouble)
		case 3:
			this.AddScore(scoreTSpinTriple)
		default:
			return
		}
	} else if clearLineCount != totalLines {  // Normal clear
		switch clearLineCount {
		case 1:
			this.AddScore(scoreLineSingle)
		case 2:
			this.AddScore(scoreLineDouble)
		case 3:
			this.AddScore(scoreLineTriple)
		case 4:
			this.AddScore(scoreLineQuad)
		default:
			return
		}	
	}

	this.combo++
	if this.combo > 1 {
		this.AddScore(scoreCombo)
	}

	// Increase level
	this.linesCleared += clearLineCount
	this.level = (this.linesCleared / 10) + 1  // level is 1 indexed for score math :)

	// Freeze on frame for a moment.
	this.freezeFrameDuration = this.freezeFrameMax
}

func (this *GameState) CheckTSpinValidity() bool {
	// Based on: https://tetris.fandom.com/wiki/T-Spin
	// Must be a T, that was rotated
	if !this.wasLastInputRotation || this.current.character != 'T' {
		return false
	}

	// Piece must be immobile
	if this.current.DoesPieceFit(Coordinate{1, 0}) ||
	this.current.DoesPieceFit(Coordinate{-1, 0}) ||
	this.current.DoesPieceFit(Coordinate{0, -1}) {
		return false
	}

	// Must have corners locked.
	x := this.current.position.x
	y := this.current.position.y
	freeSpaces := 0
	for xOff := 0; xOff < 2; xOff++ {
		for yOff := 0; yOff < 2; yOff++ {
			if (this.CellAvailable(x + (xOff * 2), y + (yOff) * 2)) {
				freeSpaces++
			}
		}
	}
	if freeSpaces > 1 {
		return false
	}
	return true
	
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
		log = NewScoreEvent(score_type)
		this.scoreLog[score_type] = log
	}
	
	this.score += log.IncreaseScore(this.level)
}

func (this *GameState) CellAvailable(x_pos, y_pos int) bool {
	fieldIndex := x_pos + (y_pos * this.current.state.fieldWidth)

	if x_pos >= this.current.state.fieldWidth || x_pos < 0 {  // x bounds check
		return false
	} else if y_pos >= this.current.state.fieldHeight || y_pos < 0 {  // y bounds check
		return false
	} else if this.current.state.field[fieldIndex] != ' ' { // Collision check
		return false
	}

	return true
}