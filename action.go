package gocfr

import "fmt"

type Move int8

const (
	NoMove Move = iota
	Check
	Bet
	Raise
	Call
	Fold
	DealPublicCard
	DealPrivateCards
)

const PreFlopBetSize float64 = 10.
const PostFlopBetSize float64 = 10.
const Ante float64 = 5.0

type Action struct {
	player ActionMaker
	move   Move
}

type ActionsCache struct {
	actions []Action
}

func (a Action) String() string {
	return fmt.Sprintf("%v:%v", a.player, a.move)
}

func (m Move) String() string {
	switch m {
	case Check:
		return "Check"
	case Bet:
		return "Bet"
	case Call:
		return "Call"
	case Fold:
		return "Fold"
	case Raise:
		return "Raise"
	case DealPrivateCards:
		return "DealPrivateCards"
	case DealPublicCard:
		return "DealPublicCard"
	}
	return "Undefined"
}
