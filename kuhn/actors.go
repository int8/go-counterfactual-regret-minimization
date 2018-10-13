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
	Id      ActorId
	Card    *Card
	Stack   float64
	Actions []Action
}

func (player *Player) GetId() ActorId {
	return player.Id
}

func (player *Player) UpdateStack(stack float64) {
	player.Stack = stack
}

func (chance *Chance) Clone() *Chance {
	return &Chance{id: chance.id, deck: chance.deck.Clone()}
}

func (player *Player) Clone() *Player {
	return &Player{Card: player.Card, Id: player.Id, Stack: player.Stack, Actions: nil}
}

func (player *Player) Opponent() ActorId {
	return -player.Id
}

func (player *Player) CollectPrivateCard(card *Card) {
	player.Card = card
}

func (player *Player) PlaceBet(table *Table, betSize float64) {
	table.AddToPot(betSize)
	player.Stack -= betSize
}

func (player *Player) EvaluateHand(table *Table) int8 {
	return int8((*player).Card.Name)
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
