package gocfr

import "fmt"

type Move int8

const (
	NoMove Move = iota
	DealPublicCard
	DealPrivateCards
	Fold
	Check
	Bet
	Call
	Raise
)

const PreFlopBetSize float64 = 10.
const PostFlopBetSize float64 = 20.
const Ante float64 = 5.0

type Action struct {
	player Actor
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
