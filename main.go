package main

import (
	"fmt"
	"github.com/rivo/tview"
	"github.com/gdamore/tcell/v2"
)

func GenerateView(app *tview.Application) *tview.TextView {
	view := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWrap(false).
		SetChangedFunc(func() { app.Draw() })
	view.SetBorder(true)
	return view
}

func GenerateGame(program_state *ProgramState) tview.Primitive {

	// UI code my beloathed
	upNextField := GenerateView(program_state.app)
	tetrisField := GenerateView(program_state.app)	
	heldField := GenerateView(program_state.app)
	scoreField := GenerateView(program_state.app)
	levelField := GenerateView(program_state.app)
	state := NewGameState(program_state, tetrisField, upNextField, heldField, scoreField, levelField)
	// i don't like this :(
	gameWidth := state.heldWidth + 2 + state.fieldWidth + 2 + state.upNextWidth + 2

	heldFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(heldField, state.heldHeight + 2, 1, false).
		AddItem(levelField, 0, 1, false)
		
	topFlex := tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(heldFlex, state.heldWidth + 2, 1, false).
		AddItem(tetrisField, state.fieldWidth + 2, 1, false).
		AddItem(upNextField, state.upNextWidth + 2, 1, false).
		AddItem(nil, 0, 1, false)
	bottomFlex := tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(scoreField, gameWidth + 2, 1, false).
		AddItem(nil, 0, 1, false)

	gameFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(topFlex, state.fieldHeight + 2, 1, false).
		AddItem(bottomFlex, state.scoreHeight + 2, 1, false)
	

	// modal := tview.NewModal().SetText("Game over!\nPress r to restart")
	// pages := tview.NewPages().
	// 	AddPage("game", paddedFlex, false, true).
	// 	AddPage("gameOver", modal, false, false)
	

	paddedFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().
			AddItem(nil, 0, 1, false).
			AddItem(gameFlex, gameWidth, 1, false).
			AddItem(nil, 0, 1, false),
		state.scoreHeight + 2 + state.fieldHeight + 2, 1, false).
		AddItem(nil, 0, 1, false)

	
	paddedFlex.SetInputCapture(state.InputCapture)
	go state.UpdateLoop()
	return paddedFlex
}

func GenerateMenu(program_state *ProgramState) tview.Primitive {
	logoField := GenerateView(program_state.app)
	logoField.SetText(GetMenuText(program_state))

	paddedFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(nil, 0, 1, false).
		// Score field
		AddItem(tview.NewFlex().
			AddItem(nil, 0, 1, false).
			AddItem(logoField, 40, 1, false).
			AddItem(nil, 0, 1, false),
		30, 1, false).
		AddItem(nil, 0, 1, false)
	

	paddedFlex.SetInputCapture(func (event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'r' {
			program_state.targetState = programGame
			program_state.gracefulExit = true
			program_state.app.Stop()
		} else if event.Rune() == 'q' {
			program_state.ascii = !program_state.ascii
			logoField.SetText(GetMenuText(program_state))
		}
		return nil
	})
	
	return paddedFlex
}

func GenerateGameOver(program_state *ProgramState) tview.Primitive {
	field := GenerateView(program_state.app)
	field.SetText("GAME OVER!\nPress r to restart.\n\ni couldn't figure out how to overlay\nthis atop gameplay :)\n")

	paddedFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(nil, 0, 1, false).
		// Score field
		AddItem(tview.NewFlex().
			AddItem(nil, 0, 1, false).
			AddItem(field, 40, 1, false).
			AddItem(nil, 0, 1, false),
		10, 1, false).
		AddItem(nil, 0, 1, false)
	

	paddedFlex.SetInputCapture(func (event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'r' {
			program_state.targetState = programMenu
			program_state.gracefulExit = true
			program_state.app.Stop()
		}
		return nil
	})
	
	return paddedFlex
}

func GetLogo(program_state *ProgramState) string {
	logo := "" +
	"oxxooooooooooooooxxo" + "\n" +
	"oxxxxooooooooooxxxxo" + "\n" +
	"oxxxxxxo....oxxxxxxo" + "\n" +
	"ooxxxxx......xxxxxoo" + "\n" +
	"oooxxx........xxxooo" + "\n" +
	"ooooxx........xxoooo" + "\n" +
	"xx..o..........o..xx" + "\n" +
	"xxx..............xxx" + "\n" +
	"oxxx............xxxo" + "\n" +
	"oox....`....`....xoo" + "\n" +
	"ooox...`....`...xooo" + "\n" +
	"ooooo..........ooooo" + "\n" +
	"ooooox........xooooo" + "\n" +
	"ooooxx........xxoooo" + "\n" +
	"oooooo........oooooo" + "\n" +
	"ooooooo......ooooooo" + "\n" +
	"ooooooo.`..`.ooooooo" + "\n" +
	"oooooooo....oooooooo" + "\n"
	var content string
	for index := range logo {
		currChar := string(logo[index])
		// Not efficient to do this for each character, could add chains.
		switch (currChar) {
		case "x":
			content += "[:#BEC9DC]"
		case ".":
			content += "[:#A088C2]"
		case "`":
			content += "[:#45283C]"
		case "o":
			content += "[-:-:-]"
			currChar = " "
		}
		if currChar == "\n" {
			content += "\n"
		} else if program_state.ascii {
			content += currChar + currChar
		} else {
			content += "  "
		}
	}

	return content
}

func GetMenuText(program_state *ProgramState) string {
	body := GetLogo(program_state)
	body += "\n"
	body += "A D - Move\n"
	body += "S - Soft drop\n"
	body += "W - Hard drop\n"
	body += "space - Hold\n"
	body += "← → - Rotate\n"
	body += "r - Restart\n"
	body += "q - Toggle ascii (try it)\n"
	body += "\n"
	body += "All standard scoring works.\n"
	body += "Press r to start\n"
	return body
}

// Field
// Score
// Level
// Next pieces
// Hold piece
// Options menu?
func main() {
	var app *tview.Application
	app = tview.NewApplication()
	state := NewProgramState(app)

	// programWidth = 40
	// programHeight = 30
	
	for state.gracefulExit {
		state.gracefulExit = false
		var flex tview.Primitive
		// Run menu
		switch state.targetState {
		case programMenu:
			flex = GenerateMenu(state)
			break
		case programGame:
			flex = GenerateGame(state)
		case programGameOver:
			flex = GenerateGameOver(state)
		default:
			panic("Unknown state")
		}

		if flex == nil {
			fmt.Print("here")
			break
		}
		

		if err := app.SetRoot(flex, true).Run(); err != nil {
			panic(err)
		}
	}
}
