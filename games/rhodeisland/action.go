package rhodeisland

import (
	"github.com/int8/gopoker/acting"
	"github.com/int8/gopoker/cards"
)

type PlayerAction struct {
	name acting.ActionName
}

var CheckAction = PlayerAction{acting.Check}
var BetAction = PlayerAction{acting.Bet}
var CallAction = PlayerAction{acting.Call}
var RaiseAction = PlayerAction{acting.Raise}
var FoldAction = PlayerAction{acting.Fold}

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

type DealPublicCardAction struct {
	Card *cards.Card
}

func (a DealPublicCardAction) Name() acting.ActionName {
	return acting.DealPublicCards

}
