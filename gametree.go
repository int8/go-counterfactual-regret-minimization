package gocfr

import "errors"

// GameState
type GameState struct {
	round       Round
	parent      *GameState
	causingMove Move
	table       *Table
	actors      map[ActorId]Actor
	nextToMove  ActorId
	terminal    bool
}

func (state *GameState) CurrentActor() Actor {
	return state.actors[state.nextToMove]
}

func (state *GameState) BetSize() float64 {
	if state.round < Flop {
		return PreFlopBetSize
	}
	return PostFlopBetSize
}

func CreateRoot(playerA *Player, playerB *Player) *GameState {
	chance := &Chance{id: ChanceId, deck: CreateFullDeck(true)}

	actors := map[ActorId]Actor{PlayerA: playerA, PlayerB: playerB, ChanceId: chance}
	table := &Table{pot: 0, cards: []Card{}}

	return &GameState{round: Start, table: table,
		actors: actors, nextToMove: ChanceId, causingMove: NoMove}
}

func (state *GameState) CreateChild(round Round, move Move, nextToMove ActorId, terminal bool) *GameState {
	child := GameState{round: round,
		parent: state, causingMove: move, terminal: terminal,
		table: state.table.Clone(), actors: cloneActorsMap(state.actors), nextToMove: nextToMove}
	return &child
}

func (state *GameState) IsTerminal() bool {
	return state.terminal
}

func (state *GameState) Evaluate() float64 {
	if state.IsTerminal() {
		if state.causingMove == Fold {
			return float64(-state.parent.nextToMove) * state.table.pot
		}
		playerAHandValueVector := state.actors[PlayerA].(*Player).EvaluateHand(state.table)
		playerBHandValueVector := state.actors[PlayerB].(*Player).EvaluateHand(state.table)
		for i := range playerAHandValueVector {
			if playerAHandValueVector[i] == playerBHandValueVector[i] {
				continue
			}
			if playerAHandValueVector[i] > playerBHandValueVector[i] {
				return state.table.pot
			} else {
				return -state.table.pot
			}
		}
		return 0.0
	}
	//TODO: make it idiomatic one day..
	panic(errors.New("GameState is not terminal"))
}
