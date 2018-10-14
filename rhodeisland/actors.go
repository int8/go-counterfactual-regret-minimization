package rhodeisland

import (
	"errors"
	. "github.com/int8/gopoker"
)

type Chance struct {
	id   ActorId
	deck Deck
}

func (chance *Chance) GetId() ActorId {
	return chance.id
}

type Player struct {
	id      ActorId
	card    *Card
	stack   float32
	actions []Action
}

func (player *Player) GetId() ActorId {
	return player.id
}

func (player *Player) UpdateStack(stack float32) {
	player.stack = stack
}

func (chance *Chance) Clone() *Chance {
	return &Chance{id: chance.id, deck: chance.deck.Clone()}
}

func (player *Player) Clone() *Player {
	return &Player{card: player.card, id: player.id, stack: player.stack, actions: nil}
}

func (player *Player) Opponent() ActorId {
	return -player.id
}

func (player *Player) CollectPrivateCard(card *Card) {
	player.card = card
}

func (player *Player) PlaceBet(table *Table, betSize float32) {
	table.AddToPot(betSize)
	player.stack -= betSize
}

func (player *Player) EvaluateHand(table *Table) []int8 {

	var flush, three, pair, straight, ownCard int8

	if (*player).card.Suit == table.Cards[0].Suit && (*player).card.Suit == table.Cards[1].Suit {
		flush = 1
	}

	if ((*player).card.Name == table.Cards[0].Name) && ((*player).card.Name == table.Cards[1].Name) {
		three = 1
	}

	if (((*player).card.Name == table.Cards[0].Name) || ((*player).card.Name == table.Cards[1].Name)) || table.Cards[0].Name == table.Cards[1].Name {
		pair = 1
	}

	if pair == 0 && cardsDiffersByTwo([]Card{*player.card, table.Cards[0], table.Cards[1]}) {
		straight = 1
	}

	ownCard = int8((*player).card.Name)

	return []int8{straight * flush, three, straight, flush, pair, ownCard}
}

func (player *Player) String() string {
	if player.id == 1 {
		return "A"
	} else if player.id == -1 {
		return "B"
	} else {
		return "Chance"
	}
	//TODO: not idiomatic !
	panic(errors.New("Code not reachable."))
}
