package kuhn

import (
	"errors"
	. "github.com/int8/gopoker"
)

// KuhnGameState - Kuhn Poker Game State
type KuhnGameState struct {
	round         Round
	parent        *KuhnGameState
	causingAction Action
	table         *Table
	actors        map[ActorId]Actor
	nextToMove    ActorId
	terminal      bool
}

func (state *KuhnGameState) Act(action Action) GameState {
	switch state.actors[state.nextToMove].(type) {
	case *Chance:
		return state.actAsChance(action)
	case *Player:
		return state.actAsPlayer(action)
	}
	return nil
}

func (state *KuhnGameState) Actions() []Action {

	switch state.actors[state.nextToMove].(type) {
	case *Chance:
		return state.chanceActions(state.actors[state.nextToMove].(*Chance))
	case *Player:
		return state.playerActions(state.actors[state.nextToMove].(*Player))
	}
	return nil
}

func (state *KuhnGameState) IsChance() bool {
	return state.nextToMove == ChanceId
}

func (state *KuhnGameState) IsTerminal() bool {
	return state.terminal
}

func (state *KuhnGameState) Parent() GameState {
	return state.parent
}

func (state *KuhnGameState) CurrentActor() Actor {
	return state.actors[state.nextToMove]
}

//TODO: test it carefully
func (state *KuhnGameState) Evaluate() float64 {
	currentActor := state.playerActor(state.CurrentActor().GetId())
	currentActorOpponent := state.playerActor(-state.CurrentActor().GetId())
	if state.IsTerminal() {
		if state.causingAction.Name() == Fold {
			currentActor.UpdateStack(state.table.Pot)
			return float64(-state.parent.nextToMove) * state.table.Pot
		}
		currentActorHandValue := currentActor.EvaluateHand(state.table)
		currentActorOpponentHandValue := currentActorOpponent.EvaluateHand(state.table)

		if currentActorHandValue > currentActorOpponentHandValue {
			currentActor.UpdateStack(state.table.Pot)
			return float64(currentActor.GetId()) * state.table.Pot
		} else {
			currentActorOpponent.UpdateStack(state.table.Pot)
			return float64(currentActor.GetId()) * state.table.Pot
		}
	}
	currentActor.UpdateStack(state.table.Pot / 2)
	currentActorOpponent.UpdateStack(state.table.Pot / 2)
	return 0.0
	panic(errors.New("RIGameState is not terminal"))
}

func (state *KuhnGameState) InformationSet() InformationSet {

	privateCardName := byte(state.actors[state.nextToMove].(*Player).card.Name)
	privateCardSuit := byte(state.actors[state.nextToMove].(*Player).card.Suit)

	informationSet := [InformationSetSize]byte{privateCardName, privateCardSuit}
	// there is no more than 50 actions overall
	currentState := state
	for i := 2; currentState.round != Start; i++ {
		informationSet[i] = byte(currentState.causingAction.Name())
		currentState = currentState.parent
		if currentState == nil {
			break
		}
	}
	return InformationSet(informationSet)
}

func (state *KuhnGameState) stack(actor ActorId) float64 {
	return state.actors[actor].(*Player).stack
}

func (state *KuhnGameState) actAsChance(action Action) GameState {
	var child *KuhnGameState
	if action.Name() == DealPrivateCards {
		child = state.dealPrivateCards(action.(DealPrivateCardsAction).CardA, action.(DealPrivateCardsAction).CardB)
	}
	return child
}

func (state *KuhnGameState) actAsPlayer(action Action) GameState {

	var child *KuhnGameState

	if !actionInSlice(action, state.Actions()) {
		panic("action not available")
	}
	actor := state.CurrentActor()

	defer func() {
		if action.Name() == Call || action.Name() == Bet {
			child.actors[actor.GetId()].(*Player).PlaceBet(child.table, BetSize)
		}
		if action.Name() == Fold {
			child.actors[-state.nextToMove].(*Player).PlaceBet(child.table, -BetSize)
		}
	}()

	if action.Name() == Fold || action.Name() == Call || (action.Name() == Check && state.causingAction.Name() == Check) {
		child = createChild(state, state.round, action, state.CurrentActor().(*Player).Opponent(), true)
		return child
	}

	child = createChild(state, state.round, action, state.CurrentActor().(*Player).Opponent(), false)
	return child
}

