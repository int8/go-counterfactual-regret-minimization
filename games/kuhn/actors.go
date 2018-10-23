package kuhn

import (
	"errors"
	. "github.com/int8/gopoker"
	"github.com/int8/gopoker/cards"
	"github.com/int8/gopoker/table"
)

type Chance struct {
	id   ActorID
	deck cards.Deck
}

func (chance *Chance) GetID() ActorID {
	return chance.id
}

type Player struct {
	Id      ActorID
	Card    *cards.Card
	Stack   float32
	Actions []Action
}

func (player *Player) GetID() ActorID {
	return player.Id
}

func (player *Player) UpdateStack(stack float32) {
	player.Stack = stack
}

func (chance *Chance) Clone() *Chance {
	return &Chance{id: chance.id, deck: chance.deck.Clone()}
}

func (player *Player) Clone() *Player {
	return &Player{Card: player.Card, Id: player.Id, Stack: player.Stack, Actions: nil}
}

func (player *Player) Opponent() ActorID {
	return -player.Id
}

func (player *Player) CollectPrivateCard(card *cards.Card) {
	player.Card = card
}

func (player *Player) PlaceBet(table *table.PokerTable, betSize float32) {
	table.AddToPot(betSize)
	player.Stack -= betSize
}

func (player *Player) EvaluateHand(table *table.PokerTable) int8 {
	return cards.CardSymbol2Int((*player).Card.Symbol)
}

func (player *Player) String() string {
	if player.Id == 1 {
		return "A"
	} else if player.Id == -1 {
		return "B"
	} else {
		return "Chance"
	}
	//TODO: not idiomatic !
	panic(errors.New("Code not reachable."))
}
