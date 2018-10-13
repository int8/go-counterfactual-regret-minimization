package kuhn

import . "github.com/int8/gopoker"

type PlayerAction struct {
	name ActionName
}

func (a PlayerAction) Name() ActionName {
	return a.name
}

type DealPrivateCardsAction struct {
	CardA *Card
	CardB *Card
}

func (a DealPrivateCardsAction) Name() ActionName {
	return DealPrivateCards
}
