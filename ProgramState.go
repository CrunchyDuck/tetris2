package main 

import (
	"github.com/rivo/tview"
)

type PState int64
const (
	programMenu PState = iota
	programGame PState = iota
	programGameOver PState = iota
)

type ProgramState struct {
	gracefulExit bool
	targetState PState
	ascii bool
	app *tview.Application
}

func NewProgramState(app *tview.Application) *ProgramState {
	this := ProgramState{}
	this.app = app
	this.gracefulExit = true
	this.ascii = true
	this.targetState = programMenu
	return &this
}