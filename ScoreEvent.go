package main

import "strconv"

// Scores """"enum""""
type Score int64
const (
	scoreSoftDrop Score = iota
	scoreHardDrop Score = iota
	
	scoreLineSingle Score = iota
	scoreLineDouble Score = iota
	scoreLineTriple Score = iota
	scoreLineQuad Score = iota

	scoreTSpinSingle Score = iota
	scoreTSpinDouble Score = iota
	scoreTSpinTriple Score = iota

	scorePerfectSingle Score = iota
	scorePerfectDouble Score = iota
	scorePerfectTriple Score = iota
	scorePerfectQuad Score = iota

	scoreCombo Score = iota
)

type ScoreEvent struct {
	scoreType Score
	
	disappearDelay int
	disappearDelayMax int
	count int
	value int
}

func NewScoreEvent(score_type Score) *ScoreEvent {
	this := ScoreEvent{}
	this.scoreType = score_type
	return &this
}

func (this *ScoreEvent) GetText() string {
	var name string
	switch (this.scoreType) {
	case scoreSoftDrop:
		name = "Soft Drop"
	case scoreHardDrop:
		name = "Hard Drop"

	case scoreLineSingle:
		name = "Single Line"
	case scoreLineDouble:
		name = "Double Line"
	case scoreLineTriple:
		name = "Triple Line!"
	case scoreLineQuad:
		name = "Tetris!!"

	case scoreTSpinSingle:
		name = "T-Spin Single"
	case scoreTSpinDouble:
		name = "T-Spin Double!"
	case scoreTSpinTriple:
		name = "T-Spin Triple!!"

	case scorePerfectSingle:
		name = "PC Single!"
	case scorePerfectDouble:
		name = "PC Double!"
	case scorePerfectTriple:
		name = "PC Triple!!"
	case scorePerfectQuad:
		name = "PC Tetris!!!"

	case scoreCombo:
		name = "Line combo"

	default:
		panic("Unknown ScoreEvent type")
	}

	count := " x" + strconv.Itoa(this.count)
	return name + count
}

func (this *ScoreEvent) GetScore() int {
	return this.count * this.value
}