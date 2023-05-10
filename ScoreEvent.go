package main

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
	baseScore int
	totalScore int
	name string
}

func NewScoreEvent(score_type Score) *ScoreEvent {
	this := ScoreEvent{}
	this.scoreType = score_type
	this.disappearDelayMax = 300
	this.disappearDelay = this.disappearDelayMax
	this.name, this.baseScore = GetScoreData(score_type)
	return &this
}

func GetScoreData(score_type Score) (name string, value int) {
	switch (score_type) {
	case scoreSoftDrop:
		name = "Soft Drop"
		value = 1
	case scoreHardDrop:
		name = "Hard Drop"
		value = 2

	case scoreLineSingle:
		name = "Single Line"
		value = 100
	case scoreLineDouble:
		name = "Double Line"
		value = 300
	case scoreLineTriple:
		name = "Triple Line!"
		value = 500
	case scoreLineQuad:
		name = "Tetris!!"
		value = 800

	case scoreTSpinSingle:
		name = "T-Spin Single"
		value = 800
	case scoreTSpinDouble:
		name = "T-Spin Double!"
		value = 1200
	case scoreTSpinTriple:
		name = "T-Spin Triple!!"
		value = 1600

	case scorePerfectSingle:
		name = "PC Single!"
		value = 800
	case scorePerfectDouble:
		name = "PC Double!"
		value = 1200
	case scorePerfectTriple:
		name = "PC Triple!!"
		value = 1800
	case scorePerfectQuad:
		name = "PC Tetris!!!"
		value = 2000

	case scoreCombo:
		name = "Line combo"
		value = 50

	default:
		panic("Unknown ScoreEvent type")
	}

	return name, value
}

func (this *ScoreEvent) IncreaseScore(level int) int {
	// Based on: https://tetris.wiki/Scoring
	this.disappearDelay = this.disappearDelayMax

	var scoreIncrease int
	if this.scoreType == scoreCombo {
		scoreIncrease = this.baseScore * this.count * level
	} else if this.scoreType == scoreSoftDrop || this.scoreType == scoreHardDrop { // These don't get level bonuses.
		scoreIncrease = this.baseScore
	} else {
		scoreIncrease = this.baseScore * level
	}
	this.count++
	this.totalScore += scoreIncrease

	return scoreIncrease
}