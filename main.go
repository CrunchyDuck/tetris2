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
	state := NewGameState()

	var app *tview.Application
	app = tview.NewApplication()

	tetrisField := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWrap(false).
		SetChangedFunc(func() { app.Draw() })
	tetrisField.SetBorder(true)

	rootFlex := tview.NewFlex()
	vertFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(nil, 0, 1, false).
		AddItem(tetrisField, state.height + 2, 1, false).
		AddItem(nil, 0, 1, false)

	rootFlex.
		AddItem(nil, 0, 1, false).
		AddItem(vertFlex, state.width + 2, 1, false).
		AddItem(nil, 0, 1, false)
	rootFlex.SetInputCapture(state.InputCapture)
	
	go state.UpdateLoop(tetrisField)
	if err := app.SetRoot(rootFlex, true).SetFocus(rootFlex).Run(); err != nil {
		panic(err)
	}
}