func (state *KuhnGameState) betSize() float64 {
	return 1.0
}

func root(playerA *Player, playerB *Player) *KuhnGameState {
	chance := &Chance{id: ChanceId, deck: CreateKuhnDeck(true)}

	actors := map[ActorId]Actor{PlayerA: playerA, PlayerB: playerB, ChanceId: chance}
	table := &Table{Pot: 0, Cards: []Card{}}

	return &KuhnGameState{round: Start, table: table,
		actors: actors, nextToMove: ChanceId, causingAction: nil}
}

func createChild(blueprint *KuhnGameState, round Round, Action Action, nextToMove ActorId, terminal bool) *KuhnGameState {
	child := KuhnGameState{round: round,
		parent: blueprint, causingAction: Action, terminal: terminal,
		table: blueprint.table.Clone(), actors: cloneActorsMap(blueprint.actors), nextToMove: nextToMove}
	return &child
}

func (state *KuhnGameState) dealPrivateCards(cardA *Card, cardB *Card) *KuhnGameState {

	child := createChild(state, state.round.NextRound(), DealPrivateCardsAction{cardA, cardB}, PlayerA, false)
	// important to deal using child deck / not current chance deck
	child.actors[PlayerA].(*Player).PlaceBet(child.table, Ante)
	child.actors[PlayerB].(*Player).PlaceBet(child.table, Ante)
	child.actors[PlayerA].(*Player).CollectPrivateCard(cardA)
	child.actors[PlayerB].(*Player).CollectPrivateCard(cardB)
	child.actors[ChanceId].(*Chance).deck.RemoveCard(cardA)
	child.actors[ChanceId].(*Chance).deck.RemoveCard(cardB)

	return child
}

func (state *KuhnGameState) chanceActions(chance *Chance) []Action {
	if state.round == Start {
		deckSize := int(chance.deck.CardsLeft())
		actions := make([]Action, deckSize*(deckSize-1))
		i := 0
		remainingCards := chance.deck.RemainingCards()
		for _, cardA := range remainingCards {
			for _, cardB := range remainingCards {
				{
					if cardA != cardB {
						actions[i] = DealPrivateCardsAction{cardA, cardB}
						i++
					}
				}
			}
		}
		return actions
	}
	return nil
}

func (state *KuhnGameState) playerActions(player *Player) []Action {

	if state.causingAction.Name() == Fold {
		player.actions = []Action{}
		return player.actions
	}

	betSize := state.betSize()
	opponentStack := state.stack(player.Opponent())
	allowedToBet := (player.stack >= betSize) && (opponentStack >= betSize)

	// whenever betting round is over (CALL OR CHECK->CHECK)
	bettingRoundEnded := state.causingAction.Name() == Call || (state.causingAction.Name() == Check && state.parent.causingAction.Name() == Check)
	if bettingRoundEnded {
		player.actions = []Action{}
		return player.actions
	}

	// single check implies BET or CHECK
	if state.causingAction.Name() == Check && state.parent.causingAction.Name() != Check {
		player.actions = []Action{CheckAction}
		if allowedToBet {
			player.actions = append(player.actions, BetAction)
		}
		return player.actions
	}

	if state.causingAction.Name() == Bet {
		player.actions = []Action{CallAction, FoldAction}
		return player.actions
	}

	if state.causingAction.Name() == DealPrivateCards || state.causingAction.Name() == DealPublicCards {
		player.actions = []Action{CheckAction}
		if allowedToBet {
			player.actions = append(player.actions, BetAction)
		}
		return player.actions
	}
	panic(errors.New("Code not reachable."))
}

func (state *KuhnGameState) playerActor(id ActorId) *Player {
	return state.actors[id].(*Player)
}
