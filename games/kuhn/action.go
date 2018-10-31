package kuhn

import (
	"github.com/int8/go-counterfactual-regret-minimization/acting"
)
import "github.com/int8/go-counterfactual-regret-minimization/cards"

type PlayerAction struct {
	name acting.ActionName
}

func (a PlayerAction) Name() acting.ActionName {
	return a.name
}

type DealPrivateCardsAction struct {
	CardA *cards.Card
	CardB *cards.Card
}

func (a DealPrivateCardsAction) Name() acting.ActionName {
	return acting.DealPrivateCards
}

var (
	CheckAction = PlayerAction{acting.Check}
	BetAction   = PlayerAction{acting.Bet}
	CallAction  = PlayerAction{acting.Call}
	FoldAction  = PlayerAction{acting.Fold}
)
