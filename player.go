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
	Act(state *GameState, move Move) *GameState
	GetAvailableMoves(state *GameState) []Move
}

type Chance struct {
	id   ActorId
	deck *FullDeck
}

func (chance *Chance) Act(state *GameState, move Move) (child *GameState) {

	if move == DealPublicCard {
		child = chance.dealPublicCard(state)
	}

	if move == DealPrivateCards {
		child = chance.dealPrivateCards(state)
	}
	return child
}

func (chance *Chance) dealPublicCard(state *GameState) *GameState {

	child := state.CreateChild(state.round.NextRound(), DealPublicCard, PlayerA, false)
	// important to deal using child deck / not current chance deck
	child.table.DropPublicCard(child.actors[ChanceId].(*Chance).deck.DealNextCard())
	return child
}

func (chance *Chance) dealPrivateCards(state *GameState) *GameState {

	child := state.CreateChild(state.round.NextRound(), DealPrivateCards, PlayerA, false)
	// important to deal using child deck / not current chance deck
	child.actors[PlayerA].(*Player).PlaceBet(child.table, Ante)
	child.actors[PlayerB].(*Player).PlaceBet(child.table, Ante)
	child.actors[PlayerA].(*Player).CollectPrivateCard(child.actors[ChanceId].(*Chance).deck.DealNextCard())
	child.actors[PlayerB].(*Player).CollectPrivateCard(child.actors[ChanceId].(*Chance).deck.DealNextCard())
	return child
}

func (chance *Chance) GetAvailableMoves(state *GameState) []Move {
	if state.round == Start {
		return []Move{DealPrivateCards}
	}
	if !state.terminal {
		return []Move{DealPublicCard}
	}
	return []Move{}
}

type Player struct {
	id    ActorId
	card  *Card
	stack float64
	moves []Move
}

func (chance *Chance) Clone() *Chance {
	return &Chance{id: chance.id, deck: chance.deck.Clone()}
}

//TODO: it is getting messy, think of structuring it better
func (player *Player) Act(state *GameState, move Move) (child *GameState) {

	betSize := state.betSize()

	defer func() {
		if move == Call || move == Bet {
			child.actors[player.id].(*Player).PlaceBet(child.table, betSize)
		}
		if move == Raise {
			child.actors[player.id].(*Player).PlaceBet(child.table, 2*betSize)
		}
	}()

	if state.round == Turn && (move == Call || (move == Check && state.causingMove == Check)) {
		child = state.CreateChild(state.round, move, NoActor, true)
		return
	}

	if move == Fold {
		child = state.CreateChild(state.round, move, NoActor, true)
		if child.round < Flop {
			child.actors[-state.nextToMove].(*Player).PlaceBet(child.table, -PreFlopBetSize)
		} else {
			child.actors[-state.nextToMove].(*Player).PlaceBet(child.table, -PostFlopBetSize)
		}
		return
	}

	if move == Call || (move == Check && state.causingMove == Check) {
		child = state.CreateChild(state.round, move, ChanceId, false)
		return
	}

	child = state.CreateChild(state.round, move, player.Opponent(), false)
	return

}

func (player *Player) GetAvailableMoves(state *GameState) []Move {
	player.computeAvailableActions(state)
	return player.moves
}

func (player *Player) Clone() *Player {
	return &Player{card: player.card, id: player.id, stack: player.stack, moves: nil}
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

func (player *Player) computeAvailableActions(state *GameState) {

	if player.moves != nil {
		return
	}

	if state.causingMove == Fold {
		player.moves = []Move{}
		return
	}
	betSize := state.betSize()

	opponentStack := state.actors[player.Opponent()].(*Player).stack

	allowedToBet := (player.stack >= betSize) && (opponentStack >= betSize)
	allowedToRaise := (player.stack >= 2*betSize) && (opponentStack >= 2*betSize)

	// whenever betting roung is over (CALL OR CHECK->CHECK)
	bettingRoundEnded := state.causingMove == Call || (state.causingMove == Check && state.parent.causingMove == Check)
	if bettingRoundEnded {
		player.moves = []Move{}
		return
	}

	// single check implies BET or CHECK
	if state.causingMove == Check && state.parent.causingMove != Check {
		player.moves = []Move{Check}
		if allowedToBet {
			player.moves = append(player.moves, Bet)
		}
		return
	}

	// if RAISE/BET, you can CALL FOLD or RAISE (unless there has been 6 prior raises - 3 for each player)
	if state.causingMove == Bet || state.causingMove == Raise {
		player.moves = []Move{Call, Fold}
		if countPriorRaises(state) < MaxRaises && allowedToRaise {
			player.moves = append(player.moves, Raise)
		}
		return
	}

	if state.causingMove == DealPrivateCards || state.causingMove == DealPublicCard {
		player.moves = []Move{Check}
		if allowedToBet {
			player.moves = append(player.moves, Bet)
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
