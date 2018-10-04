package gocfr

type Strategy func(state RhodeIslandGameState) map[Action]float64

type HeadsUpGame struct {
	root *TwoPlayersGameNode
	players []Player
	playersStacks map[Player]float64
	table PokerTable
	strategyProfile map[Player]Strategy
}



