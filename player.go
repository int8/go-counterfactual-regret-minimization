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

//TODO: now separate chance istance is creater for every game state this is very inefficient
//TODO: Chance is rather heavy, it keeps the whole deck. Consider using the same instance
//TODO: with methods like nextCardExceptOf(excludedCards []Card) etc, keeping dealt cards would then be enough

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

	state.actors[PlayerA].(*Player).PlaceBet(state.table, Ante)
	state.actors[PlayerB].(*Player).PlaceBet(state.table, Ante)
	child := state.CreateChild(state.round.NextRound(), DealPrivateCards, PlayerA, false)
	// important to deal using child deck / not current chance deck
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
	cards []Card
	stack float64
	moves []Move
}

func (chance *Chance) Clone() *Chance {
	return &Chance{id: chance.id, deck: chance.deck.Clone()}
}

func (player *Player) Act(state *GameState, move Move) *GameState {

	table := state.table
	betSize := state.BetSize()

	if move == Call || move == Bet || move == Raise {
		player.PlaceBet(table, betSize)
	}

	if move == Fold || (state.round == Turn && (move == Call || (move == Check && state.causingMove == Check))) {
		return state.CreateChild(state.round, move, NoActor, true)
	}

	if move == Call || (move == Check && state.causingMove == Check) {
		return state.CreateChild(state.round, move, ChanceId, false)
	}

	return state.CreateChild(state.round, move, player.Opponent(), false)
}

func (player *Player) GetAvailableMoves(state *GameState) []Move {
	player.computeAvailableActions(state)
	return player.moves
}

func (player *Player) Clone() *Player {
	cards := make([]Card, len(player.cards))
	copy(cards, player.cards)
	return &Player{cards: cards, id: player.id, stack: player.stack, moves: nil}
}

func (player *Player) Opponent() ActorId {
	return -player.id
}

func (player *Player) CollectPrivateCard(card Card) {
	player.cards = append(player.cards, card)
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
	betSize := state.BetSize()

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
		if countPriorRaises(state) < 6 && allowedToRaise {
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
	panic(errors.New("This code should not be reachable."))
}

func (player Player) String() string {
	if player.id == 1 {
		return "A"
	} else if player.id == -1 {
		return "B"
	} else {
		return "Chance"
	}
	panic(errors.New("This code should not be reachable."))
}
