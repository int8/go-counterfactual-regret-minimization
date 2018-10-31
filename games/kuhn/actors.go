package kuhn

import (
	"errors"
	"github.com/int8/go-counterfactual-regret-minimization/acting"
	"github.com/int8/go-counterfactual-regret-minimization/cards"
	"github.com/int8/go-counterfactual-regret-minimization/table"
)

type Chance struct {
	id   acting.ActorID
	deck cards.Deck
}

func (chance *Chance) GetID() acting.ActorID {
	return chance.id
}

type Player struct {
	Id      acting.ActorID
	Card    *cards.Card
	Stack   float32
	Actions []acting.Action
}

func (player *Player) GetID() acting.ActorID {
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

func (player *Player) Opponent() acting.ActorID {
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
	panic(errors.New("code not reachable"))
}
