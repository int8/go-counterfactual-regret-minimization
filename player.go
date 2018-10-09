package gocfr

import (
	"errors"
)

type ActorId int8

const (
	PlayerA  ActorId = 1
	PlayerB          = -PlayerA
	ChanceId         = 0
	NoActor          = 100
)

type Actor interface {
	Act(state *RIGameState, action Action) *RIGameState
	GetAvailableActions(state *RIGameState) []Action
}

type Chance struct {
	id   ActorId
	deck *FullDeck
}

func (chance *Chance) Act(state *RIGameState, action Action) (child *RIGameState) {

	if action == DealPublicCard {
		child = chance.dealPublicCard(state)
	}

	if action == DealPrivateCards {
		child = chance.dealPrivateCards(state)
	}
	return child
}

func (chance *Chance) dealPublicCard(state *RIGameState) *RIGameState {

	child := state.CreateChild(state.round.NextRound(), DealPublicCard, PlayerA, false)
	// important to deal using child deck / not current chance deck
	child.table.DropPublicCard(child.actors[ChanceId].(*Chance).deck.DealNextCard())
	return child
}

func (chance *Chance) dealPrivateCards(state *RIGameState) *RIGameState {

	child := state.CreateChild(state.round.NextRound(), DealPrivateCards, PlayerA, false)
	// important to deal using child deck / not current chance deck
	child.actors[PlayerA].(*Player).PlaceBet(child.table, Ante)
	child.actors[PlayerB].(*Player).PlaceBet(child.table, Ante)
	child.actors[PlayerA].(*Player).CollectPrivateCard(child.actors[ChanceId].(*Chance).deck.DealNextCard())
	child.actors[PlayerB].(*Player).CollectPrivateCard(child.actors[ChanceId].(*Chance).deck.DealNextCard())
	return child
}

func (chance *Chance) GetAvailableActions(state *RIGameState) []Action {
	if state.round == Start {
		return []Action{DealPrivateCards}
	}
	if !state.terminal {
		return []Action{DealPublicCard}
	}
	return []Action{}
}

type Player struct {
	id      ActorId
	card    *Card
	stack   float64
	actions []Action
}

func (chance *Chance) Clone() *Chance {
	return &Chance{id: chance.id, deck: chance.deck.Clone()}
}

//TODO: it is getting messy, think of structuring it better
func (player *Player) Act(state *RIGameState, action Action) (child *RIGameState) {

	if !actionInSlice(action, player.GetAvailableActions(state)) {
		panic("action not available")
	}

	betSize := state.betSize()

	defer func() {
		if action == Call || action == Bet {
			child.actors[player.id].(*Player).PlaceBet(child.table, betSize)
		}
		if action == Raise {
			child.actors[player.id].(*Player).PlaceBet(child.table, 2*betSize)
		}
		if action == Fold {
			if child.round < Flop {
				child.actors[-state.nextToMove].(*Player).PlaceBet(child.table, -PreFlopBetSize)
			} else {
				child.actors[-state.nextToMove].(*Player).PlaceBet(child.table, -PostFlopBetSize)
			}
		}
	}()

	if action == Fold || (state.round == Turn && (action == Call || (action == Check && state.causingAction == Check))) {
		child = state.CreateChild(state.round, action, player.Opponent(), true)
		return
	}

	if action == Call || (action == Check && state.causingAction == Check) {
		child = state.CreateChild(state.round, action, ChanceId, false)
		return
	}

	child = state.CreateChild(state.round, action, player.Opponent(), false)
	return

}

func (player *Player) GetAvailableActions(state *RIGameState) []Action {
	player.computeAvailableActions(state)
	return player.actions
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

func (player *Player) computeAvailableActions(state *RIGameState) {

	if player.actions != nil {
		return
	}

	if state.causingAction == Fold {
		player.actions = []Action{}
		return
	}
	betSize := state.betSize()

	opponentStack := state.actors[player.Opponent()].(*Player).stack

	allowedToBet := (player.stack >= betSize) && (opponentStack >= betSize)
	allowedToRaise := (player.stack >= 2*betSize) && (opponentStack >= 2*betSize)

	// whenever betting roung is over (CALL OR CHECK->CHECK)
	bettingRoundEnded := state.causingAction == Call || (state.causingAction == Check && state.parent.causingAction == Check)
	if bettingRoundEnded {
		player.actions = []Action{}
		return
	}

	// single check implies BET or CHECK
	if state.causingAction == Check && state.parent.causingAction != Check {
		player.actions = []Action{Check}
		if allowedToBet {
			player.actions = append(player.actions, Bet)
		}
		return
	}

	// if RAISE/BET, you can CALL FOLD or RAISE (unless there has been 6 prior raises - 3 for each player)
	if state.causingAction == Bet || state.causingAction == Raise {
		player.actions = []Action{Call, Fold}
		priorRaisesInCurrentRound := countPriorRaisesPerRound(state, state.round)
		if priorRaisesInCurrentRound < MaxRaises && allowedToRaise {
			player.actions = append(player.actions, Raise)
		}
		return
	}

	if state.causingAction == DealPrivateCards || state.causingAction == DealPublicCard {
		player.actions = []Action{Check}
		if allowedToBet {
			player.actions = append(player.actions, Bet)
		}
		return
	}
	//TODO: not idiomatic !
	panic(errors.New("Code not reachable."))
}

func (player *Player) EvaluateHand(table *Table) []int8 {

	var flush, three, pair, straight, ownCard int8

	if (*player).card.suit == table.cards[0].suit && (*player).card.suit == table.cards[1].suit {
		flush = 1
	}

	if ((*player).card.name == table.cards[0].name) && ((*player).card.name == table.cards[1].name) {
		three = 1
	}

	if (((*player).card.name == table.cards[0].name) || ((*player).card.name == table.cards[1].name)) || table.cards[0].name == table.cards[1].name {
		pair = 1
	}

	if pair == 0 && cardsDiffersByTwo([]Card{*player.card, table.cards[0], table.cards[1]}) {
		straight = 1
	}

	ownCard = int8((*player).card.name)

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
