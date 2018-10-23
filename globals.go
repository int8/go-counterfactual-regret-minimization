package gopoker

const (
	PlayerA  ActorID = 1
	PlayerB          = -PlayerA
	ChanceId         = 0
)

const (
	Start Round = iota
	PreFlop
	Flop
	Turn
	River
	End
)

var (
	NoAction         ActionName = to3BinArray(0)
	DealPublicCards  ActionName = to3BinArray(1)
	DealPrivateCards ActionName = to3BinArray(2)
	Fold             ActionName = to3BinArray(3)
	Check            ActionName = to3BinArray(4)
	Bet              ActionName = to3BinArray(5)
	Call             ActionName = to3BinArray(6)
	Raise            ActionName = to3BinArray(7)
)
