package kuhn

import . "github.com/int8/gopoker"
import "github.com/int8/gopoker/cards"

type PlayerAction struct {
	name ActionName
}

func (a PlayerAction) Name() ActionName {
	return a.name
}

type DealPrivateCardsAction struct {
	CardA *cards.Card
	CardB *cards.Card
}

func (a DealPrivateCardsAction) Name() ActionName {
	return DealPrivateCards
}
