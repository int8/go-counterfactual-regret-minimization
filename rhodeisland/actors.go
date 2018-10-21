package rhodeisland

import (
	"errors"
	. "github.com/int8/gopoker"
)

type Chance struct {
	id   ActorID
	deck Deck
}

func (chance *Chance) GetID() ActorID {
	return chance.id
}

type Player struct {
	Id      ActorID
	Card    *Card
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

func (player *Player) CollectPrivateCard(card *Card) {
	player.Card = card
}

func (player *Player) PlaceBet(table *Table, betSize float32) {
	table.AddToPot(betSize)
	player.Stack -= betSize
}

func (player *Player) EvaluateHand(table *Table) []int8 {

	var flush, three, pair, straight, ownCard int8

	if (*player).Card.Suit == table.Cards[0].Suit && (*player).Card.Suit == table.Cards[1].Suit {
		flush = 1
	}

	if ((*player).Card.Symbol == table.Cards[0].Symbol) && ((*player).Card.Symbol == table.Cards[1].Symbol) {
		three = 1
	}

	if (((*player).Card.Symbol == table.Cards[0].Symbol) || ((*player).Card.Symbol == table.Cards[1].Symbol)) || table.Cards[0].Symbol == table.Cards[1].Symbol {
		pair = 1
	}

	if pair == 0 && cardsDiffersByTwo([]Card{*player.Card, table.Cards[0], table.Cards[1]}) {
		straight = 1
	}

	ownCard = CardSymbol2Int((*player).Card.Symbol)

	return []int8{straight * flush, three, straight, flush, pair, ownCard}
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
