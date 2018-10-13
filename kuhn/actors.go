package kuhn

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
	stack   float64
	actions []Action
}

func (player *Player) GetId() ActorId {
	return player.id
}

func (player *Player) UpdateStack(stack float64) {
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

func (player *Player) PlaceBet(table *Table, betSize float64) {
	table.AddToPot(betSize)
	player.stack -= betSize
}

func (player *Player) EvaluateHand(table *Table) int8 {
	return int8((*player).card.Name)
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
