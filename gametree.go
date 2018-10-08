package gocfr

import (
	"errors"
)

const RIMaxActionsPerGame int = 3 * (MaxRaises + 3) // CHECK + BET + MaxRaises + CALL + Chance

// RIGameState
type RIGameState struct {
	round         Round
	parent        *RIGameState
	causingAction Action
	table         *Table
	actors        map[ActorId]Actor
	nextToMove    ActorId
	terminal      bool
}

func (state *RIGameState) Child(action Action) GameStateHolder {
	return state.actors[state.nextToMove].Act(state, action)
}

func (state *RIGameState) Actions() []Action {
	return state.actors[state.nextToMove].GetAvailableActions(state)
}

func (state *RIGameState) IsChance() bool {
	return state.nextToMove == ChanceId
}

func (state *RIGameState) IsTerminal() bool {
	return state.terminal
}

func (state *RIGameState) CurrentActor() Actor {
	return state.actors[state.nextToMove]
}

func (state *RIGameState) Evaluate() float64 {
	if state.IsTerminal() {
		if state.causingAction == Fold {
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
	panic(errors.New("RIGameState is not terminal"))
}

func (state *RIGameState) CurrentInformationSet() InformationSet {

	privateCardName := byte(state.actors[state.nextToMove].(*Player).card.name)
	privateCardSuit := byte(state.actors[state.nextToMove].(*Player).card.suit)
	flopCardName := byte(NoCardName)
	flopCardSuit := byte(NoCardSuit)
	turnCardName := byte(NoCardName)
	turnCardSuit := byte(NoCardSuit)

	if len(state.table.cards) > 0 {
		flopCardName = byte(state.table.cards[0].name)
		flopCardSuit = byte(state.table.cards[0].suit)
	}

	if len(state.table.cards) > 1 {
		turnCardName = byte(state.table.cards[1].name)
		turnCardSuit = byte(state.table.cards[1].suit)
	}
	informationSet := [InformationSetSize]byte{privateCardSuit, privateCardName, flopCardName, flopCardSuit, turnCardName, turnCardSuit}
	// there is no more than 50 actions overall
	currentState := state
	for i := 6; i < 6+RIMaxActionsPerGame; i++ {
		informationSet[i] = byte(currentState.causingAction)
		currentState = currentState.parent
		if currentState == nil {
			break
		}
	}
	return InformationSet(informationSet)
}

func (state *RIGameState) betSize() float64 {
	if state.round < Flop {
		return PreFlopBetSize
	}
	return PostFlopBetSize
}

func CreateRoot(playerA *Player, playerB *Player) *RIGameState {
	chance := &Chance{id: ChanceId, deck: CreateFullDeck(true)}

	actors := map[ActorId]Actor{PlayerA: playerA, PlayerB: playerB, ChanceId: chance}
	table := &Table{pot: 0, cards: []Card{}}

	return &RIGameState{round: Start, table: table,
		actors: actors, nextToMove: ChanceId, causingAction: NoAction}
}

func (state *RIGameState) CreateChild(round Round, Action Action, nextToMove ActorId, terminal bool) *RIGameState {
	child := RIGameState{round: round,
		parent: state, causingAction: Action, terminal: terminal,
		table: state.table.Clone(), actors: cloneActorsMap(state.actors), nextToMove: nextToMove}
	return &child
}
