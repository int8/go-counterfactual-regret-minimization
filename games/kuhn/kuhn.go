package kuhn

import (
	"errors"
	. "github.com/int8/gopoker"
	"github.com/int8/gopoker/cards"
	"github.com/int8/gopoker/table"
)

// KuhnGameState - Kuhn Poker Game State
type KuhnGameState struct {
	round         Round
	parent        *KuhnGameState
	causingAction Action
	table         *table.PokerTable
	actors        map[ActorID]Actor
	nextToMove    ActorID
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
func (state *KuhnGameState) Evaluate() float32 {
	currentActor := state.playerActor(state.CurrentActor().GetID())
	currentActorOpponent := state.playerActor(-state.CurrentActor().GetID())
	if state.IsTerminal() {
		if state.causingAction.Name() == Fold {
			currentActor.UpdateStack(currentActor.Stack + state.table.Pot)
			return float32(currentActor.GetID()) * (state.table.Pot / 2)
		}
		currentActorHandValue := currentActor.EvaluateHand(state.table)
		currentActorOpponentHandValue := currentActorOpponent.EvaluateHand(state.table)

		if currentActorHandValue > currentActorOpponentHandValue {
			currentActor.UpdateStack(currentActor.Stack + state.table.Pot)
			return float32(currentActor.GetID()) * (state.table.Pot / 2)
		} else {
			currentActorOpponent.UpdateStack(currentActorOpponent.Stack + state.table.Pot)
			return float32(currentActorOpponent.GetID()) * (state.table.Pot / 2)
		}
	}
	currentActor.UpdateStack(currentActor.Stack + state.table.Pot/2)
	currentActorOpponent.UpdateStack(currentActorOpponent.Stack + state.table.Pot/2)
	return 0.0
}

func (state *KuhnGameState) InformationSet() InformationSet {

	privateCardSymbol := state.actors[state.nextToMove].(*Player).Card.Symbol
	privateCardSuit := state.actors[state.nextToMove].(*Player).Card.Suit
	informationSet := [InformationSetSizeBytes]byte{}

	informationSetBool := [InformationSetSize]bool{
		privateCardSymbol[0], privateCardSymbol[1], privateCardSymbol[2], privateCardSymbol[3],
		privateCardSuit[0], privateCardSuit[1], privateCardSuit[2],
	}
	currentState := state
	for i := 7; currentState.round != Start; i += 3 {
		actionName := currentState.causingAction.Name()
		informationSetBool[i] = actionName[0]
		informationSetBool[i+1] = actionName[1]
		informationSetBool[i+2] = actionName[2]
		currentState = currentState.parent
		if currentState == nil {
			break
		}
	}
	for i := 0; i < InformationSetSizeBytes; i++ {
		informationSet[i] = CreateByte(informationSetBool[(i * 8):((i + 1) * 8)])
	}

	return InformationSet(informationSet)
}

func (state *KuhnGameState) stack(actor ActorID) float32 {
	return state.actors[actor].(*Player).Stack
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
			child.playerActor(actor.GetID()).PlaceBet(child.table, BetSize)
		}
		if action.Name() == Fold {
			//opponent of folding player can now take his bet back
			child.actors[state.CurrentActor().(*Player).Opponent()].(*Player).PlaceBet(child.table, -BetSize)
		}
	}()

	if action.Name() == Fold || action.Name() == Call || (action.Name() == Check && state.causingAction.Name() == Check) {
		child = createChild(state, state.round, action, state.CurrentActor().(*Player).Opponent(), true)
		return child
	}

	child = createChild(state, state.round, action, state.CurrentActor().(*Player).Opponent(), false)
	return child
}

func (state *KuhnGameState) betSize() float32 {
	return 1.0
}

func Root(playerA *Player, playerB *Player) *KuhnGameState {
	chance := &Chance{id: ChanceId, deck: CreateKuhnDeck(true)}

	actors := map[ActorID]Actor{PlayerA: playerA, PlayerB: playerB, ChanceId: chance}
	table := &table.PokerTable{Pot: 0, Cards: []cards.Card{}}

	return &KuhnGameState{round: Start, table: table,
		actors: actors, nextToMove: ChanceId, causingAction: nil}
}

func createChild(blueprint *KuhnGameState, round Round, Action Action, nextToMove ActorID, terminal bool) *KuhnGameState {
	child := KuhnGameState{round: round,
		parent: blueprint, causingAction: Action, terminal: terminal,
		table: blueprint.table.Clone(), actors: cloneActorsMap(blueprint.actors), nextToMove: nextToMove}
	return &child
}

func (state *KuhnGameState) dealPrivateCards(cardA *cards.Card, cardB *cards.Card) *KuhnGameState {

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
		player.Actions = []Action{}
		return player.Actions
	}

	betSize := state.betSize()
	opponentStack := state.stack(player.Opponent())
	allowedToBet := (player.Stack >= betSize) && (opponentStack >= betSize)

	// whenever betting round is over (CALL OR CHECK->CHECK)
	bettingRoundEnded := state.causingAction.Name() == Call || (state.causingAction.Name() == Check && state.parent.causingAction.Name() == Check)
	if bettingRoundEnded {
		player.Actions = []Action{}
		return player.Actions
	}

	// single check implies BET or CHECK
	if state.causingAction.Name() == Check && state.parent.causingAction.Name() != Check {
		player.Actions = []Action{CheckAction}
		if allowedToBet {
			player.Actions = append(player.Actions, BetAction)
		}
		return player.Actions
	}

	if state.causingAction.Name() == Bet {
		player.Actions = []Action{CallAction, FoldAction}
		return player.Actions
	}

	if state.causingAction.Name() == DealPrivateCards || state.causingAction.Name() == DealPublicCards {
		player.Actions = []Action{CheckAction}
		if allowedToBet {
			player.Actions = append(player.Actions, BetAction)
		}
		return player.Actions
	}
	panic(errors.New("Code not reachable."))
}

func (state *KuhnGameState) playerActor(id ActorID) *Player {
	return state.actors[id].(*Player)
}
