package main

import (
	// "fmt"
	// "strconv"
	// "strings"
	// "time"
	// "math"

	"github.com/rivo/tview"
)


// Field
// Score
// Level
// Next pieces
// Hold piece
// Options menu?
func main() {
	var app *tview.Application
	app = tview.NewApplication()

	// UI code my beloathed
	tetrisField := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWrap(false).
		SetChangedFunc(func() { app.Draw() })
	tetrisField.SetBorder(true)
	
	upNextField := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWrap(false).
		SetChangedFunc(func() { app.Draw() })
	upNextField.SetBorder(true)

	heldField := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWrap(false).
		SetChangedFunc(func() { app.Draw() })
	heldField.SetBorder(true)

	scoreField := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWrap(false).
		SetChangedFunc(func() { app.Draw() })
	scoreField.SetBorder(true)

	state := NewGameState(tetrisField, upNextField, heldField, scoreField)


	heldFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(heldField, state.heldHeight + 2, 1, false)

	tetrisFlex := tview.NewFlex().
		AddItem(heldFlex, state.heldWidth + 2, 1, false).
		AddItem(tetrisField, state.fieldWidth + 2, 1, false).
		AddItem(upNextField, state.upNextWidth + 2, 1, false)
	
	mainFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(tetrisFlex, state.fieldHeight + 2, 1, false)
		// AddItem(scoreField, state.scoreHeight + 2, 1, false)
	_, _, mainFlexWidth, mainFlexHeight := mainFlex.GetRect()
	
	paddedFlexInner := tview.NewFlex().
		// AddItem(nil, 0, 1, false).
		AddItem(mainFlex, mainFlexWidth, 0, false)
		// AddItem(nil, 0, 1, false)

	paddedFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		// AddItem(nil, 0, 1, false).
		AddItem(paddedFlexInner, mainFlexHeight, 0, false)
		// AddItem(nil, 0, 1, false)
		paddedFlex.SetInputCapture(state.InputCapture)
	
	go state.UpdateLoop()
	if err := app.SetRoot(paddedFlex, true).SetFocus(paddedFlex).Run(); err != nil {
		panic(err)
	}
}